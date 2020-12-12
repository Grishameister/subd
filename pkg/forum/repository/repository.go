package repository

import (
	"context"
	"errors"
	"github.com/Grishameister/subd/configs/config"
	"github.com/Grishameister/subd/internal/database"
	"github.com/Grishameister/subd/pkg/domain"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx"
)

type Repo struct {
	db database.IDbConn
}

func New(db database.IDbConn) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) CreateForum(f *domain.Forum) (domain.Forum, error) {
	var owner string
	if err := r.db.QueryRow(context.Background(), "insert into forums (slug, owner, title) "+
		"values ($1, $2, $3) returning owner", f.Slug, f.Owner, f.Title).Scan(&owner); err != nil {
		config.Lg("forum", "CreateForum").Error(err.Error())
		if pqErr, ok := err.(*pgconn.PgError); ok {
			if pqErr.Code == "23505" {
				fr, _ := r.GetForum(f.Slug)
				return fr, errors.New("slug exists")
			}
			if pqErr.Code == "23503" {
				return domain.Forum{}, errors.New("user not found")
			}
		}
		return domain.Forum{}, err
	}
	return *f, nil
}

func (r *Repo) GetForum(slug string) (domain.Forum, error) {
	f := domain.Forum{}
	f.Slug = slug
	if err := r.db.QueryRow(context.Background(), "select owner, title, threads, "+
		"posts from forums where slug=$1", f.Slug).Scan(&f.Owner, &f.Title, &f.Threads, &f.Posts); err != nil {
		config.Lg("forum", "GetForum").Error(err.Error())
		if err.Error() == pgx.ErrNoRows.Error() {
			return domain.Forum{}, errors.New("forum not found")
		}
		return domain.Forum{}, err
	}
	return f, nil
}
