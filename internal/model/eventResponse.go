package model

type StateResponse struct {
	State       string `json:"state,omitempty"`
	Description string `json:"description,omitempty"`
}

type IdResponse struct {
	Id string `json:"id"`
}
