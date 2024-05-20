package repositories

import (
	"xkcd/internal/core/domain"
	"xkcd/internal/parse"
	"xkcd/internal/search"
	"xkcd/internal/utils"
	"xkcd/pkg/database"
	"xkcd/pkg/xkcd"
)

type SQLiteRepository struct {
	db     database.SQLite
	client xkcd.Client
}

func NewSQLiteRepository(db database.SQLite, client xkcd.Client) *SQLiteRepository {
	return &SQLiteRepository{db: db, client: client}
}

func (repo *SQLiteRepository) List(proposal string) ([]domain.Comic, error) {
	results := search.DBRelevantComics(repo.db, proposal)
	comics := make([]domain.Comic, 0)
	for _, result := range results {
		if result.Weight > 0 {
			comics = append(comics, domain.Comic{
				URL: result.Url,
			})
		}
	}
	return comics, nil
}

func (repo *SQLiteRepository) Update() (domain.UpdateResponse, error) {
	newComics := utils.CheckNewComics(repo.db, repo.client)
	comics := parse.ParallelParseComics(repo.client, newComics, 6)
	repo.db.Adds(comics)
	response := domain.UpdateResponse{
		Total: repo.client.ComicsCount,
		New:   len(newComics),
	}
	return response, nil
}
