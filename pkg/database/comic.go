package database

type Comic struct {
	Id       int    `json:"id" db:"id"`
	Url      string `json:"url" db:"url"`
	Keywords string `db:"keywords`
}

func (c *Comic) ToDBEntry() (int, map[string]interface{}) {
	return c.Id, map[string]interface{}{
		"url":      c.Url,
		"keywords": c.Keywords,
	}
}

func ComicFromDBEntry(id int, aboutComic map[string]interface{}) Comic {
	var keywords string
	keywordsInterface, ok := aboutComic["keywords"]
	if !ok {
		keywords = ""
	}

	keywords, ok = keywordsInterface.(string)
	if !ok {
		for _, keyword := range keywordsInterface.([]interface{}) {
			keywords = keyword.(string)
		}
	}

	return Comic{
		Id:       id,
		Url:      aboutComic["url"].(string),
		Keywords: keywords,
	}
}
