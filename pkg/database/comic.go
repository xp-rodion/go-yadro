package database

import (
	"encoding/json"
	"strconv"
	"xkcd/pkg/words"
)

type Comic struct {
	Id       int    `json:"num"`
	Url      string `json:"img"`
	Alt      string `json:"alt"`
	Keywords []string
}

func JsonToComic(byteComic []byte) (Comic, error) {
	comic := new(Comic)
	err := json.Unmarshal(byteComic, &comic)
	if err != nil {
		return *comic, err
	}
	comic.Keywords = words.Stemmer(comic.Alt)
	return *comic, nil
}

func (c *Comic) ComicToMap() map[string]map[string]interface{} {
	comic := map[string]map[string]interface{}{
		strconv.Itoa(c.Id): {
			"url":      c.Url,
			"keywords": c.Keywords,
		},
	}
	return comic
}
