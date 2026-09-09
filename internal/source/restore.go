package source

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/shruggietech/cueson/internal/model"
)

type stagedAsset struct {
	plan       destinationPlan
	stagePath  string
	stageInfo  fs.FileInfo
	backupPath string
	committed  fs.FileInfo
}

type transactionHooks struct {
	fail func(point, path string) error
}

const (
	stagePattern  = ".cueson-stage-*"
	backupPattern = ".cueson-backup-*"
)

func (hooks transactionHooks) at(point, path string) error {
	if hooks.fail == nil {
		return nil
	}
	return hooks.fail(point, path)
}

// Restore validates, plans, stages, commits, verifies, and applies metadata to a complete bundle.
func Restore(ctx context.Context, document model.Document, options RestoreOptions) (report Report, result error) {
	return restoreWithHooks(ctx, document, options, transactionHooks{})
}

func restoreWithHooks(ctx context.Context, document model.Document, options RestoreOptions, hooks transactionHooks) (report Report, result error) {
	if err := document.Validate(); err != nil {
		return report, fmt.Errorf("validate Cue JSON semantics: %w", err)
	}
	assets, err := prepareAssetsContext(ctx, document)
	if err != nil {
		return report, err
	}
	plans, err := planDestinations(assets, options)
	if err != nil {
		return report, err
	}
	staged := make([]stagedAsset, 0, len(plans))
	defer func() {
		if result != nil {
			failures := []error{result}
			if rollbackErr := rollback(staged, hooks); rollbackErr != nil {
				failures = append(failures, fmt.Errorf("rollback: %w", rollbackErr))
			}
			if cleanupErr := cleanupFailedStages(staged, hooks); cleanupErr != nil {
				failures = append(failures, fmt.Errorf("cleanup failed staging files: %w", cleanupErr))
			}
			result = errors.Join(failures...)
			return
		}
		for index := range staged {
			if staged[index].stagePath != "" {
				_ = os.Remove(staged[index].stagePath)
			}
		}
	}()

	for index := range plans {
		if err := ctx.Err(); err != nil {
			return report, fmt.Errorf("restore canceled: %w", err)
		}
		item, err := stage(ctx, plans[index], hooks)
		if err != nil {
			return report, err
		}
		staged = append(staged, item)
	}
	for index := range staged {
		if err := ctx.Err(); err != nil {
			return report, fmt.Errorf("restore canceled: %w", err)
		}
		if err := commit(&staged[index], hooks); err != nil {
			return report, err
		}
	}

	report.Assets = make([]AssetResult, 0, len(staged))
	for index := range staged {
		item := &staged[index]
		if err := hooks.at("final-verify", item.plan.destination); err != nil {
			return report, fmt.Errorf("verify final destination %q: %w", item.plan.destination, err)
		}
		count, digest, err := inspectFileContext(ctx, item.plan.destination)
		if err != nil {
			return report, fmt.Errorf("verify final destination %q: %w", item.plan.destination, err)
		}
		if count != item.plan.asset.asset.Size.Bytes || digest != item.plan.asset.asset.Hashes.SHA256 {
			return report, fmt.Errorf("final destination %q failed byte-integrity verification", item.plan.destination)
		}
		timestamps, warnings, err := timestampResults(item.plan.destination, item.plan.asset.asset.Timestamps, options.Metadata)
		if err != nil {
			return report, err
		}
		report.Warnings = append(report.Warnings, warnings...)
		report.Assets = append(report.Assets, AssetResult{AssetID: item.plan.asset.asset.ID, Destination: item.plan.destination, Bytes: count, SHA256: digest, Timestamps: timestamps})
	}

	for index := range staged {
		item := &staged[index]
		if item.backupPath != "" {
			if err := hooks.at("cleanup-backup", item.backupPath); err != nil {
				report.Warnings = append(report.Warnings, fmt.Sprintf("remove accepted rollback backup %q: %v", item.backupPath, err))
			} else if err := os.Remove(item.backupPath); err != nil {
				report.Warnings = append(report.Warnings, fmt.Sprintf("remove accepted rollback backup %q: %v", item.backupPath, err))
			} else {
				item.backupPath = ""
			}
		}
		if item.stagePath != "" {
			if err := hooks.at("cleanup-stage", item.stagePath); err != nil {
				report.Warnings = append(report.Warnings, fmt.Sprintf("remove accepted staging link %q: %v", item.stagePath, err))
			} else if err := os.Remove(item.stagePath); err != nil {
				report.Warnings = append(report.Warnings, fmt.Sprintf("remove accepted staging link %q: %v", item.stagePath, err))
			} else {
				item.stagePath = ""
			}
		}
	}
	return report, nil
}

