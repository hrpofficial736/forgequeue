package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func SetupPostgres(ctx context.Context) *pgx.Conn {
	client, err := pgx.Connect(ctx, "")

	if err != nil {
		panic(err)
	}

	// create jobs table
	_, err = client.Query(ctx, "")

	return client
}
