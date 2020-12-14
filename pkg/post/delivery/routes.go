package delivery

import (
	"github.com/Grishameister/subd/internal/database"
	frRepo "github.com/Grishameister/subd/pkg/forum/repository"
	"github.com/Grishameister/subd/pkg/post/repository"
	"github.com/Grishameister/subd/pkg/post/usecase"
	thrRepo "github.com/Grishameister/subd/pkg/thread/repository"
	urRepo "github.com/Grishameister/subd/pkg/user/repository"
	"github.com/gin-gonic/gin"
)

func AddPostsRoutes(r *gin.Engine, db database.IDbConn) {
	tr := thrRepo.New(db)
	fr := frRepo.New(db)
	ur := urRepo.New(db)
	rep := repository.New(db, tr, fr, ur)
	uc := usecase.New(rep)
	handler := New(uc)

	r.POST("/api/thread/:slug_or_id/create", handler.CreatePosts)

	r.GET("/api/thread/:slug_or_id/posts", handler.GetPosts)
	r.GET("/api/post/:id/details", handler.GetPost)
	r.POST("/api/post/:id/details", handler.UpdateMessage)
}
