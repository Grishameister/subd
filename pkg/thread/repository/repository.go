package repository

import (
	"context"
	"errors"
	"github.com/Grishameister/subd/configs/config"
	"github.com/Grishameister/subd/internal/database"
	"github.com/Grishameister/subd/pkg/domain"
	"github.com/jackc/pgconn"
	"log"
	"strconv"
)

type Repo struct {
	db database.IDbConn
}

func New(db database.IDbConn) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) CreateInForum(t *domain.Thread) (domain.Thread, error) {
	start := "insert into threads (forum, author, message, title"
	end := " values($1, $2, $3, $4"
	var placeholders []interface{}
	i := 4
	placeholders = append(placeholders, t.Forum, t.Author, t.Message, t.Title)
	if t.Slug != "" {
		i++
		start += ", slug"
		end += ", $5"
		placeholders = append(placeholders, t.Slug)
	}

	if !t.Created.IsZero() {
		i++
		start += ", created"
		end += ", $" + strconv.Itoa(i)
		placeholders = append(placeholders, t.Created)
	}
	start += ") "
	end += ") returning id, created"

	log.Println(start + end)

	if err := r.db.QueryRow(context.Background(), start+end, placeholders...).Scan(&t.Id, &t.Created); err != nil {
		config.Lg("forum", "CreateForum").Error(err.Error())
		if pqErr, ok := err.(*pgconn.PgError); ok {
			if pqErr.Code == "23505" {
				tr, _ := r.GetThreadBySlugOrId(t.Slug)
				return tr, errors.New("thread exists")
			}
			if pqErr.Code == "23503" {
				return *t, errors.New("user or forum not found")
			}
		}
		return domain.Thread{}, err
	}
	return *t, nil
}

func (r *Repo) GetThreadBySlugOrId(slugOrId string) (domain.Thread, error) {
	var where string
	var t domain.Thread
	if _, err := strconv.Atoi(slugOrId); err != nil {
		where = " slug = $1"
	} else {
		where = " id = $1"
	}

	if err := r.db.QueryRow(context.Background(), "select id, author, message, forum, title, created,"+
		"case when slug is null then '' else slug end, votes from threads where"+where,
		slugOrId).Scan(&t.Id, &t.Author, &t.Message, &t.Forum, &t.Title, &t.Created, &t.Slug, &t.Votes); err != nil {
		config.Lg("thread", "GetThreadSlugOrId").Error(err.Error())
		return t, err
	}
	return t, nil
}
