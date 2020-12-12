package usecase

import (
	"github.com/Grishameister/subd/pkg/domain"
	"github.com/Grishameister/subd/pkg/forum"
)

type UseCase struct {
	r forum.IRepo
}

func New(r forum.IRepo) *UseCase {
	return &UseCase{
		r: r,
	}
}

func (uc *UseCase) CreateForum(f *domain.Forum) (domain.Forum, error) {
	return uc.r.CreateForum(f)
}

func (uc *UseCase) GetForum(slug string) (domain.Forum, error) {
	return uc.r.GetForum(slug)
}
