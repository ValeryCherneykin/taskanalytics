package pg

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type pg struct {
	dbc *pgxpool.Pool
}

// func NewDB(dbc *pgxpool.Pool) db.DB {
// 	return &pg{
// 		dbc: dbc,
// 	}
// }

func (p *pg) Ping(ctx context.Context) error {
	return p.dbc.Ping(ctx)
}

func (p *pg) Close() {
	p.dbc.Close()
}
