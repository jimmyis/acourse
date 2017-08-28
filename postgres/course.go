package postgres

import (
	"context"
	"database/sql"

	"github.com/acoshift/acourse/course"
)

func NewCourseRepository(db *sql.DB) (course.Repository, error) {
	r := &courseRepository{db}

	return r, nil
}

type courseRepository struct {
	db *sql.DB
}

func (r *courseRepository) Count(ctx context.Context) (int64, error) {
	var cnt int64
	err := r.db.QueryRowContext(ctx, `select count(*) from courses`).Scan(&cnt)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}
