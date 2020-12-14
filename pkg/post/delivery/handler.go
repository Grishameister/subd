package delivery

import (
	"encoding/json"
	"github.com/Grishameister/subd/pkg/domain"
	"github.com/Grishameister/subd/pkg/post"
	"github.com/Grishameister/subd/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	uc post.IUsecase
}

func New(uc post.IUsecase) *Handler {
	return &Handler{
		uc: uc,
	}
}

func (h *Handler) CreatePosts(c *gin.Context) {
	slugOrId := c.Param("slug_or_id")
	var posts []*domain.Post

	decoder := json.NewDecoder(c.Request.Body)

	if err := decoder.Decode(&posts); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.Error{Error: err.Error()})
		return
	}

	postsr, err := h.uc.CreatePosts(slugOrId, posts)

	if err != nil {
		if err.Error() == "post not found" {
			c.AbortWithStatusJSON(http.StatusConflict, utils.Error{Error: err.Error()})
			return
		}
		if err.Error() == "thread not found" {
			c.AbortWithStatusJSON(http.StatusNotFound, utils.Error{Error: err.Error()})
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusCreated, postsr)
}

func (h *Handler) GetPosts(c *gin.Context) {
	slugOrId := c.Param("slug_or_id")

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

	sort := c.Query("sort")

	posts, err := h.uc.GetPosts(slugOrId, limit, since, sort, order)

	if err != nil {
		if err.Error() == "thread not found" {
			c.AbortWithStatusJSON(http.StatusNotFound, utils.Error{Error: err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *Handler) GetPost(c *gin.Context) {
	id := c.Param("id")

	related := c.Query("related")

	p, err := h.uc.GetPost(id, related)

	if err != nil {
		if err.Error() == "post not found" {
			c.AbortWithStatusJSON(http.StatusNotFound, utils.Error{Error: err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *Handler) UpdateMessage(c *gin.Context) {
	id := c.Param("id")

	m := domain.PostUpdate{}

	if err := c.BindJSON(&m); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.Error{Error: err.Error()})
		return
	}

	p, err := h.uc.UpdatePost(id, m.Message)

	if err != nil {
		if err.Error() == "post not found" {
			c.AbortWithStatusJSON(http.StatusNotFound, utils.Error{Error: err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}
