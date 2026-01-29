package persistance

type Pagination struct {
	Limit  int
	Offset int
}

type SortingDirection string

const (
	Ascending  SortingDirection = "asc"
	Descending SortingDirection = "desc"
)

func DefaultPagination() Pagination {
	return Pagination{Limit: 100, Offset: 0}
}

type Sorting struct {
	OrderBy   string
	Direction SortingDirection
}
