package usecase

import (
	"github.com/Grishameister/subd/pkg/domain"
	"github.com/Grishameister/subd/pkg/service"
)

type Usecase struct {
	r service.IRepo
}

func New(r service.IRepo) *Usecase {
	return &Usecase{
		r: r,
	}
}

func (uc *Usecase) Clear() error {
	return uc.r.Clear()
}

func (uc *Usecase) Status() (domain.Status, error) {
	return uc.r.Status()
}
