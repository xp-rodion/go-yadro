package xkcd

type Entry struct {
	Id    int    `json:"num"`
	Url   string `json:"img"`
	Alt   string `json:"alt"`
	Title string `json:"title"`
}
