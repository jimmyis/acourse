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

func (r *courseRepository) Store(ctx context.Context, course *course.Course) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if len(course.ID) == 0 {
		err = tx.QueryRowContext(ctx, `
			insert into courses (
				title, short_desc, long_desc, image, start,
				url, type, price, discount, enroll_detail
			) values (
				$1, $2, $3, $4, $5,
				$6, $7, $8, $9, $10
			) returning id
		`).Scan(&course.ID)
		if err != nil {
			return err
		}
	} else {
		_, err = tx.ExecContext(ctx, `
			update courses
				set
					title = $2,
					short_desc = $3,
					long_desc = $4,
					image = $5,
					start = $6,
					url = $7,
					type = $8,
					price = $9,
					discount = $10,
					enroll_detail = $11,
					updated_at = now()
				where id = $1
		`,
			course.ID,
			course.Title,
			course.ShortDesc,
			course.Desc,
			course.Image,
			course.Start,
			course.URL,
			course.Type,
			course.Price,
			course.Discount,
			course.EnrollDetail,
		)
		if err != nil {
			return err
		}
	}
	_, err = tx.ExecContext(ctx, `
		insert into course_options (
			course_id,
			public, enroll, attend, assignment, discount
		) values (
			$1,
			$2, $3, $4, $5, $6
		) on conflict (course_id) do update
			set
				public = excluded.public,
				enroll = excluded.enroll,
				attend = excluded.attend,
				assignment = excluded.assignment,
				discount = excluded.discount
	`)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (r *courseRepository) FindID(ctx context.Context, id string) (*course.Course, error) {
	var x course.Course
	err := r.db.QueryRowContext(ctx, `
		select
			c.id, c.title, c.short_desc, c.long_desc, c.image, c.start,
			c.url, c.type, c.price, c.discount, c.enroll_detail,
			u.id, u.username, u.name, u.image,
			opt.public, opt.enroll, opt.attend, opt.assignment, opt.discount
		from
			courses as c
			left join users as u on c.user_id = u.id
			left join course_options as opt on c.id = opt.course_id
		where id = $1
	`, id).Scan(
		&x.ID, &x.Title, &x.ShortDesc, &x.Desc, &x.Image, &x.Start,
		&x.URL, &x.Type, &x.Price, &x.Discount, &x.EnrollDetail,
		&x.Owner.ID, &x.Owner.Username, &x.Owner.Name, &x.Owner.Image,
		&x.Option.Public, &x.Option.Enroll, &x.Option.Attend, &x.Option.Assignment, &x.Option.Discount,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
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
