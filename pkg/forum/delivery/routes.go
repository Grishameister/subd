package delivery

import (
	"github.com/Grishameister/subd/internal/database"
	"github.com/Grishameister/subd/pkg/forum/repository"
	"github.com/Grishameister/subd/pkg/forum/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AddForumRoutes(r *gin.Engine, db database.IDbConn) {
	rep := repository.New(db)
	uc := usecase.New(rep)
	handler := New(uc)

	r.POST("/api/forum/:slug", func(c *gin.Context) {
		if strings.HasPrefix(c.Request.RequestURI, "/api/forum/create") {
			handler.CreateForum(c)
			return
		}
		c.AbortWithStatus(http.StatusNotFound)
	})

	r.GET("/api/forum/:slug/details", handler.GetForum)
	r.GET("/api/forum/:slug/threads", handler.GetThreads)
	r.GET("/api/forum/:slug/users", handler.GetUsers)
}
