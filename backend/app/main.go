package main

import (

	"codejam.io/config"
	"codejam.io/database"

	//"codejam.io/logging"
	"flag"

	"codejam.io/server"
	"github.com/gin-gonic/gin"
)

func main() {
	debugArg := flag.Bool("debug", false, "Enable debug mode.")
	flag.Parse()
	// TODO setup logger

	// Disables the debug logging... Comment this out to enable debug logging for GIN
	gin.SetMode(gin.ReleaseMode)

	// Un-comment this to view a logging output example
	//logger := logging.NewLogger(logging.Options{Name: "Main", Level: logging.DEBUG})
	//logger.Debug("Testing debug")
	//logger.Info("Testing info")
	//logger.Warn("Testing warn")
	//logger.Error("Testing error")
	//logger.Critical("Testing critical")

	config := new(config.Config)
	config.LoadFromFile("config.toml")

	database.Initialize(config.Database)
	database.RunMigrations()

	server := server.Server{
		Config: *config,
		Debug: *debugArg,
	}
	server.StartServer()
}
