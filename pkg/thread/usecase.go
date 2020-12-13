package thread

import "github.com/Grishameister/subd/pkg/domain"

type IUseCase interface {
	CreateInForum(t *domain.Thread) (domain.Thread, error)
	GetThreadBySlugOrId(slugOrId string) (domain.Thread, error)
	UpdateThread(slugOrId string, t *domain.ThreadUpdate) (domain.Thread, error)
}
