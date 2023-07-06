package server

import (
	"net/http"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/influenzanet/counter-service/pkg/types"
)

type HTTPServer struct {
	port int
	meta *types.Meta
	registries map[string]types.RegistryService
}

func NewHTTPServer(port int, registries map[string]types.RegistryService, meta *types.Meta) *HTTPServer {
	if(port == 0) {
		port = 5021
	}
	return &HTTPServer{port: port, registries: registries, meta:meta}
}

func (h *HTTPServer) Start() {
	
	r := gin.Default()

	r.Use(cors.Default())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/whoami", func(c *gin.Context) {
		c.JSON(http.StatusOK, h.meta)
	})

	r.GET("/study/:studyKey", func(c *gin.Context) {

		studyKey := c.Param("studyKey")

		if(studyKey == ""|| len(studyKey) > 64) {
			c.AbortWithStatusJSON(404, map[string]string{"error": "bad studyKey"})
			return
		}

		registry, ok := h.registries[studyKey]
		if(!ok) {
			c.AbortWithStatusJSON(404, map[string]string{"error": "Unknown studyKey"})
			return
		}

		counters := registry.Read()
		c.JSON(200, counters)
	})

	r.Run(fmt.Sprintf(":%d", h.port)) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
