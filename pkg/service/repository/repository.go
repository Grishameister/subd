package repository

import (
	"context"
	"github.com/Grishameister/subd/internal/database"
	"github.com/Grishameister/subd/pkg/domain"
)

type Repo struct {
	db database.IDbConn
}

func New(db database.IDbConn) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) Clear() error {
	if _, err := r.db.Exec(context.Background(), "truncate table votes, posts, threads, forums_users, forums, users"); err != nil {
		return err
	}
	return nil
}

func (r *Repo) Status() (domain.Status, error) {
	s := domain.Status{}

	if err := r.db.QueryRow(context.Background(), "select count(nickname) from users").Scan(&s.User); err != nil {
		return s, err
	}

	if err := r.db.QueryRow(context.Background(), "select count(slug) from forums").Scan(&s.Forum); err != nil {
		return s, err
	}
	if err := r.db.QueryRow(context.Background(), "select count(id) from threads").Scan(&s.Thread); err != nil {
		return s, err
	}
	if err := r.db.QueryRow(context.Background(), "select count(id) from posts").Scan(&s.Post); err != nil {
		return s, err
	}
	return s, nil
}
