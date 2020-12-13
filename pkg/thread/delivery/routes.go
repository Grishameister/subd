package delivery

import (
	"github.com/Grishameister/subd/internal/database"
	"github.com/Grishameister/subd/pkg/thread/repository"
	"github.com/Grishameister/subd/pkg/thread/usecase"
	"github.com/gin-gonic/gin"
)

func AddThreadRoutes(r *gin.Engine, db database.IDbConn) {
	rep := repository.New(db)
	uc := usecase.New(rep)
	handler := New(uc)

	r.POST("/api/forum/:slug/create", handler.CreateThread)
	r.GET("/api/thread/:slug_or_id/details", handler.GetThread)
}
