package repository

import (
	"context"
	"errors"
	"github.com/Grishameister/subd/configs/config"
	"github.com/Grishameister/subd/internal/database"
	"github.com/Grishameister/subd/pkg/domain"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx"
	"strconv"
	"strings"
)

type UserRepo struct {
	db database.IDbConn
}

func New(db database.IDbConn) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(u *domain.User) ([]domain.User, error) {
	var users []domain.User
	_, err := r.db.Exec(context.Background(), "insert into users (nickname, about, fullname, email) "+
		"values ($1, $2, $3,$4)", u.Nickname, u.About, u.Fullname, u.Email)
	if err != nil {
		config.Lg("user", "CreateUser").Error(err.Error())

		if pqErr, ok := err.(*pgconn.PgError); ok {
			if pqErr.Code == "23505" {
				rows, err := r.db.Query(context.Background(), "select nickname, about, fullname, email "+
					"from users where nickname = $1 or email = $2", u.Nickname, u.Email)
				if err != nil {
					config.Lg("user", "CreateSUsers").Error(err.Error())
					return nil, err
				}
				defer rows.Close()
				for rows.Next() {
					us := domain.User{}
					if err := rows.Scan(&us.Nickname, &us.About, &us.Fullname, &us.Email); err != nil {
						return nil, err
					}
					users = append(users, us)
				}
				return users, errors.New("already exists")
			}
		}
		return nil, err
	}
	return []domain.User{*u}, nil
}

func (r *UserRepo) GetUser(nickname string) (domain.User, error) {
	u := domain.User{}
	if err := r.db.QueryRow(context.Background(), "select nickname, about, fullname, "+
		"email from users where nickname = $1", nickname).Scan(&u.Nickname, &u.About, &u.Fullname, &u.Email); err != nil {
		return u, err
	}
	return u, nil
}

func (r *UserRepo) UpdateUser(u *domain.User) (domain.User, error) {
	i := 0
	query := "update users set "
	var queryParams []string
	var values []interface{}
	if u.Email != "" {
		i++
		queryParams = append(queryParams, "email=$"+strconv.Itoa(i))
		values = append(values, u.Email)
	}
	if u.About != "" {
		i++
		queryParams = append(queryParams, "about=$"+strconv.Itoa(i))
		values = append(values, u.About)
	}

	if u.Fullname != "" {
		i++
		queryParams = append(queryParams, "fullname=$"+strconv.Itoa(i))
		values = append(values, u.Fullname)
	}

	if i == 0 {
		u, err := r.GetUser(u.Nickname)
		if err != nil {
			return domain.User{}, err
		}
		return u, nil
	}

	query += strings.Join(queryParams, ",")
	i++
	query += " where nickname=$" + strconv.Itoa(i) + " returning nickname, email, fullname, about"

	values = append(values, u.Nickname)
	var nick string

	if err := r.db.QueryRow(context.Background(), query,
		values...).Scan(&nick, &u.Email, &u.Fullname, &u.About); err != nil {
		config.Lg("user", "UpdateUser").Error("UpdateUser: ", err.Error())
		if pgx.ErrNoRows.Error() == err.Error() {
			return *u, errors.New("not found user")
		}
		if err, ok := err.(*pgconn.PgError); ok {
			if err.Code == "23505" {
				return *u, errors.New("conflict")
			}
		}
		return *u, err
	}
	return *u, nil
}
