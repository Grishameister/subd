package thread

import "github.com/Grishameister/subd/pkg/domain"

type IRepo interface {
	CreateInForum(t *domain.Thread) (domain.Thread, error)
	GetThreadBySlugOrId(slugOrId string) (domain.Thread, error)
}
