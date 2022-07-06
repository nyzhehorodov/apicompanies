package company

import (
	"context"
	"fmt"

	"github.com/nyzhehorodov/apicompanies/pkg/domain/company"
)

//go:generate mockgen -destination=./mocks/service.go -package=mocks . Service

type Service interface {
	Add(company company.Company) error
	List(options company.ListOptions) (list []company.Company, count int, err error)
	Update(company *company.Company) error
	Delete(id int) error
}

type svc struct {
	repo company.Repository
}

func NewService(repo company.Repository) Service {
	return &svc{repo: repo}
}

func (s *svc) Add(company company.Company) error {
	if err := s.repo.Add(context.Background(), company); err != nil {
		return fmt.Errorf("add company: %w", err)
	}

	return nil
}

func (s *svc) List(options company.ListOptions) ([]company.Company, int, error) {
	list, err := s.repo.List(context.Background(), options)
	if err != nil {
		return nil, 0, fmt.Errorf("list companies: %w", err)
	}

	return list, len(list), nil
}

func (s *svc) Update(company *company.Company) error {
	if err := s.repo.Update(context.Background(), company); err != nil {
		return fmt.Errorf("update company: %w", err)
	}

	return nil
}

func (s *svc) Delete(id int) error {
	if err := s.repo.Delete(context.Background(), id); err != nil {
		return fmt.Errorf("delete company: %w", err)
	}

	return nil
}
