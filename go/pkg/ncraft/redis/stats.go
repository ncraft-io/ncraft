package redis

type Stats struct {
    TotalCount int64
    IdleCount  int64
    StaleCount int64
    Hits       int64
    Misses     int64
    Timeouts   int64
}
