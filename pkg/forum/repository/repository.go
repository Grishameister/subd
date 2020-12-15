package repository

import (
	"context"
	"errors"
	"github.com/Grishameister/subd/configs/config"
	"github.com/Grishameister/subd/internal/database"
	"github.com/Grishameister/subd/pkg/domain"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx"
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

func (r *Repo) CreateForum(f *domain.Forum) (domain.Forum, error) {
	if err := r.db.QueryRow(context.Background(), "insert into forums (slug, owner, title) "+
		"values ($1, (select nickname from users where nickname = $2), $3) returning owner", f.Slug, f.Owner, f.Title).Scan(&f.Owner); err != nil {
		config.Lg("forum", "CreateForum").Error(err.Error())
		if pqErr, ok := err.(*pgconn.PgError); ok {
			if pqErr.Code == "23505" {
				fr, _ := r.GetForum(f.Slug)
				return fr, errors.New("slug exists")
			}
			if pqErr.Code == "23503" || pqErr.Code == "23502" {
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
	if err := r.db.QueryRow(context.Background(), "select slug, owner, title, threads, "+
		"posts from forums where slug=$1", f.Slug).Scan(&f.Slug, &f.Owner, &f.Title, &f.Threads, &f.Posts); err != nil {
		config.Lg("forum", "GetForum").Error(err.Error())
		if err.Error() == pgx.ErrNoRows.Error() {
			return domain.Forum{}, errors.New("forum not found")
		}
		return domain.Forum{}, err
	}
	return f, nil
}

func (r *Repo) GetThreads(slug string, limit string, since string, order string) ([]domain.Thread, error) {
	threads := make([]domain.Thread, 0)

	query := "select id, author, message, forum, title, created," +
		"case when slug is null then '' else slug end, votes from threads where forum = $1"

	i := 1
	var values []interface{}
	values = append(values, slug)
	if order == "asc" {
		if since != "" {
			i++
			query += " and created >= $" + strconv.Itoa(i)
			values = append(values, since)
		}
	} else {
		if since != "" {
			i++
			query += " and created <= $" + strconv.Itoa(i)
			values = append(values, since)
		}
	}

	query += " order by created " + order
	i++
	query += " limit $" + strconv.Itoa(i)
	values = append(values, limit)

	rows, err := r.db.Query(context.Background(), query, values...)
	if err != nil {
		config.Lg("thread", "GetThreadSlugOrId").Error(err.Error())
		return nil, errors.New("forum not found")
	}
	defer rows.Close()

	for rows.Next() {
		t := domain.Thread{}
		if err := rows.Scan(&t.Id, &t.Author, &t.Message, &t.Forum, &t.Title, &t.Created, &t.Slug, &t.Votes); err != nil {
			config.Lg("forum", "GetThread").Error(err.Error())
			return nil, err
		}
		threads = append(threads, t)
	}
	if len(threads) == 0 {
		if _, err := r.GetForum(slug); err != nil {
			config.Lg("thread", "GetThreadSlugOrId").Error(err.Error())
			return nil, errors.New("forum not found")
		}
	}

	return threads, nil
}

func (r *Repo) GetUsers(slug string, limit string, since string, order string) ([]domain.User, error) {
	users := make([]domain.User, 0)

	query := "select about, email, fullname, forums_users.nickname from users join forums_users on forums_users.nickname = users.nickname " +
		"where forums_users.forum_slug = $1"

	i := 1
	var values []interface{}
	values = append(values, slug)
	if order == "asc" {
		if since != "" {
			i++
			query += " and forums_users.nickname > $" + strconv.Itoa(i)
			values = append(values, since)
		}
	} else {
		if since != "" {
			i++
			query += " and forums_users.nickname < $" + strconv.Itoa(i)
			values = append(values, since)
		}
	}

	query += " order by forums_users.nickname " + order
	i++
	query += " limit $" + strconv.Itoa(i)
	values = append(values, limit)

	rows, err := r.db.Query(context.Background(), query, values...)
	if err != nil {
		config.Lg("forum", "GetUsers").Error(err.Error())
		return nil, errors.New("forum not found")
	}
	defer rows.Close()

	for rows.Next() {
		u := domain.User{}
		if err := rows.Scan(&u.About, &u.Email, &u.Fullname, &u.Nickname); err != nil {
			config.Lg("forum", "GetUsers").Error(err.Error())
			return nil, err
		}
		log.Println(u.Nickname)
		users = append(users, u)
	}

	if len(users) == 0 {
		_, err := r.GetForum(slug)
		if err != nil {
			return nil, errors.New("forum not found")
		}
	}

	return users, nil
}
