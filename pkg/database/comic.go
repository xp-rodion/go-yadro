package database

type Comic struct {
	Id       int    `json:"id"`
	Url      string `json:"url"`
	Keywords []string
}

func (c *Comic) ToDBEntry() (int, map[string]interface{}) {
	return c.Id, map[string]interface{}{
		"url":      c.Url,
		"keywords": c.Keywords,
	}
}

func ComicFromDBEntry(id int, aboutComic map[string]interface{}) Comic {
	return Comic{
		Id:       id,
		Url:      aboutComic["url"].(string),
		Keywords: aboutComic["keywords"].([]string),
	}
}
