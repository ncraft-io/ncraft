package counter

type BatchRequest interface {
    GetBatchCount() int64
}
