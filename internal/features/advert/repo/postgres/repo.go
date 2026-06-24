package advertpostgres

import corepgxpool "github.com/egotk/golang-advert-app/internal/core/postgres/pool/pgx"

type Repo struct {
	pool *corepgxpool.Pool
}

func New(pool *corepgxpool.Pool) *Repo {
	return &Repo{
		pool: pool,
	}
}
