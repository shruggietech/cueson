package source

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type destinationPlan struct {
	asset       preparedAsset
	destination string
	existing    fs.FileInfo
}

func planDestinations(assets []preparedAsset, options RestoreOptions) ([]destinationPlan, error) {
	if options.Output != "" && strings.TrimSpace(options.Output) == "" {
		return nil, preconditionf("--output path must not be whitespace-only")
	}
	if options.OutputDir != "" && strings.TrimSpace(options.OutputDir) == "" {
		return nil, preconditionf("--output-dir path must not be whitespace-only")
	}
	if options.Output != "" && options.OutputDir != "" {
		return nil, preconditionf("--output and --output-dir are mutually exclusive")
	}
	if len(assets) == 0 {
		return nil, fmt.Errorf("source envelope contains no assets")
	}
	if len(assets) > 1 && options.Output != "" {
		return nil, preconditionf("--output accepts only a single-asset document")
	}
	if len(assets) > 1 && options.OutputDir == "" {
		return nil, preconditionf("multi-asset restoration requires --output-dir")
	}
	if options.Metadata > MetadataNone {
		return nil, preconditionf("invalid metadata mode")
	}

	if options.OutputDir != "" {
		info, err := os.Stat(options.OutputDir)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return nil, preconditionf("output directory %q does not exist", options.OutputDir)
			}
			return nil, fmt.Errorf("inspect output directory %q: %w", options.OutputDir, err)
		}
		if !info.IsDir() {
			return nil, preconditionf("output directory %q is not a directory", options.OutputDir)
		}
	}

	plans := make([]destinationPlan, 0, len(assets))
	portablePaths := make(map[string]string, len(assets))
	nativePaths := make(map[string]string, len(assets))
	for index := range assets {
		destination := assets[index].asset.FileName
		switch {
		case options.Output != "":
			destination = options.Output
		case options.OutputDir != "":
			destination = filepath.Join(options.OutputDir, assets[index].asset.FileName)
		}
		if destination == "" {
			return nil, preconditionf("output path must not be empty")
		}
		absolute, err := filepath.Abs(destination)
		if err != nil {
			return nil, preconditionf("resolve output %q: %v", destination, err)
		}
		absolute = filepath.Clean(absolute)
		parent := filepath.Dir(absolute)
		parentInfo, err := os.Stat(parent)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return nil, preconditionf("output parent %q does not exist", parent)
			}
			return nil, fmt.Errorf("inspect output parent %q: %w", parent, err)
		}
		if !parentInfo.IsDir() {
			return nil, preconditionf("output parent %q is not a directory", parent)
		}
		portableKey := portableIdentity(filepath.ToSlash(absolute))
		if prior, exists := portablePaths[portableKey]; exists {
			return nil, preconditionf("destinations %q and %q collide portably", prior, absolute)
		}
		portablePaths[portableKey] = absolute
		nativeKey := filepath.Clean(absolute)
		if prior, exists := nativePaths[nativeKey]; exists {
			return nil, preconditionf("destinations %q and %q resolve to the same path", prior, absolute)
		}
		nativePaths[nativeKey] = absolute

		var existing fs.FileInfo
		info, statErr := os.Lstat(absolute)
		switch {
		case statErr == nil:
			if info.Mode()&fs.ModeSymlink != 0 || !info.Mode().IsRegular() {
				return nil, preconditionf("destination %q is not a regular file", absolute)
			}
			if !options.Force {
				return nil, preconditionf("destination %q already exists; use --force to replace it", absolute)
			}
			for priorIndex := range plans {
				if plans[priorIndex].existing != nil && os.SameFile(info, plans[priorIndex].existing) {
					return nil, preconditionf("destinations %q and %q identify the same existing file", plans[priorIndex].destination, absolute)
				}
			}
			existing = info
		case errors.Is(statErr, fs.ErrNotExist):
		default:
			return nil, fmt.Errorf("inspect destination %q: %w", absolute, statErr)
		}
		plans = append(plans, destinationPlan{asset: assets[index], destination: absolute, existing: existing})
	}
	return plans, nil
}
