package postgresql

import (
	"context"
	"fmt"
	"log"
	"time"

	repeatable "github.com/yofukashi/e-commerce/pkg/utils"

	"github.com/yofukashi/e-commerce/internal/config"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	Pool    *pgxpool.Pool
	Builder squirrel.StatementBuilderType
}

func NewClient(ctx context.Context, maxAttempts int, sc config.StorageConfig) (*Postgres, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", sc.Username, sc.Password, sc.Host, sc.Port, sc.Database)
	var pool *pgxpool.Pool
	err := repeatable.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		var err error
		pool, err = pgxpool.Connect(ctx, dsn)
		if err != nil {
			return err
		}

		return nil
	}, maxAttempts, 5*time.Second)

	if err != nil {
		log.Fatal("error do with tries postgresql")
	}
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	return &Postgres{Pool: pool, Builder: builder}, nil
}
