package delivery

import (
	"github.com/Grishameister/subd/internal/database"
	"github.com/Grishameister/subd/pkg/service/repository"
	"github.com/Grishameister/subd/pkg/service/usecase"
	"github.com/gin-gonic/gin"
)

func AddServiceRoutes(r *gin.Engine, db database.IDbConn) {
	repo := repository.New(db)
	uc := usecase.New(repo)
	handler := New(uc)

	r.GET("/api/service/status", handler.Status)
	r.POST("/api/service/clear", handler.Clear)

}
