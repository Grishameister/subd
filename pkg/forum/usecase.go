package forum

import "github.com/Grishameister/subd/pkg/domain"

type IUseCase interface {
	CreateForum(f *domain.Forum) (domain.Forum, error)
	GetForum(slug string) (domain.Forum, error)
}
