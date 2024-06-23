package server

import (
	"embed"
	_ "embed"
	"github.com/gin-gonic/gin"
	"mime"
	"net/http"
	"path/filepath"
)

//go:embed static_files/*
var htmlFS embed.FS

func (server *Server) setupInternalStaticRoute() {
	server.Gin.NoRoute(func(ctx *gin.Context) {
		filename := ctx.Request.RequestURI
		if filename == "/" {
			filename = "/index.html"
		}

		content, err := htmlFS.ReadFile("static_files" + filename)
		if err != nil {
			logger.Error("htmlFS.ReadFile: %+v", err)
			ctx.Writer.WriteHeader(http.StatusNotFound)
			return
		}

		mimeType := mime.TypeByExtension(filepath.Ext(filename))
		if mimeType == "" {
			mimeType = "application/text"
		}

		ctx.Data(http.StatusOK, mimeType, content)
	})
}

func (server *Server) setupProxiedStaticRoute() {
	server.Gin.NoRoute(func(ctx *gin.Context) {
		filename := ctx.Request.RequestURI
		if filename == "/" {
			filename = "/index.html"
		}

		response, err := http.Get(server.Config.Server.StaticProxy + filename)
		if err != nil {
			ctx.Status(http.StatusServiceUnavailable)
			return
		}

		reader := response.Body
		contentLength := response.ContentLength
		contentType := response.Header.Get("Content-Type")

		ctx.DataFromReader(http.StatusOK, contentLength, contentType, reader, nil)
	})
}

// SetupStaticRoutes assumes that any request that doesn't match a specific router, should be treated
// as an attempt to load a static file.
func (server *Server) SetupStaticRoutes() {
	if server.Config.Server.StaticProxy != "" {
		server.setupProxiedStaticRoute()
	} else {
		server.setupInternalStaticRoute()
	}
}
