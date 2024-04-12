package database

type Comic struct {
	Id       string `json:"num"`
	Url      string `json:"img"`
	Keywords []string
}

func (c *Comic) ToDBEntry() (string, map[string]interface{}) {
	return c.Id, map[string]interface{}{
		"url":      c.Url,
		"keywords": c.Keywords,
	}
}

func ComicFromDBEntry(id string, aboutComic map[string]interface{}) Comic {
	return Comic{
		Id:       id,
		Url:      aboutComic["url"].(string),
		Keywords: aboutComic["keywords"].([]string),
	}
}
