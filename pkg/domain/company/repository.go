package company

import "context"

//go:generate mockgen -destination=./mocks/company_repository.go -package=company . Repository

type Repository interface {
	Add(ctx context.Context, company Company) error
	List(ctx context.Context, options ListOptions) (list []Company, err error)
	Update(ctx context.Context, company *Company) error
	Delete(ctx context.Context, id int) error
}