func stage(ctx context.Context, plan destinationPlan, hooks transactionHooks) (item stagedAsset, result error) {
	file, err := os.CreateTemp(filepath.Dir(plan.destination), stagePattern)
	if err != nil {
		return stagedAsset{}, fmt.Errorf("create staging file for %q: %w", plan.destination, err)
	}
	path := file.Name()
	remove := true
	defer func() {
		if remove {
			if err := hooks.at("failure-cleanup-stage", path); err != nil {
				result = errors.Join(result, fmt.Errorf("remove failed staging file %q: %w", path, err))
			} else if err := os.Remove(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
				result = errors.Join(result, fmt.Errorf("remove failed staging file %q: %w", path, err))
			}
		}
	}()
	if err := file.Chmod(0o600); err != nil {
		_ = file.Close()
		return stagedAsset{}, fmt.Errorf("secure staging file for %q: %w", plan.destination, err)
	}
	hash := sha256.New()
	if err := hooks.at("stage-write", path); err != nil {
		_ = file.Close()
		return stagedAsset{}, fmt.Errorf("write staging file for %q: %w", plan.destination, err)
	}
	count, copyErr := io.Copy(io.MultiWriter(file, hash), contextReader{ctx: ctx, reader: base64.NewDecoder(base64.StdEncoding.Strict(), strings.NewReader(plan.asset.asset.DataBase64))})
	if copyErr != nil {
		_ = file.Close()
		return stagedAsset{}, fmt.Errorf("decode source asset %q: %w", plan.asset.asset.ID, copyErr)
	}
	if syncErr := hooks.at("stage-sync", path); syncErr != nil {
		_ = file.Close()
		return stagedAsset{}, fmt.Errorf("flush staging file for %q: %w", plan.destination, syncErr)
	}
	if syncErr := file.Sync(); syncErr != nil {
		_ = file.Close()
		return stagedAsset{}, fmt.Errorf("flush staging file for %q: %w", plan.destination, syncErr)
	}
	if closeErr := hooks.at("stage-close", path); closeErr != nil {
		_ = file.Close()
		return stagedAsset{}, fmt.Errorf("close staging file for %q: %w", plan.destination, closeErr)
	}
	if closeErr := file.Close(); closeErr != nil {
		return stagedAsset{}, fmt.Errorf("close staging file for %q: %w", plan.destination, closeErr)
	}
	if count != plan.asset.asset.Size.Bytes || fmt.Sprintf("%x", hash.Sum(nil)) != plan.asset.asset.Hashes.SHA256 {
		return stagedAsset{}, fmt.Errorf("staged asset %q failed byte-integrity verification", plan.asset.asset.ID)
	}
	verifiedCount, digest, err := inspectFileContext(ctx, path)
	if err != nil || verifiedCount != count || digest != plan.asset.asset.Hashes.SHA256 {
		return stagedAsset{}, fmt.Errorf("reopen staging file for %q failed verification: %v", plan.destination, err)
	}
	if err := hooks.at("stage-inspect", path); err != nil {
		return stagedAsset{}, fmt.Errorf("inspect staging file for %q: %w", plan.destination, err)
	}
	stageInfo, err := os.Lstat(path)
	if err != nil {
		return stagedAsset{}, fmt.Errorf("inspect staging file for %q: %w", plan.destination, err)
	}
	if !stageInfo.Mode().IsRegular() {
		return stagedAsset{}, fmt.Errorf("inspect staging file for %q: staging path is not a regular file", plan.destination)
	}
	remove = false
	return stagedAsset{plan: plan, stagePath: path, stageInfo: stageInfo}, nil
}

func commit(item *stagedAsset, hooks transactionHooks) error {
	currentStage, err := os.Lstat(item.stagePath)
	if err != nil || !currentStage.Mode().IsRegular() || !os.SameFile(currentStage, item.stageInfo) {
		return fmt.Errorf("staging file for %q changed before commit", item.plan.destination)
	}
	if item.plan.existing == nil {
		if err := hooks.at("publish-link", item.plan.destination); err != nil {
			return fmt.Errorf("publish destination %q: %w", item.plan.destination, err)
		}
		if err := os.Link(item.stagePath, item.plan.destination); err != nil {
			if errors.Is(err, fs.ErrExist) {
				return fmt.Errorf("destination %q appeared after planning", item.plan.destination)
			}
			return fmt.Errorf("publish destination %q: %w", item.plan.destination, err)
		}
		item.committed = item.stageInfo
	} else {
		current, err := os.Lstat(item.plan.destination)
		if err != nil || !current.Mode().IsRegular() || !os.SameFile(current, item.plan.existing) {
			return fmt.Errorf("destination %q changed after planning", item.plan.destination)
		}
		backup, err := createBackupLink(item.plan.destination)
		if err != nil {
			return fmt.Errorf("preserve destination %q for rollback: %w", item.plan.destination, err)
		}
		item.backupPath = backup
		current, err = os.Lstat(item.plan.destination)
		if err != nil || !current.Mode().IsRegular() || !os.SameFile(current, item.plan.existing) {
			return fmt.Errorf("destination %q changed while preparing replacement", item.plan.destination)
		}
		if err := hooks.at("replace-rename", item.plan.destination); err != nil {
			return fmt.Errorf("replace destination %q: %w", item.plan.destination, err)
		}
		if err := os.Rename(item.stagePath, item.plan.destination); err != nil {
			return fmt.Errorf("replace destination %q: %w", item.plan.destination, err)
		}
		item.committed = item.stageInfo
		item.stagePath = ""
	}
	committed, err := os.Lstat(item.plan.destination)
	if err != nil || !committed.Mode().IsRegular() {
		return fmt.Errorf("inspect committed destination %q: %w", item.plan.destination, err)
	}
	if !os.SameFile(committed, item.committed) {
		return fmt.Errorf("destination %q changed during commit", item.plan.destination)
	}
	return nil
}

