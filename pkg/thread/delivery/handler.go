package delivery

import (
	"github.com/Grishameister/subd/pkg/domain"
	"github.com/Grishameister/subd/pkg/thread"
	"github.com/Grishameister/subd/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	uc thread.IUseCase
}

func New(uc thread.IUseCase) *Handler {
	return &Handler{
		uc: uc,
	}
}

func (h *Handler) CreateThread(c *gin.Context) {
	forum := c.Param("slug")
	var t domain.Thread
	if err := c.BindJSON(&t); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	t.Forum = forum

	tr, err := h.uc.CreateInForum(&t)
	if err != nil {
		if err.Error() == "thread exists" {
			c.AbortWithStatusJSON(http.StatusConflict, tr)
			return
		}
		if err.Error() == "user or forum not found" {
			c.AbortWithStatusJSON(http.StatusNotFound, utils.Error{Error: "not found user or forum"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, tr)
}

func (h *Handler) GetThread(c *gin.Context) {
	slugOrId := c.Param("slug_or_id")

	tr, err := h.uc.GetThreadBySlugOrId(slugOrId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, utils.Error{Error: "not found" + slugOrId})
		return
	}
	c.JSON(http.StatusOK, tr)
}
