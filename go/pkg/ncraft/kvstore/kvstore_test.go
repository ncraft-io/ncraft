package kvstore

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type scanStore struct {
	KvStore
	scan func([]byte, []byte, int) ([][]byte, [][]byte, error)
}

func (s scanStore) Scan(_ context.Context, start, end []byte, limit int) ([][]byte, [][]byte, error) {
	return s.scan(start, end, limit)
}

func TestScanAdvancesWithSingleItemBatches(t *testing.T) {
	calls := 0
	s := scanStore{scan: func(start, end []byte, limit int) ([][]byte, [][]byte, error) {
		calls++
		require.LessOrEqual(t, calls, 4, "scan must advance")
		require.Equal(t, 1, limit)
		require.Equal(t, []byte("z"), end)
		switch string(start) {
		case "":
			return [][]byte{[]byte("a")}, [][]byte{[]byte("A")}, nil
		case "a\x00":
			return [][]byte{[]byte("a\x00")}, [][]byte{[]byte("B")}, nil
		case "a\x00\x00":
			return nil, nil, nil
		default:
			t.Fatalf("unexpected start key: %q", start)
			return nil, nil, nil
		}
	}}
	ch := make(chan *KeysValues, 4)
	Scan(s, nil, []byte("z"), ch, 1)
	require.Len(t, ch, 2)
	require.Equal(t, [][]byte{[]byte("a")}, (<-ch).Keys)
	require.Equal(t, [][]byte{[]byte("a\x00")}, (<-ch).Keys)
}

func TestScanStopsOnErrorAndDefaultsInvalidLimit(t *testing.T) {
	for _, limit := range []int{0, -1} {
		calls := 0
		s := scanStore{scan: func(_, _ []byte, n int) ([][]byte, [][]byte, error) {
			calls++
			require.Equal(t, DefaultBatchSize, n)
			require.Equal(t, 1, calls, "persistent errors must not spin forever")
			return nil, nil, errors.New("unavailable")
		}}
		Scan(s, nil, nil, make(chan *KeysValues, 1), limit)
	}
}

func TestBucketKeyAndOptions(t *testing.T) {
	require.Equal(t, []byte("bucket.key"), BucketKey("bucket", []byte("key")))
	require.Equal(t, []byte("key"), BucketKey("", []byte("key")))
	var options struct {
		Count int `json:"count"`
	}
	require.NoError(t, ResetOptions(map[string]interface{}{"count": 2}, &options))
	require.Equal(t, 2, options.Count)
	require.Error(t, ResetOptions(map[string]interface{}{"count": "bad"}, &options))
	require.Error(t, CheckKey(nil))
	require.NoError(t, CheckKey([]byte("key")))
	require.Error(t, CheckValue(nil))
	require.NoError(t, CheckValue([]byte{}))
}

func TestStoreFactoryCanUseRegistry(t *testing.T) {
	name := t.Name()
	require.NoError(t, RegisterStore(name, func(Options) KvStore {
		// Factories can resolve/register another backend without deadlocking.
		require.Nil(t, NewStore(name+"-missing", nil))
		require.NoError(t, RegisterStore(name+"-child", func(Options) KvStore { return nil }))
		return nil
	}))
	require.Nil(t, NewStore(name, nil))
}
