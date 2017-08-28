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

func (r *courseRepository) FindID(ctx context.Context, id string) (*course.Course, error) {
	var x course.Course
	var (
		start pq.NullTime
		url   pq.NullString
	)
	err := r.db.QueryRowContext(ctx, `
		select
			c.id, c.user_id, c.title, c.short_desc, c.long_desc, c.image, c.start,
			c.url, c.type, c.price, c.discount, c.enroll_detail,
			opt.public, opt.enroll, opt.attend, opt.assignment, opt.discount
		from
			courses as c
			left join course_options as opt on c.id = opt.course_id
		where id = $1
	`, id).Scan(
		&x.ID, &x.Owner.ID, &x.Title, &x.ShortDesc, &x.Desc, &x.Image, &start,
		&url, &x.Type, &x.Price, &x.Discount, &x.EnrollDetail,
		&x.Option.Public, &x.Option.Enroll, &x.Option.Attend, &x.Option.Assignment, &x.Option.Discount,
	)
	if err != nil {
		return nil, err
	}
	x.Start = start.Time
	x.URL = url.String
	return &x, nil
}

func (r *courseRepository) Count(ctx context.Context) (int64, error) {
	var cnt int64
	err := r.db.QueryRowContext(ctx, `select count(*) from courses`).Scan(&cnt)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}
