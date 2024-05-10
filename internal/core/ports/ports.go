package ports

import "xkcd/internal/core/domain"

type ComicsRepository interface {
	List(proposal string) ([]domain.Comic, error)
	Update() (domain.UpdateResponse, error)
}

type ComicsService interface {
	List(proposal string) ([]domain.Comic, error)
	Update() (domain.UpdateResponse, error)
}
