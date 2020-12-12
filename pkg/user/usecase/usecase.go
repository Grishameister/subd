package usecase

import (
	"github.com/Grishameister/subd/pkg/domain"
	"github.com/Grishameister/subd/pkg/user"
)

type UserUsecase struct {
	r user.IRepo
}

func New(r user.IRepo) *UserUsecase {
	return &UserUsecase{
		r: r,
	}
}

func (uc *UserUsecase) CreateUser(u *domain.User) ([]domain.User, error) {
	return uc.r.CreateUser(u)
}

func (uc *UserUsecase) GetUser(nickname string) (domain.User, error) {
	return uc.r.GetUser(nickname)
}

func (uc *UserUsecase) UpdateUser(u *domain.User) (domain.User, error) {
	return uc.r.UpdateUser(u)
}
