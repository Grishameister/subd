package service

import "github.com/Grishameister/subd/pkg/domain"

type IRepo interface {
	Clear() error
	Status() (domain.Status, error)
}
