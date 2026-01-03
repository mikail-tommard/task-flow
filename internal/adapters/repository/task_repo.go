package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/mikail-tommard/task-flow/internal/domain"
)

var ErrNotFound = errors.New("not found")

type Repo struct {
	db *sql.DB
}

func New(db *sql.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) Create(ctx context.Context, task *domain.Task) (int, error) {
	const q = `
		INSERT INTO tasks (user_id, title, description, done)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var id int
	err := r.db.QueryRowContext(ctx, q, task.UserID(), task.Title(), task.Description(), task.Done()).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repo) GetByID(ctx context.Context, id int) (*domain.Task, error) {
	const q = `
		SELECT id, user_id, title, description, done
		FROM tasks
		WHERE id = $1
	`

	var (
		dbID   int
		userID int
		title  string
		desc   string
		done   bool
	)

	err := r.db.QueryRowContext(ctx, q, id).Scan(&dbID, &userID, &title, &desc, &done)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return domain.FromStorage(dbID, title, done, desc, userID), nil
}

func (r *Repo) ListByUser(ctx context.Context, userID int) ([]*domain.Task, error) {
	const q = `
		SELECT id, user_id, title, description, done
		FROM tasks
		WHERE user_id = $1
	`

	rows, err := r.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*domain.Task, 0)
	for rows.Next() {
		var (
			dbID  int
			uID   int
			title string
			desc  string
			done  bool
		)

		if err := rows.Scan(&dbID, &uID, &title, &desc, &done); err != nil {
			return nil, err
		}
		tasks = append(tasks, domain.FromStorage(dbID, title, done, desc, uID))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *Repo) Update(ctx context.Context, task *domain.Task) error {
	const q = `
		UPDATE tasks
		SET title = $1, description = $2, done = $3
		WHERE id = $4
	`

	res, err := r.db.ExecContext(ctx, q, task.Title(), task.Description(), task.Done(), task.ID())
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err == nil && n == 0 {
		return ErrNotFound
	}
	return nil
}
