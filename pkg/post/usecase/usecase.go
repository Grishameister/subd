package usecase

import (
	"github.com/Grishameister/subd/pkg/domain"
	"github.com/Grishameister/subd/pkg/post"
)

type Usecase struct {
	r post.IRepo
}

func New(r post.IRepo) *Usecase {
	return &Usecase{
		r: r,
	}
}

func (uc *Usecase) CreatePosts(slugOrId string, posts []*domain.Post) ([]*domain.Post, error) {
	return uc.r.CreatePosts(slugOrId, posts)
}

func (uc *Usecase) GetPosts(slugOrId string, limit string, since string, sort string, order string) ([]domain.Post, error) {
	return uc.r.GetPosts(slugOrId, limit, since, sort, order)
}

func (uc *Usecase) GetPost(id string, related string) (domain.PostFull, error) {
	return uc.r.GetPost(id, related)
}

func (uc *Usecase) UpdatePost(id string, message string) (domain.Post, error) {
	return uc.r.UpdatePost(id, message)
}
