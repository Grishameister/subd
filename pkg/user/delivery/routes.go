package delivery

import (
	"github.com/Grishameister/subd/internal/database"
	"github.com/Grishameister/subd/pkg/user/repository"
	"github.com/Grishameister/subd/pkg/user/usecase"
	"github.com/gin-gonic/gin"
)

func AddUserRoutes(r *gin.Engine, db database.IDbConn) {
	rep := repository.New(db)
	uc := usecase.New(rep)
	handler := New(uc)

	r.POST("/api/user/:nickname/create", handler.CreateUser)
	r.POST("/api/user/:nickname/profile", handler.UpdateUser)
	r.GET("/api/user/:nickname/profile", handler.GetUser)
}
