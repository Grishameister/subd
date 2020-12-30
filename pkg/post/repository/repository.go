package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/Grishameister/subd/configs/config"
	"github.com/Grishameister/subd/internal/database"
	"github.com/Grishameister/subd/pkg/domain"
	"github.com/Grishameister/subd/pkg/forum"
	"github.com/Grishameister/subd/pkg/thread"
	"github.com/Grishameister/subd/pkg/user"
	"strconv"
	"strings"
)

type Repo struct {
	db database.IDbConn
	tr thread.IRepo
	fr forum.IRepo
	ur user.IRepo
}

func New(db database.IDbConn, tr thread.IRepo, fr forum.IRepo, ur user.IRepo) *Repo {
	return &Repo{
		db: db,
		tr: tr,
		fr: fr,
		ur: ur,
	}
}

func (r *Repo) checkParent(idThread int, id int64) error {
	err := r.db.QueryRow(context.Background(),
		"select id from posts where thread = $1 and id = $2", idThread, id).Scan(&id)
	if err != nil {
		return errors.New("post not found")
	}
	return nil
}

func (r *Repo) checkAuthor(author string) error {
	err := r.db.QueryRow(context.Background(),
		`select nickname from users where nickname = $1`, author).Scan(&author)

	if err != nil {
		return errors.New("user not found")
	}
	return nil
}

func (r *Repo) CreatePosts(slugOrId string, posts []*domain.Post) ([]*domain.Post, error) {
	t, err := r.tr.GetThreadBySlugOrId(slugOrId)
	if err != nil {
		return nil, errors.New("thread not found")
	}
	var b strings.Builder

	if len(posts) == 0 {
		return posts, nil
	}

	query := "insert into posts(author, forum, message, parent, thread, post_path) values "

	b.WriteString(query)
	values := make([]interface{}, 0, len(posts)*5)

	for i, post := range posts {
		post.Thread = int64(t.Id)
		post.ForumSlug = t.Forum

		if err := r.checkParent(int(post.Thread), int64(post.Parent)); err != nil && post.Parent != 0 {
			return nil, err
		}

		if err := r.checkAuthor(post.Author); err != nil {
			return nil, err
		}

		offset := i * 5
		b.WriteString(fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, "+
			"array_append((SELECT post_path FROM posts WHERE id = $%d), (SELECT last_value FROM posts_id_seq)))",
			offset+1, offset+2, offset+3, offset+4, offset+5, offset+4))
		if i != len(posts)-1 {
			b.WriteString(",")
		}
		values = append(values, post.Author, post.ForumSlug, post.Message, post.Parent, post.Thread)
	}

	b.WriteString(" returning id, created")

	rows, err := r.db.Query(context.Background(), b.String(), values...)

	if err != nil {
		config.Lg("post", "CreatePosts").Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	j := 0
	for rows.Next() {
		if err := rows.Scan(&posts[j].Id, &posts[j].Created); err != nil {
			config.Lg("posts", "CreatePosts").Error(err.Error())
			return nil, err
		}
		j++
	}
	return posts, nil
}

