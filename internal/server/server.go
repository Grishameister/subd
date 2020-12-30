package server

import (
	"context"
	"fmt"
	"github.com/Grishameister/subd/configs/config"
	"github.com/Grishameister/subd/internal/database"
	forumDelivery "github.com/Grishameister/subd/pkg/forum/delivery"
	postsDelivery "github.com/Grishameister/subd/pkg/post/delivery"
	serviceDelivery "github.com/Grishameister/subd/pkg/service/delivery"
	threadDelivery "github.com/Grishameister/subd/pkg/thread/delivery"
	"github.com/Grishameister/subd/pkg/user/delivery"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Server struct {
	logFile *os.File
	server  *http.Server
}

func New(config *config.Config, db database.IDbConn) *Server {
	gin.SetMode(gin.ReleaseMode)
	logFile := setupGinLogger()

	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20

	delivery.AddUserRoutes(r, db)
	threadDelivery.AddThreadRoutes(r, db)
	forumDelivery.AddForumRoutes(r, db)
	postsDelivery.AddPostsRoutes(r, db)
	serviceDelivery.AddServiceRoutes(r, db)

	return &Server{
		logFile: logFile,
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%s", config.Web.Server.Address, config.Web.Server.Port),
			Handler: r,
		},
	}
}

func (s *Server) Run() {
	defer s.logFile.Close()

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	config.Lg("server", "Run").Info("Server listening on " + s.server.Addr)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	config.Lg("server", "Run").Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		config.Lg("server", "Run").Fatal("Server forced to shutdown:", err)
	}
}

func setupGinLogger() *os.File {
	switch strings.ToLower(config.Conf.Logger.GinLevel) {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	case "debug":
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	if !config.Conf.Logger.StdoutLog {
		file, err := os.OpenFile(config.Conf.Logger.GinFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			config.Lg("server", "setupGinLogger").Fatal("Failed to log to file, using default stderr")
			return nil
		}

		gin.DefaultWriter = io.MultiWriter(file)
		return file
	} else {
		return nil
	}
}
