package api

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	ctx context.Context
)

// Server 配置服务器
type Server struct {
	Addr   string
	Router *gin.Engine
}

// Run 运行配置服务器
func (s *Server) Run(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	log.Printf("Controller Server started. %s\n", s.Addr)

	srv := &http.Server{
		Addr:    s.Addr,
		Handler: s.Router,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	srvCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(srvCtx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Controller Server exiting.")
}

func (s *Server) InitRouters(f embed.FS) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	api := router.Group("/api")
	api.POST("/conversion-status", ConversionStatus)
	api.POST("/convert-file", ConvertFile)
	router.GET("/download-file/:flag", DownloadFile)
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"action": s.Addr,
		})
	})
	ht := template.Must(template.New("").ParseFS(f, "templates/*.tmpl"))

	router.SetHTMLTemplate(ht)
	s.Router = router
	return router
}

// NewServer 实例化配置服务器
func NewServer(ip, port string, _ctx context.Context) *Server {
	addr := fmt.Sprintf(
		"%s:%s",
		ip,
		port,
	)
	server := &Server{
		Addr: addr,
	}
	ctx = _ctx
	return server
}