func (r *Repo) GetPosts(slugOrId string, limit string, since string, sort string, order string) ([]domain.Post, error) {
	t, err := r.tr.GetThreadBySlugOrId(slugOrId)
	if err != nil {
		return nil, errors.New("thread not found")
	}
	var b strings.Builder

	b.WriteString("select author, created, forum, message, parent, thread, id, isEdit from posts where thread=$1 ")

	i := 1
	var values []interface{}
	values = append(values, t.Id)

	switch sort {
	case "", "flat":
		if order == "asc" {
			if since != "" {
				i++
				b.WriteString(" and id > $" + strconv.Itoa(i))
				values = append(values, since)
			}
		} else {
			if since != "" {
				i++
				b.WriteString(" and id < $" + strconv.Itoa(i))
				values = append(values, since)
			}
		}
		b.WriteString(" order by id " + order)
		i++
		b.WriteString(" limit $" + strconv.Itoa(i))
		values = append(values, limit)
	case "tree":
		if order == "asc" {
			if since != "" {
				i++
				b.WriteString(" and post_path > (select post_path FROM posts WHERE id = $" + strconv.Itoa(i))
				b.WriteString(")")
				values = append(values, since)
			}
		} else {
			if since != "" {
				i++
				b.WriteString(" and post_path < (select post_path from posts where id=$" + strconv.Itoa(i))
				b.WriteString(")")
				values = append(values, since)
			}
		}
		b.WriteString(" order by post_path " + order)
		i++
		b.WriteString(" limit $" + strconv.Itoa(i))
		values = append(values, limit)
	case "parent_tree":
		b.WriteString(" and post_path[1] IN (SELECT post_path[1] from posts where thread = $1 and array_length(post_path, 1) = 1")
		if order == "asc" {
			if since != "" {
				i++
				b.WriteString(" and post_path[1] > (select post_path[1] from posts WHERE id=$" + strconv.Itoa(i))
				b.WriteString(")")
				values = append(values, since)
			}
		} else {
			if since != "" {
				i++
				b.WriteString(" and post_path[1] < (select post_path[1] from posts where id=$" + strconv.Itoa(i))
				b.WriteString(")")
				values = append(values, since)
			}
		}
		b.WriteString(" order by post_path " + order)
		i++
		b.WriteString(" limit $" + strconv.Itoa(i))
		values = append(values, limit)
		b.WriteString(") order by post_path[1] " + order)
		b.WriteString(", post_path[2:]")
	}

	rows, err := r.db.Query(context.Background(), b.String(), values...)

	if err != nil {
		config.Lg("post", "GetPosts").Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	posts := make([]domain.Post, 0)
	for rows.Next() {
		p := domain.Post{
			Parent: 0,
		}
		if err := rows.Scan(&p.Author, &p.Created, &p.ForumSlug, &p.Message, &p.Parent, &p.Thread, &p.Id, &p.IsEdited); err != nil {
			config.Lg("post", "GetPosts").Error(err.Error())
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *Repo) GetPost(id string, related string) (domain.PostFull, error) {
	p := domain.PostFull{}

	if err := r.db.QueryRow(context.Background(), "select "+
		"author, created, forum, message, parent, thread, id, isEdit from posts where id=$1", id).
		Scan(&p.Post.Author, &p.Post.Created, &p.Post.ForumSlug, &p.Post.Message,
			&p.Post.Parent, &p.Post.Thread, &p.Post.Id, &p.Post.IsEdited); err != nil {
		config.Lg("post", "GetPost").Error(err.Error())
		return p, errors.New("post not found")
	}

	if strings.Contains(related, "thread") {
		tr, _ := r.tr.GetThreadBySlugOrId(strconv.FormatInt(p.Post.Thread, 10))
		p.Thread = &tr
	}

	if strings.Contains(related, "forum") {
		fr, _ := r.fr.GetForum(p.Post.ForumSlug)
		p.Forum = &fr
	}
	if strings.Contains(related, "user") {
		ur, _ := r.ur.GetUser(p.Post.Author)
		p.Profile = &ur
	}
	return p, nil
}

func (r *Repo) UpdatePost(id string, message string) (domain.Post, error) {
	var p domain.Post

	if message == "" {
		if err := r.db.QueryRow(context.Background(), "select "+
			"author, created, forum, message, parent, thread, id, isEdit from posts where id=$1", id).
			Scan(&p.Author, &p.Created, &p.ForumSlug, &p.Message,
				&p.Parent, &p.Thread, &p.Id, &p.IsEdited); err != nil {
			config.Lg("post", "UpdatePost").Error(err.Error())
			return p, err
		}
		return p, nil
	}

	if err := r.db.QueryRow(context.Background(), "update posts set isEdit = "+
		"(case when message = $1 then false else true end), message = $1 where id = $2 "+
		"returning author, created, forum, message, parent, thread, id, isEdit", message, id).
		Scan(&p.Author, &p.Created, &p.ForumSlug, &p.Message,
			&p.Parent, &p.Thread, &p.Id, &p.IsEdited); err != nil {
		config.Lg("post", "UpdatePost").Error(err.Error())
		return p, errors.New("post not found")
	}
	return p, nil
}
