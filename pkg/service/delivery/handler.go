package delivery

import (
	"github.com/Grishameister/subd/pkg/service"
	"github.com/Grishameister/subd/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	uc service.IUsecase
}

func New(uc service.IUsecase) *Handler {
	return &Handler{
		uc: uc,
	}
}

func (h *Handler) Clear(c *gin.Context) {
	if err := h.uc.Clear(); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.Error{Error: err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) Status(c *gin.Context) {
	s, err := h.uc.Status()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}
