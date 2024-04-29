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
	keywordsInterface, ok := aboutComic["keywords"]
	if !ok {
		keywordsInterface = []interface{}{}
	}

	keywords := make([]string, len(keywordsInterface.([]interface{})))
	for i, keyword := range keywordsInterface.([]interface{}) {
		keywords[i] = keyword.(string)
	}

	return Comic{
		Id:       id,
		Url:      aboutComic["url"].(string),
		Keywords: keywords,
	}
}
