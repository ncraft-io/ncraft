package bbolt

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/ncraft-io/ncraft/go/pkg/ncraft/kvstore"
	"github.com/stretchr/testify/require"
	bolt "go.etcd.io/bbolt"
)

func newTestStore(t *testing.T) kvstore.KvStore {
	t.Helper()
	s, err := NewStore(Options{Path: filepath.Join(t.TempDir(), "test.db")})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, s.Close()) })
	return s
}

func TestStoreCRUDAndBatch(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	require.Error(t, s.Put(ctx, nil, []byte("v")))
	require.NoError(t, s.Put(ctx, []byte("a"), []byte("one")))
	got, err := s.Get(ctx, []byte("a"))
	require.NoError(t, err)
	require.Equal(t, []byte("one"), got)
	require.NoError(t, s.BatchPut(ctx, [][]byte{[]byte("b"), []byte("c")}, [][]byte{[]byte("two"), {}}))
	values, err := s.BatchGet(ctx, [][]byte{[]byte("b"), []byte("missing"), []byte("c"), []byte("a")})
	require.NoError(t, err)
	require.Equal(t, [][]byte{[]byte("two"), nil, {}, []byte("one")}, values)
	require.Error(t, s.BatchPut(ctx, [][]byte{[]byte("d")}, nil))
	require.Error(t, s.BatchPut(ctx, [][]byte{[]byte("d"), nil}, [][]byte{[]byte("four"), []byte("invalid")}))
	got, err = s.Get(ctx, []byte("d"))
	require.NoError(t, err)
	require.Nil(t, got, "invalid batch must not partially write")
	_, err = s.BatchGet(ctx, [][]byte{nil, []byte("a")})
	require.Error(t, err)
	require.NoError(t, s.Delete(ctx, []byte("a")))
	require.NoError(t, s.Delete(ctx, []byte("missing")))
	got, err = s.Get(ctx, []byte("a"))
	require.NoError(t, err)
	require.Nil(t, got)
}

func TestScanRangeAndLimit(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	require.NoError(t, s.BatchPut(ctx, [][]byte{[]byte("a"), []byte("b"), []byte("c")}, [][]byte{[]byte("A"), []byte("B"), []byte("C")}))
	for _, tc := range []struct {
		name, start, end string
		limit            int
		want             [][]byte
	}{
		{"one", "a", "", 1, [][]byte{[]byte("a")}},
		{"exclusive end", "a", "c", 10, [][]byte{[]byte("a"), []byte("b")}},
		{"seek", "aa", "", 10, [][]byte{[]byte("b"), []byte("c")}},
		{"empty", "c", "b", 10, nil},
		{"zero", "", "", 0, nil},
		{"negative", "", "", -1, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var end []byte
			if tc.end != "" {
				end = []byte(tc.end)
			}
			keys, values, err := s.Scan(ctx, []byte(tc.start), end, tc.limit)
			require.NoError(t, err)
			require.Equal(t, tc.want, keys)
			require.Len(t, values, len(keys))
		})
	}
	keys, values, err := s.Scan(ctx, nil, nil, 3)
	require.NoError(t, err)
	keys[0][0] = 'z'
	values[0][0] = 'Z'
	got, err := s.Get(ctx, []byte("a"))
	require.NoError(t, err)
	require.Equal(t, []byte("A"), got, "scan buffers must be owned by caller")
}

func TestReadErrorsAfterClose(t *testing.T) {
	s := newTestStore(t)
	require.NoError(t, s.Close())
	_, err := s.Get(context.Background(), []byte("key"))
	require.ErrorIs(t, err, bolt.ErrDatabaseNotOpen)
	_, err = s.BatchGet(context.Background(), [][]byte{[]byte("key")})
	require.ErrorIs(t, err, bolt.ErrDatabaseNotOpen)
}
