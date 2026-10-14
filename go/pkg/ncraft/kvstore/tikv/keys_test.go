package tikv

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tikv/client-go/rawkv"
)

type rangeClient struct {
	rawClient
	start, end []byte
	limit      int
	err        error
}

func (c *rangeClient) Scan(_ context.Context, start, end []byte, limit int, _ ...rawkv.ScanOption) ([][]byte, [][]byte, error) {
	c.start, c.end, c.limit = start, end, limit
	return [][]byte{[]byte("bucket.a"), []byte("bucket.b")}, [][]byte{[]byte("A"), []byte("B")}, c.err
}
func (c *rangeClient) DeleteRange(_ context.Context, start, end []byte) error {
	c.start, c.end = start, end
	return c.err
}

func TestBucketKeysAndRanges(t *testing.T) {
	c := &rangeClient{}
	s := &Store{cli: c, bucketName: "bucket"}
	require.Equal(t, []byte("bucket.a"), s.Key([]byte("a")))
	require.Equal(t, [][]byte{[]byte("bucket.a"), []byte("bucket.b")}, s.Keys([][]byte{[]byte("a"), []byte("b")}))
	keys, values, err := s.Scan(context.Background(), []byte("a"), nil, 2)
	require.NoError(t, err)
	require.Equal(t, [][]byte{[]byte("a"), []byte("b")}, keys)
	require.Equal(t, [][]byte{[]byte("A"), []byte("B")}, values)
	require.Equal(t, []byte("bucket.a"), c.start)
	require.Equal(t, []byte("bucket/"), c.end)
	require.Equal(t, 2, c.limit)
	require.NoError(t, s.DeleteRange(context.Background(), nil, nil))
	require.Equal(t, []byte("bucket."), c.start)
	require.Equal(t, []byte("bucket/"), c.end)
	require.NoError(t, s.DeleteRange(context.Background(), []byte("a"), []byte("c")))
	require.Equal(t, []byte("bucket.c"), c.end)
	s.bucketName = ""
	require.NoError(t, s.DeleteRange(context.Background(), nil, nil))
	require.Nil(t, c.end)
}

func TestScanErrorsAndLimits(t *testing.T) {
	expected := errors.New("unavailable")
	c := &rangeClient{err: expected}
	s := &Store{cli: c, bucketName: "bucket"}
	keys, values, err := s.Scan(context.Background(), nil, nil, 1)
	require.ErrorIs(t, err, expected)
	require.Nil(t, keys)
	require.Nil(t, values)
	for _, limit := range []int{0, -1} {
		keys, values, err = s.Scan(context.Background(), nil, nil, limit)
		require.NoError(t, err)
		require.Empty(t, keys)
		require.Empty(t, values)
	}
	require.ErrorIs(t, s.DeleteRange(context.Background(), nil, nil), expected)
}
