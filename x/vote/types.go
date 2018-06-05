package vote

type AccumulatedProjectVote struct {
	ProjectID string `json:"project_id"`
	Ammount   int64  `json:"ammount"`
}
