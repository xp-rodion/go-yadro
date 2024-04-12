package xkcd

type Entry struct {
	Id  string `json:"num"`
	Url string `json:"img"`
	Alt string `json:"alt"`
}
