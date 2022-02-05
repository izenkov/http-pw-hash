package api

type Stats struct {
	Total   int64 `json:"total"`
	Average int64 `json:"average"`
}

func NewStats(total int64, average int64) *Stats {
	s := &Stats{
		Total:   total,
		Average: average,
	}
	return s
}
