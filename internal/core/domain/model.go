package domain

type Comic struct {
	URL string `json:"url"`
}

type UpdateResponse struct {
	Total int `json:"total"`
	New   int `json:"new"`
}
