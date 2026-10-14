package badgerdb

import (
	"context"
	"github.com/dgraph-io/badger"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	s, err := NewStore(Options{Dir: t.TempDir(), BucketName: "test"})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, s.Close()) })
	ctx := context.Background()
	require.NoError(t, s.Put(ctx, []byte("key"), []byte("value")))
	got, err := s.Get(ctx, []byte("key"))
	require.NoError(t, err)
	require.Equal(t, []byte("value"), got)
	require.NoError(t, s.Delete(ctx, []byte("key")))
	got, err = s.Get(ctx, []byte("key"))
	require.NoError(t, err)
	require.Nil(t, got)
	require.Error(t, s.Put(ctx, nil, nil))
	_, err = s.Get(ctx, nil)
	require.Error(t, err)
	require.Error(t, s.Delete(ctx, nil))
	require.NoError(t, s.(*Store).db.Update(func(txn *badger.Txn) error { return txn.Set([]byte("other.key"), []byte("other")) }))
	keys, _, err := s.Scan(ctx, nil, nil, 10)
	require.NoError(t, err)
	require.Empty(t, keys, "scan must not cross bucket boundaries")
}
