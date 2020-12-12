package delivery

import (
	"github.com/Grishameister/subd/internal/database"
	"github.com/Grishameister/subd/pkg/forum/repository"
	"github.com/Grishameister/subd/pkg/forum/usecase"
	"github.com/gin-gonic/gin"
)

func AddForumRoutes(r *gin.Engine, db database.IDbConn) {
	rep := repository.New(db)
	uc := usecase.New(rep)
	handler := New(uc)

	r.POST("/api/forum/create", handler.CreateForum)
	r.GET("/api/forum/:slug/details", handler.GetForum)
}
