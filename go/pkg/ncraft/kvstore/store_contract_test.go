package kvstore_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/ncraft-io/ncraft/go/pkg/ncraft/kvstore"
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/kvstore/badgerdb"
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/kvstore/bbolt"
	"github.com/stretchr/testify/require"
)

func TestLocalStoreContract(t *testing.T) {
	for _, backend := range []string{"bbolt", "badgerdb"} {
		t.Run(backend, func(t *testing.T) {
			var s kvstore.KvStore
			var err error
			dir := t.TempDir()
			if backend == "bbolt" {
				s, err = bbolt.NewStore(bbolt.Options{Path: filepath.Join(dir, "test.db")})
			} else {
				s, err = badgerdb.NewStore(badgerdb.Options{Dir: dir, BucketName: "test"})
			}
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, s.Close()) })
			ctx := context.Background()
			keys := [][]byte{[]byte("a"), []byte("b"), []byte("c"), []byte("d")}
			values := [][]byte{[]byte("A"), []byte("B"), []byte("C"), []byte("D")}
			require.NoError(t, s.BatchPut(ctx, keys, values))
			got, err := s.BatchGet(ctx, [][]byte{keys[2], []byte("missing"), keys[0]})
			require.NoError(t, err)
			require.Equal(t, [][]byte{values[2], nil, values[0]}, got)
			require.Error(t, s.BatchPut(ctx, keys, nil))
			require.Error(t, s.BatchPut(ctx, [][]byte{[]byte("new"), nil}, [][]byte{[]byte("v"), []byte("v")}))
			value, err := s.Get(ctx, []byte("new"))
			require.NoError(t, err)
			require.Nil(t, value)
			require.Error(t, s.BatchDelete(ctx, [][]byte{keys[0], nil}))
			value, err = s.Get(ctx, keys[0])
			require.NoError(t, err)
			require.Equal(t, values[0], value, "invalid delete must be atomic")
			scanned, contents, err := s.Scan(ctx, keys[1], keys[3], 1)
			require.NoError(t, err)
			require.Equal(t, keys[1:2], scanned)
			require.Equal(t, values[1:2], contents)
			require.NoError(t, s.DeleteRange(ctx, keys[1], keys[3]))
			scanned, _, err = s.Scan(ctx, nil, nil, 10)
			require.NoError(t, err)
			require.Equal(t, [][]byte{keys[0], keys[3]}, scanned)
			require.NoError(t, s.BatchDelete(ctx, [][]byte{keys[0], []byte("missing")}))
			require.NoError(t, s.DeleteRange(ctx, keys[3], nil))
			scanned, _, err = s.Scan(ctx, nil, nil, 10)
			require.NoError(t, err)
			require.Empty(t, scanned)
		})
	}
}

func TestBboltRegistration(t *testing.T) {
	for _, name := range []string{"bbolt", "bblot"} {
		s := kvstore.NewStore(name, kvstore.Options{"path": filepath.Join(t.TempDir(), "test.db")})
		require.NotNil(t, s)
		require.NoError(t, s.Close())
	}
	require.Nil(t, kvstore.NewStore("missing", nil))
	require.Nil(t, kvstore.NewStore("bbolt", kvstore.Options{"path": []string{"invalid"}}))
}
