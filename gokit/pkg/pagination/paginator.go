package pagination

type Paginator interface {
	GetTotalCount() int32
	GetNextPageToken() string
}

// Paginater deprecated
type Paginater interface {
	GetTotalCount() int32
	GetNextPageToken() string
}
