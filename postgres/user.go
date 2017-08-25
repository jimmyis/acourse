package postgres

import (
	"context"
	"database/sql"

	"github.com/acoshift/acourse/user"
)

// NewUserRepository returns a new instance of a Postgres user repository
func NewUserRepository(db *sql.DB) (user.Repository, error) {
	r := &userRepository{db}

	// create table
	_, err := db.Exec(`
		create table if not exists users (
			id varchar not null,
			username varchar not null,
			name varchar not null,
			email varchar,
			about_me varchar not null default '',
			image varchar not null default '',
			created_at timestamp not null default now(),
			updated_at timestamp not null default now(),
			primary key (id)
		);
		create unique index if not exists users_username_idx on users (username);
		create unique index if not exists users_email_idx on users (email);
		create index if not exists users_created_at_idx on users (created_at desc);

		create table if not exists roles (
			user_id varchar,
			admin bool not null default false,
			instructor bool not null default false,
			created_at timestamp not null default now(),
			updated_at timestamp not null default now(),
			primary key (user_id),
			foreign key (user_id) references users (id)
		);
		create index if not exists roles_admin_idx on roles (admin);
		create index if not exists roles_instructor_idx on roles (instructor);
	`)
	if err != nil {
		return nil, err
	}

	return r, nil
}

type userRepository struct {
	db *sql.DB
}

func (r *userRepository) Store(ctx context.Context, user *user.User) error {
	_, err := r.db.ExecContext(ctx, `
		insert into users
			(
				id,
				name,
				username,
				email,
				about_me,
				image,
				updated_at
			)
		values
			(
				id = @id,
				name = @name,
				username = @username,
				email = @email,
				about_me = @about_me,
				image = @image,
				updated_at = now()
			)
		on conflict do update set
			id = excluded.id,
			name = excluded.name,
			username = excluded.username,
			email = excluded.email,
			about_me = excluded.about_me,
			image = excluded.image,
			updated_at = excluded.updated_at,
	`,
		sql.Named("id", user.ID),
		sql.Named("name", user.Name),
		sql.Named("username", user.Username),
		sql.Named("email", sql.NullString{String: user.Email, Valid: len(user.Email) > 0}),
		sql.Named("about_me", user.ID),
		sql.Named("image", user.ID),
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) FindID(ctx context.Context, id string) (*user.User, error) {
	var x user.User
	err := r.db.QueryRowContext(ctx, `
		select
			users.id,
			users.name,
			users.username,
			coalesce(users.email, ''),
			users.about_me,
			users.image,
			users.created_at,
			users.updated_at,
			coalesce(roles.admin, false),
			coalesce(roles.instructor, false)
		from users
			left join roles on users.id = roles.user_id
		where users.id = $1
	`,
		id, // 1
	).Scan(
		&x.ID,
		&x.Name,
		&x.Username,
		&x.Email,
		&x.AboutMe,
		&x.Image,
		&x.CreatedAt,
		&x.UpdatedAt,
		&x.Role.Admin,
		&x.Role.Instructor,
	)
	if err == sql.ErrNoRows {
		return nil, user.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &x, nil
}

func (r *userRepository) FindEmail(ctx context.Context, email string) (*user.User, error) {
	var x user.User
	err := r.db.QueryRowContext(ctx, `
		select
			users.id,
			users.name,
			users.username,
			coalesce(users.email, ''),
			users.about_me,
			users.image,
			users.created_at,
			users.updated_at,
			coalesce(roles.admin, false),
			coalesce(roles.instructor, false)
		from users
		where users.email = @email
	`,
		sql.Named("email", email),
	).Scan(
		&x.ID,
		&x.Name,
		&x.Username,
		&x.Email,
		&x.AboutMe,
		&x.Image,
		&x.CreatedAt,
		&x.UpdatedAt,
		&x.Role.Admin,
		&x.Role.Instructor,
	)
	if err == sql.ErrNoRows {
		return nil, user.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &x, nil
}

func (r *userRepository) FindUsername(ctx context.Context, username string) (*user.User, error) {
	var x user.User
	err := r.db.QueryRowContext(ctx, `
		select
			users.id,
			users.name,
			users.username,
			coalesce(users.email, ''),
			users.about_me,
			users.image,
			users.created_at,
			users.updated_at,
			coalesce(roles.admin, false),
			coalesce(roles.instructor, false)
		from users
		where users.username = @username
	`,
		sql.Named("username", username),
	).Scan(
		&x.ID,
		&x.Name,
		&x.Username,
		&x.Email,
		&x.AboutMe,
		&x.Image,
		&x.CreatedAt,
		&x.UpdatedAt,
		&x.Role.Admin,
		&x.Role.Instructor,
	)
	if err == sql.ErrNoRows {
		return nil, user.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &x, nil
}

func (r *userRepository) List(ctx context.Context, limit, offset int64) ([]*user.User, error) {
	xs := make([]*user.User, 0)
	rows, err := r.db.QueryContext(ctx, `
		select
			users.id,
			users.name,
			users.username,
			coalesce(users.email, ''),
			users.about_me,
			users.image,
			users.created_at,
			users.updated_at,
			coalesce(roles.admin, false),
			coalesce(roles.instructor, false)
		from users
			left join roles on users.id = roles.user_id
		order by users.created_at desc
		limit @limit
		offset @offset
	`,
		sql.Named("limit", limit),
		sql.Named("offset", offset),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var x user.User
		err = rows.Scan(
			&x.ID,
			&x.Name,
			&x.Username,
			&x.Email,
			&x.AboutMe,
			&x.Image,
			&x.CreatedAt,
			&x.UpdatedAt,
			&x.Role.Admin,
			&x.Role.Instructor,
		)
		if err != nil {
			return nil, err
		}
		xs = append(xs, &x)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return xs, nil
}

func (r *userRepository) Count(ctx context.Context) (int64, error) {
	var cnt int64
	err := r.db.QueryRowContext(ctx, `select count(*) from users`).Scan(&cnt)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}
