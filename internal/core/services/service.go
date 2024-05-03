package services

import (
	"errors"
	"xkcd/internal/core/domain"
	"xkcd/internal/core/ports"
)

type Service struct {
	comicsRepository ports.ComicsRepository
}

func NewService(comicsRepository ports.ComicsRepository) *Service {
	return &Service{
		comicsRepository: comicsRepository,
	}
}

func (srv *Service) List(proposal string) ([]domain.Comic, error) {
	comics, err := srv.comicsRepository.List(proposal)
	if err != nil {
		return []domain.Comic{}, errors.New("list comics from repository has failed")
	}
	return comics, nil
}
func (srv *Service) Update() (domain.UpdateResponse, error) {
	response, err := srv.comicsRepository.Update()
	if err != nil {
		return domain.UpdateResponse{}, errors.New("update comics failed")
	}
	return response, nil
}
