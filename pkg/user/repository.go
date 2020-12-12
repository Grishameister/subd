package user

import (
	"github.com/Grishameister/subd/pkg/domain"
)

type IRepo interface {
	CreateUser(u *domain.User) ([]domain.User, error)
	GetUser(nickname string) (domain.User, error)
	UpdateUser(u *domain.User) (domain.User, error)
}
