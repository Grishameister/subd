package delivery

import (
	"github.com/Grishameister/subd/pkg/domain"
	"github.com/Grishameister/subd/pkg/user"
	"github.com/Grishameister/subd/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	uc user.IUseCase
}

func New(uc user.IUseCase) *Handler {
	return &Handler{
		uc: uc,
	}
}

func (h *Handler) CreateUser(c *gin.Context) {
	nick := c.Param("nickname")

	u := domain.User{}
	if err := c.BindJSON(&u); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	u.Nickname = nick

	users, err := h.uc.CreateUser(&u)
	if err != nil {
		if err.Error() == "already exists" {
			c.AbortWithStatusJSON(http.StatusConflict, users)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusCreated, users[0])
}

func (h *Handler) GetUser(c *gin.Context) {
	nick := c.Param("nickname")

	u, err := h.uc.GetUser(nick)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, utils.Error{Error: "Can't find user with id " + nick})
		return
	}
	c.JSON(http.StatusOK, u)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	nick := c.Param("nickname")

	u := domain.User{}
	if err := c.BindJSON(&u); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.Error{Error: err.Error()})
		return
	}

	u.Nickname = nick

	ur, err := h.uc.UpdateUser(&u)

	if err != nil {
		if err.Error() == "not found user" {
			c.AbortWithStatusJSON(http.StatusNotFound, utils.Error{Error: "Can't find user with id " + nick})
			return
		}
		if err.Error() == "conflict" {
			c.AbortWithStatusJSON(http.StatusConflict, utils.Error{Error: "Conflict status in " + nick})
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, ur)
}
