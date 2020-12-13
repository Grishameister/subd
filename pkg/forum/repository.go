package forum

import "github.com/Grishameister/subd/pkg/domain"

type IRepo interface {
	CreateForum(f *domain.Forum) (domain.Forum, error)
	GetForum(slug string) (domain.Forum, error)
	GetThreads(slug string, limit string, since string, order string) ([]domain.Thread, error)
	GetUsers(slug string, limit string, since string, order string) ([]domain.User, error)
}
