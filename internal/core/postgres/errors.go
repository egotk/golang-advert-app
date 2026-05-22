package corepostgres

import (
	"errors"
	"fmt"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	codeUniqueViolation = "23505"
	codeCheckViolation  = "23514"
)

func MapError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return coreerrors.ErrNotFound
	}

	var pgErr *pgconn.PgError

	if !errors.As(err, &pgErr) {
		return err
	}
	switch pgErr.Code {
	case codeUniqueViolation:
		return fmt.Errorf("%s: %w", pgErr.ConstraintName, coreerrors.ErrConflict)
	case codeCheckViolation:
		return fmt.Errorf("%s: %w", pgErr.ConstraintName, coreerrors.ErrInvalidArgument)
	default:
		return err
	}
}