func createBackupLink(destination string) (string, error) {
	for attempt := 0; attempt < 8; attempt++ {
		placeholder, err := os.CreateTemp(filepath.Dir(destination), backupPattern)
		if err != nil {
			return "", err
		}
		path := placeholder.Name()
		if closeErr := placeholder.Close(); closeErr != nil {
			_ = os.Remove(path)
			return "", closeErr
		}
		if removeErr := os.Remove(path); removeErr != nil {
			return "", removeErr
		}
		if linkErr := os.Link(destination, path); linkErr == nil {
			return path, nil
		} else if !errors.Is(linkErr, fs.ErrExist) {
			return "", linkErr
		}
	}
	return "", fmt.Errorf("could not reserve an exclusive rollback path")
}

func rollback(staged []stagedAsset, hooks transactionHooks) error {
	var failures []error
	for index := len(staged) - 1; index >= 0; index-- {
		item := &staged[index]
		if item.committed == nil {
			if item.backupPath != "" {
				if err := os.Remove(item.backupPath); err != nil && !errors.Is(err, fs.ErrNotExist) {
					failures = append(failures, err)
				}
				item.backupPath = ""
			}
			continue
		}
		current, err := os.Lstat(item.plan.destination)
		if err != nil || !current.Mode().IsRegular() || !os.SameFile(current, item.committed) {
			failures = append(failures, fmt.Errorf("destination %q changed before rollback; backup retained at %q", item.plan.destination, item.backupPath))
			continue
		}
		if item.backupPath == "" {
			if err := hooks.at("rollback-remove-new", item.plan.destination); err != nil {
				failures = append(failures, err)
			} else if err := os.Remove(item.plan.destination); err != nil {
				failures = append(failures, err)
			}
			continue
		}
		if err := hooks.at("rollback-remove-replacement", item.plan.destination); err != nil {
			failures = append(failures, err)
			continue
		}
		if err := os.Remove(item.plan.destination); err != nil {
			failures = append(failures, err)
			continue
		}
		if err := hooks.at("rollback-restore-backup", item.plan.destination); err != nil {
			failures = append(failures, fmt.Errorf("restore backup %q: %w", item.backupPath, err))
			continue
		}
		if err := os.Rename(item.backupPath, item.plan.destination); err != nil {
			failures = append(failures, fmt.Errorf("restore backup %q: %w", item.backupPath, err))
			continue
		}
		item.backupPath = ""
	}
	return errors.Join(failures...)
}

func cleanupFailedStages(staged []stagedAsset, hooks transactionHooks) error {
	var failures []error
	for index := range staged {
		if staged[index].stagePath == "" {
			continue
		}
		path := staged[index].stagePath
		if err := hooks.at("failure-cleanup-stage", path); err != nil {
			failures = append(failures, fmt.Errorf("remove %q: %w", path, err))
			continue
		}
		if err := os.Remove(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
			failures = append(failures, fmt.Errorf("remove %q: %w", path, err))
			continue
		}
		staged[index].stagePath = ""
	}
	return errors.Join(failures...)
}

func inspectFile(path string) (int64, string, error) {
	return inspectFileContext(context.Background(), path)
}

func inspectFileContext(ctx context.Context, path string) (int64, string, error) {
	file, err := openRegularNoFollow(path, false)
	if err != nil {
		return 0, "", err
	}
	hash := sha256.New()
	count, copyErr := io.Copy(hash, contextReader{ctx: ctx, reader: file})
	closeErr := file.Close()
	if copyErr != nil {
		return 0, "", copyErr
	}
	if closeErr != nil {
		return 0, "", closeErr
	}
	return count, fmt.Sprintf("%x", hash.Sum(nil)), nil
}

type contextReader struct {
	ctx    context.Context
	reader io.Reader
}

func (reader contextReader) Read(payload []byte) (int, error) {
	if err := reader.ctx.Err(); err != nil {
		return 0, err
	}
	count, err := reader.reader.Read(payload)
	if contextErr := reader.ctx.Err(); contextErr != nil {
		return count, contextErr
	}
	return count, err
}
