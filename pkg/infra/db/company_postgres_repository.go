package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/nyzhehorodov/apicompanies/pkg/domain/company"
)

type CompanyPostgresRepository struct {
	conn *pgxpool.Pool
}

func NewCompanyPostgresRepository(conn *pgxpool.Pool) *CompanyPostgresRepository {
	return &CompanyPostgresRepository{
		conn: conn,
	}
}

func (r *CompanyPostgresRepository) Add(ctx context.Context, raw company.Company) error {
	query := "INSERT INTO companies " +
		"(code, name, country, website, phone) " +
		"VALUES ($1, $2, $3, $4, $5) RETURNING id"

	err := r.conn.QueryRow(ctx, query, raw, raw.Code, raw.Name, raw.Country, raw.Website, raw.Phone).Scan(&raw.ID)
	if err != nil {
		return fmt.Errorf("query exec: %w", err)
	}

	return nil
}

func (r *CompanyPostgresRepository) List(ctx context.Context, opts company.ListOptions) ([]company.Company, error) {
	query := "SELECT id, code, name, country, website, phone " +
		"FROM companies"

	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	var res []company.Company
	for rows.Next() {
		var row company.Company
		err := rows.Scan(&row.ID, &row.Code, &row.Name, &row.Country, &row.Website, &row.Phone)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		res = append(res, row)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return res, nil
}

func (r *CompanyPostgresRepository) Update(ctx context.Context, row *company.Company) error {
	query := "UPDATE companies SET code = $1, name = $2, country = $3, website = $4, phone = $5 " +
		"WHERE id = $6"

	_, err := r.conn.Exec(ctx, query, row.Code, row.Name, row.Country, row.Website, row.Phone, row.ID)
	if err != nil {
		return fmt.Errorf("query exec: %w", err)
	}

	return nil
}

func (r *CompanyPostgresRepository) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM companies WHERE id = $1"

	_, err := r.conn.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("query exec: %w", err)
	}

	return nil
}
