package usecase

import (
	"github.com/Grishameister/subd/pkg/domain"
	"github.com/Grishameister/subd/pkg/thread"
)

type Usecase struct {
	r thread.IRepo
}

func New(r thread.IRepo) *Usecase {
	return &Usecase{
		r: r,
	}
}

func (uc *Usecase) CreateInForum(t *domain.Thread) (domain.Thread, error) {
	return uc.r.CreateInForum(t)
}
func (uc *Usecase) GetThreadBySlugOrId(slugOrId string) (domain.Thread, error) {
	return uc.r.GetThreadBySlugOrId(slugOrId)
}
