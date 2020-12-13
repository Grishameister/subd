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

func (uc *UseCase) GetThreads(slug string, limit string, since string, order string) ([]domain.Thread, error) {
	return uc.r.GetThreads(slug, limit, since, order)
}

func (uc *UseCase) GetUsers(slug string, limit string, since string, order string) ([]domain.User, error) {
	return uc.r.GetUsers(slug, limit, since, order)
}
