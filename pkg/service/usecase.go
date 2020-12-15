package service

import "github.com/Grishameister/subd/pkg/domain"

type IUsecase interface {
	Clear() error
	Status() (domain.Status, error)
}
