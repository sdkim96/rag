package graph

type Node struct {
	ID    string                        `json:"id"`
	fnPtr func(StateSchema) StateSchema `json:"-"`
}

type Edge struct {
	ID       string `json:"id"`
	SourceID string `json:"source"`
	TargetID string `json:"target"`
}
