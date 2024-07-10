package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/PolyAbit/auth/internal/models"
)

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.sqlite.IsAdmin"

	stmt, err := s.db.Prepare("SELECT is_admin FROM users WHERE id = ?")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, userID)

	var isAdmin bool
	err = row.Scan(&isAdmin)
	if errors.Is(err, sql.ErrNoRows) {
		return false, fmt.Errorf("%s: %w", op, models.ErrUserNotFound)
	}
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}
