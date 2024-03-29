package kvstore

import (
    "context"
    "github.com/mojo-lang/core/go/pkg/mojo/core"
    "github.com/ncraft-io/ncraft/go/pkg/ncraft/logs"
    "sync"
)

const DefaultBatchSize = 10000

type Options = core.Options

type KeysValues struct {
    Keys   [][]byte
    Values [][]byte
}

type KvStore interface {
    // Close the storage
    Close() error

    Get(ctx context.Context, key []byte) ([]byte, error)

    BatchGet(ctx context.Context, keys [][]byte) ([][]byte, error)

    Put(ctx context.Context, key, value []byte) error

    BatchPut(ctx context.Context, keys, values [][]byte) error

    Delete(ctx context.Context, key []byte) error

    BatchDelete(ctx context.Context, keys [][]byte) error

    // Scan
    // GetRange
    Scan(ctx context.Context, startKey, endKey []byte, limit int) (keys [][]byte, values [][]byte, err error)

    // DeleteRange delete the range key
    DeleteRange(ctx context.Context, startKey []byte, endKey []byte) error
}

var (
    stores map[string]func(Options) KvStore
    lock   sync.Mutex
)

// RegisterStore register a reader function to process a command
func RegisterStore(name string, reader func(Options) KvStore) error {
    lock.Lock()
    defer lock.Unlock()

    if stores == nil {
        stores = make(map[string]func(Options) KvStore)
    }

    if _, ok := stores[name]; ok {
        logs.Warnw("register the same KvStore", "name", name)
        return nil
    }

    stores[name] = reader
    return nil
}

func NewStore(name string, options Options) KvStore {
    lock.Lock()
    defer lock.Unlock()

    if store, ok := stores[name]; ok {
        return store(options)
    }
    return nil
}

func Scan(kv KvStore, startKey, endKey []byte, valuesChan chan<- *KeysValues, batchLimit int) {
    if batchLimit == 0 {
        batchLimit = DefaultBatchSize
    }
    for {
        keys, values, err := kv.Scan(context.TODO(), startKey, endKey, batchLimit)
        if err != nil {
            logs.Errorf("failed to scan start from %s, end %s, err: %s", startKey, endKey, err)
            continue
        }
        keysNum := len(keys)
        if keysNum == 0 {
            break
        }

        if keysNum < batchLimit {
            valuesChan <- &KeysValues{
                Keys:   keys,
                Values: values,
            }
            break
        }

        // avoid to getting duplicated value by minus one
        startKey = keys[len(keys)-1]
        valuesChan <- &KeysValues{
            Keys:   keys[0 : len(keys)-1],
            Values: values[0 : len(values)-1],
        }
    }
}
