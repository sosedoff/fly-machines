package testdata

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	machines "github.com/sosedoff/fly-machines"
)

func fixture(path string) string {
	data, err := os.ReadFile(filepath.Join("testdata", path+".json"))
	if err != nil {
		panic(err)
	}
	return string(data)
}

func Server(appName string) *httptest.Server {
	srv := gin.New()
	srv.Use(func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
	})

	requireMachine := func(c *gin.Context) {
		if c.Param("id") == "foo" {
			c.AbortWithStatusJSON(404, gin.H{"error": "machine does not exist"})
		}
	}

	api := srv.Group("/v1/apps/" + appName)
	{
		api.POST("/machines/:id/lease", requireMachine, func(c *gin.Context) {
			c.String(200, fixture("create_lease"))
		})

		api.GET("/machines", func(c *gin.Context) {
			c.String(200, fixture("list"))
		})

		api.GET("/machines/:id", requireMachine, func(c *gin.Context) {
			c.String(200, fixture("get"))
		})

		api.POST("/machines", func(c *gin.Context) {
			input := machines.CreateInput{}
			if err := c.BindJSON(&input); err != nil {
				panic(err)
			}
			if input.Config == nil {
				c.AbortWithStatusJSON(400, gin.H{"error": "no config provided"})
				return
			}

			switch input.Name {
			case "fatal": // simulate crash
				c.AbortWithStatusJSON(500, gin.H{"error": "something went wrong"})
				return
			case "timeout": // simulate increased latency
				time.Sleep(time.Second)
			}
			c.String(200, fixture("get"))
		})

		api.GET("/machines/:id/wait", requireMachine, func(c *gin.Context) {
			c.JSON(200, gin.H{"ok": true})
		})

		api.POST("/machines/:id/stop", requireMachine, func(c *gin.Context) {
			c.JSON(200, gin.H{"ok": true})
		})

		api.DELETE("/machines/:id", requireMachine, func(c *gin.Context) {
			c.JSON(200, gin.H{"ok": true})
		})
	}

	return httptest.NewServer(srv.Handler())
}
