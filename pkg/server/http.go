package server

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/influenzanet/counter-service/pkg/types"
)

type HTTPServer struct {
	port        int
	counters    map[string]types.CounterService
	platform    string
	extra       []string
	root        []string
	metaAuthKey string
}

func NewHTTPServer(config types.ServiceConfig, counters map[string]types.CounterService) *HTTPServer {
	port := config.Port
	if port == 0 {
		port = 5021
	}

	root := make([]string, 0, 1)
	extra := make([]string, 0)
	for _, counter := range counters {
		name := counter.Name()
		if counter.IsRoot() {
			root = append(root, name)
		} else {
			if counter.IsPublic() {
				extra = append(extra, name)
			}
		}
	}

	return &HTTPServer{
		port:        port,
		platform:    config.Platform,
		counters:    counters,
		root:        root,
		extra:       extra,
		metaAuthKey: config.MetaAuthKey,
	}
}

func (h *HTTPServer) Start() {

	r := gin.Default()

	r.Use(cors.Default())

	r.GET("/", func(c *gin.Context) {

		rootCounters := make(map[string]types.CounterData)

		for _, name := range h.root {
			registry := h.counters[name]
			data := registry.Data()
			rootCounters[name] = data
		}

		r := RootResponse{
			Extra:    h.extra,
			Platform: h.platform,
			Counters: rootCounters,
		}

		c.JSON(200, r)

	})

	if h.metaAuthKey != "" {

		authMiddleWare := gin.BasicAuth(gin.Accounts{
			"admin": h.metaAuthKey,
		})

		r.GET("/meta.json", authMiddleWare, func(c *gin.Context) {
			m := make([]types.CounterServiceDefinition, 0, len(h.counters))
			for _, counter := range h.counters {
				m = append(m, counter.Definition())
			}
			c.JSON(200, m)
		})
	}

	r.GET("/counter/:name", func(c *gin.Context) {

		name := c.Param("name")

		if name == "" || len(name) > 64 {
			c.AbortWithStatusJSON(404, map[string]string{"error": "bad counter name"})
			return
		}

		registry, ok := h.counters[name]
		if !ok {
			c.AbortWithStatusJSON(404, map[string]string{"error": "Unknown counter"})
			return
		}

		counters := registry.Data()
		c.JSON(200, counters)
	})

	r.Run(fmt.Sprintf(":%d", h.port)) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
