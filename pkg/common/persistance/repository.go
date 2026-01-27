package persistance

type Pagination struct {
	Limit  int
	Offset int
}

type SortingDirection string

const (
	Ascending        SortingDirection = "asc"
	Descending       SortingDirection = "desc"
	DefaultDirection                  = Ascending
)

type Sorting struct {
	OrderBy   string
	Direction SortingDirection
}
