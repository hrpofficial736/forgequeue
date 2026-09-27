package postgres

import "github.com/jackc/pgx/v5"

type PostgresInterface interface {
	AddRow(any) bool
	DeleteRow(any) bool
	UpdateRow(any) bool
	FetchRow(any) bool
}

type psql struct {
	conn *pgx.Conn
}

func (p *psql) AddRow(data any) bool {
	return true
}

func (p *psql) DeleteRow(data any) bool {
	return true
}

func (p *psql) UpdateRow(data any) bool {
	return true
}

func (p *psql) FetchRow(data any) bool {
	return true
}
