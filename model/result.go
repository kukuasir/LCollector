package model

type Result struct {
	Status  int64  `json:"status"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	Total   int64  `json:"total,omitempty"`
}
