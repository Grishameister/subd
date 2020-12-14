package post

import "github.com/Grishameister/subd/pkg/domain"

type IRepo interface {
	CreatePosts(slugOrId string, posts []*domain.Post) ([]*domain.Post, error)
	GetPosts(slugOrId string, limit string, since string, sort string, order string) ([]domain.Post, error)
	GetPost(id string, related string) (domain.PostFull, error)
	UpdatePost(id string, message string) (domain.Post, error)
}
