package delivery

import (
	"github.com/Grishameister/subd/pkg/domain"
	"github.com/Grishameister/subd/pkg/forum"
	"github.com/Grishameister/subd/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	uc forum.IUseCase
}

func New(uc forum.IUseCase) *Handler {
	return &Handler{
		uc: uc,
	}
}

func (h *Handler) CreateForum(c *gin.Context) {
	f := domain.Forum{}

	if err := c.BindJSON(&f); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	fr, err := h.uc.CreateForum(&f)

	if err != nil {
		if err.Error() == "slug exists" {
			c.AbortWithStatusJSON(http.StatusConflict, fr)
			return
		}
		if err.Error() == "user not found" {
			c.AbortWithStatusJSON(http.StatusNotFound, utils.Error{Error: "user not found is " + f.Owner})
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusCreated, f)
}

func (h *Handler) GetForum(c *gin.Context) {
	slug := c.Param("slug")

	f, err := h.uc.GetForum(slug)

	if err != nil {
		if err.Error() == "forum not found" {
			c.AbortWithStatusJSON(http.StatusNotFound, utils.Error{Error: "forum not found is " + slug})
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, f)
}

func (h *Handler) GetThreads(c *gin.Context) {
	slug := c.Param("slug")

	limit := c.Query("limit")

	if limit == "" {
		limit = "100"
	}

	since := c.Query("since")

	order := c.Query("desc")

	if order == "false" || order == "" {
		order = "asc"
	} else {
		order = "desc"
	}

	threads, err := h.uc.GetThreads(slug, limit, since, order)

	if err != nil {
		if err.Error() == "forum not found" {
			c.AbortWithStatusJSON(http.StatusNotFound, utils.Error{Error: err.Error()})
			return
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
	c.JSON(http.StatusOK, threads)
}

func (h *Handler) GetUsers(c *gin.Context) {
	slug := c.Param("slug")

	limit := c.Query("limit")

	if limit == "" {
		limit = "100"
	}

	since := c.Query("since")

	order := c.Query("desc")

	if order == "false" || order == "" {
		order = "asc"
	} else {
		order = "desc"
	}

	users, err := h.uc.GetUsers(slug, limit, since, order)

	if err != nil {
		if err.Error() == "forum not found" {
			c.AbortWithStatusJSON(http.StatusNotFound, utils.Error{Error: err.Error()})
			return
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
	c.JSON(http.StatusOK, users)
}