package main

import (
	"github.com/Grishameister/subd/configs/config"
	"github.com/Grishameister/subd/internal/database"
	"github.com/Grishameister/subd/internal/server"
)

func main() {
	config.Conf = config.NewConfig()

	logger := config.Logger{}
	logger.Init()
	defer logger.Cleanup()



	dbConn := database.NewDB(&config.Conf.Db)
	if err := dbConn.Open(); err != nil {
		config.Lg("main", "main").Fatal("Connection refused")
		return
	}
	defer dbConn.Close()
	config.Lg("main", "main").Info("Connected to DB")

	srv := server.New(config.Conf, dbConn)

	srv.Run()

	config.Lg("main", "main").Info("Server stopped")
}