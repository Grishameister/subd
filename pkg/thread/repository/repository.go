package repository

import (
	"context"
	"errors"
	"github.com/Grishameister/subd/configs/config"
	"github.com/Grishameister/subd/internal/database"
	"github.com/Grishameister/subd/pkg/domain"
	"github.com/jackc/pgconn"
	"strconv"
	"strings"
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
	end := " values((select slug from forums where slug = $1), $2, $3, $4"
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
	end += ") returning id, created, forum"

	if err := r.db.QueryRow(context.Background(), start+end, placeholders...).Scan(&t.Id, &t.Created, &t.Forum); err != nil {
		config.Lg("forum", "CreateForum").Error(err.Error())
		if pqErr, ok := err.(*pgconn.PgError); ok {
			if pqErr.Code == "23505" {
				tr, _ := r.GetThreadBySlugOrId(t.Slug)
				return tr, errors.New("thread exists")
			}
			if pqErr.Code == "23503" || pqErr.Code == "23502" {
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

func (r *Repo) UpdateThread(slugOrId string, t *domain.ThreadUpdate) (domain.Thread, error) {
	var where string
	var tr domain.Thread

	query := "update threads set "
	var queryParams []string
	var values []interface{}
	i := 0
	if t.Message != "" {
		i++
		queryParams = append(queryParams, "message=$"+strconv.Itoa(i))
		values = append(values, t.Message)
	}
	if t.Title != "" {
		i++
		queryParams = append(queryParams, "title=$"+strconv.Itoa(i))
		values = append(values, t.Title)
	}
	if i == 0 {
		tr, err := r.GetThreadBySlugOrId(slugOrId)
		if err != nil {
			return domain.Thread{}, err
		}
		return tr, nil
	}

	query += strings.Join(queryParams, ",")

	i++
	values = append(values, slugOrId)

	if _, err := strconv.Atoi(slugOrId); err != nil {
		where = " where slug =$" + strconv.Itoa(i)
	} else {
		where = " where id =$" + strconv.Itoa(i)
	}
	query += where

	query += "returning id, author, forum, created, case when slug is null then '' else slug end, title, message, votes"

	if err := r.db.QueryRow(context.Background(), query, values...).
		Scan(&tr.Id, &tr.Author, &tr.Forum, &tr.Created, &tr.Slug, &tr.Title, &tr.Message, &tr.Votes); err != nil {
		config.Lg("thread", "UpdateThread").Error(err.Error())
		return tr, errors.New("not found thread")
	}

	return tr, nil
}

func (r *Repo) VoteThread(slugOrId string, v *domain.Vote) (domain.Thread, error) {
	t, err := r.GetThreadBySlugOrId(slugOrId)

	if err != nil {
		return t, errors.New("thread not found")
	}

	if _, err := r.db.Exec(context.Background(), "insert into votes (author, thread, vote) values ($1, $2, $3) "+
		"on conflict (thread, author) do update set vote = $3", v.ProfileNickname, t.Id, v.Voice); err != nil {
		config.Lg("thread", "VoteThread").Error(err.Error())
		if pqErr, ok := err.(*pgconn.PgError); ok {
			if pqErr.Code == "23505" {
				tr, _ := r.GetThreadBySlugOrId(t.Slug)
				return tr, errors.New("thread exists")
			}
			if pqErr.Code == "23503" || pqErr.Code == "23502" {
				return t, errors.New("user not found")
			}
		}
		return t, err
	}

	t, err = r.GetThreadBySlugOrId(slugOrId)

	if err != nil {
		return t, errors.New("thread not found")
	}

	return t, nil
}
