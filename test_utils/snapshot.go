package test_utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	zoptions "github.com/zeet-dev/pkg/utils/options"
)

type SnapshotOption struct {
	UpdateSnapshots bool
	ReadFunc        func(io.ReadCloser) (any, error)
	CompareFunc     func(expected, actual any) bool
	DiffFunc        func(expected, actual any) string
}

var DefaultSnapshotOptions = SnapshotOption{
	UpdateSnapshots: lo.Contains([]string{"true", "1"}, os.Getenv("UPDATE_SNAPSHOTS")),
	ReadFunc: func(f io.ReadCloser) (any, error) {
		data, err := io.ReadAll(f)
		return data, err
	},
	CompareFunc: func(expected, actual any) bool {
		return cmp.Equal(expected, actual)
	},
	DiffFunc: func(expected, actual any) string {
		return cmp.Diff(expected, actual)
	},
}

func snapshot(snapshotFilePath string, actualBytes []byte, opts SnapshotOption) (io.ReadCloser, error) {
	// If UPDATE_SNAPSHOTS is true, update the snapshot file
	if opts.UpdateSnapshots {
		dir := filepath.Dir(snapshotFilePath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, errors.Wrap(err, "failed to create snapshot directory")
			}
		}

		err := os.WriteFile(snapshotFilePath, actualBytes, 0644)
		if err != nil {
			return nil, errors.Wrap(err, "failed to write updated snapshot to file")
		}
		fmt.Printf("Snapshot updated: %s\n", snapshotFilePath)
	}

	// Read expected YAML from file
	reader, err := os.Open(snapshotFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open expected file")
	}

	return reader, nil
}

func SnapshotT(t *testing.T, snapshotFilePath string, actualBytes []byte, options ...zoptions.MustOption[SnapshotOption]) {
	opts := zoptions.MustNewWithDefaults(DefaultSnapshotOptions, options...)

	expectedReader, err := snapshot(snapshotFilePath, actualBytes, opts)
	require.NoError(t, err)

	expectedData, err := opts.ReadFunc(expectedReader)
	require.NoError(t, err)

	actualData, err := opts.ReadFunc(io.NopCloser(bytes.NewReader(actualBytes)))
	require.NoError(t, err)

	// Compare using the provided compare func
	equal := opts.CompareFunc(expectedData, actualData)
	if !equal {
		t.Errorf("Snapshot comparison failed:\n%s", opts.DiffFunc(expectedData, actualData))
	}
}
