package source

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/shruggietech/cueson/internal/model"
)

func prepareAssets(document model.Document) ([]preparedAsset, error) {
	return prepareAssetsContext(context.Background(), document)
}

func prepareAssetsContext(ctx context.Context, document model.Document) ([]preparedAsset, error) {
	prepared := make([]preparedAsset, 0, len(document.Source.Assets))
	identities := make(map[string]string, len(document.Source.Assets))
	for index := range document.Source.Assets {
		asset := document.Source.Assets[index]
		if err := validateSafeBasename(asset.FileName); err != nil {
			return nil, fmt.Errorf("source.assets[%d]: %w", index, err)
		}
		key := portableIdentity(asset.FileName)
		if prior, exists := identities[key]; exists {
			return nil, fmt.Errorf("source assets %q and %q have colliding portable basenames", prior, asset.FileName)
		}
		identities[key] = asset.FileName
		count, digest, err := inspectEncodedContext(ctx, asset.DataBase64)
		if err != nil {
			return nil, fmt.Errorf("source asset %q data_base64: %w", asset.ID, err)
		}
		if count != asset.Size.Bytes {
			return nil, fmt.Errorf("source asset %q decoded length is %d, want %d", asset.ID, count, asset.Size.Bytes)
		}
		if digest != asset.Hashes.SHA256 {
			return nil, fmt.Errorf("source asset %q SHA-256 is %s, want %s", asset.ID, digest, asset.Hashes.SHA256)
		}
		prepared = append(prepared, preparedAsset{asset: asset})
	}
	return prepared, nil
}

func inspectEncoded(encoded string) (int64, string, error) {
	return inspectEncodedContext(context.Background(), encoded)
}

func inspectEncodedContext(ctx context.Context, encoded string) (int64, string, error) {
	if strings.ContainsAny(encoded, "\r\n") || len(encoded)%4 != 0 {
		return 0, "", fmt.Errorf("encoding is not canonical standard base64")
	}
	hash := sha256.New()
	count, err := io.Copy(hash, contextReader{ctx: ctx, reader: base64.NewDecoder(base64.StdEncoding.Strict(), strings.NewReader(encoded))})
	if err != nil {
		return 0, "", fmt.Errorf("decode canonical base64: %w", err)
	}
	return count, fmt.Sprintf("%x", hash.Sum(nil)), nil
}
