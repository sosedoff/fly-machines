package machines_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	machines "github.com/sosedoff/fly-machines"
	"github.com/stretchr/testify/require"
)

func TestListContext(t *testing.T) {
	list, err := client.ListContext(context.Background(), nil)
	require.NoError(t, err)
	require.Equal(t, 1, len(list))
	require.Equal(t, "4d89040f431938", list[0].ID)
	require.Equal(t, "winter-cloud-3782", list[0].Name)
}

func TestGetContext(t *testing.T) {
	_, err := client.GetContext(context.Background(), nil)
	require.Equal(t, err, machines.ErrInputRequired)

	_, err = client.GetContext(context.Background(), &machines.GetInput{})
	require.Equal(t, err, machines.ErrMachineIDRequired)

	_, err = client.GetContext(context.Background(), &machines.GetInput{ID: "foo"})
	require.Equal(t, "machine does not exist", err.Error())

	machine, err := client.GetContext(context.Background(), &machines.GetInput{ID: "1"})
	require.Nil(t, err)
	require.Equal(t, "4d89040f431938", machine.ID)
	require.Equal(t, "winter-cloud-3782", machine.Name)
}

func TestCreateContext(t *testing.T) {
	_, err := client.CreateContext(context.Background(), nil)
	require.Equal(t, machines.ErrInputRequired, err)

	_, err = client.CreateContext(context.Background(), &machines.CreateInput{})
	require.Equal(t, "no config provided", err.Error())

	machine, err := client.CreateContext(context.Background(), &machines.CreateInput{Config: &machines.Config{}})
	require.NoError(t, err)
	require.Equal(t, "4d89040f431938", machine.ID)
	require.Equal(t, "winter-cloud-3782", machine.Name)

	_, err = client.CreateContext(context.Background(), &machines.CreateInput{Name: "fatal", Config: &machines.Config{}})
	require.Equal(t, "something went wrong", err.Error())
}

func TestLeaseContext(t *testing.T) {
	_, err := client.LeaseContext(context.Background(), nil)
	require.Equal(t, err, machines.ErrInputRequired)

	_, err = client.LeaseContext(context.Background(), &machines.LeaseInput{})
	require.Equal(t, err, machines.ErrMachineIDRequired)

	_, err = client.LeaseContext(context.Background(), &machines.LeaseInput{ID: "foo"})
	require.Equal(t, "machine does not exist", err.Error())

	lease, err := client.LeaseContext(context.Background(), &machines.LeaseInput{ID: "1"})
	require.NoError(t, err)
	require.Equal(t, &machines.Lease{
		Nonce:     "1234",
		ExpiresAt: 1679456889,
		Owner:     "owner@corp.com",
	}, lease)
}

func TestDeleteContext(t *testing.T) {
	err := client.DeleteContext(context.Background(), nil)
	require.Equal(t, machines.ErrInputRequired, err)

	err = client.DeleteContext(context.Background(), &machines.DeleteInput{})
	require.Equal(t, machines.ErrMachineIDRequired, err)

	err = client.DeleteContext(context.Background(), &machines.DeleteInput{ID: "foo"})
	require.Equal(t, "machine does not exist", err.Error())

	err = client.DeleteContext(context.Background(), &machines.DeleteInput{ID: "1"})
	require.NoError(t, err)
}

func TestWaitContext(t *testing.T) {
	err := client.WaitContext(context.Background(), nil)
	require.Equal(t, machines.ErrInputRequired, err)

	err = client.WaitContext(context.Background(), &machines.WaitInput{})
	require.Equal(t, machines.ErrMachineIDRequired, err)

	err = client.WaitContext(context.Background(), &machines.WaitInput{ID: "foo"})
	require.Equal(t, machines.ErrInvalidWaitState, err)

	err = client.WaitContext(context.Background(), &machines.WaitInput{ID: "foo", State: machines.StateStarted})
	require.Equal(t, "machine does not exist", err.Error())

	// err = client.WaitContext(context.Background(), &machines.WaitInput{ID: "1"})
	// require.NoError(t, err)
}

var (
	client *machines.Client
	server *http.Server
)

func init() {
	client = testClient()
	server = testServer()

	go func() {
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()
}

func testClient() *machines.Client {
	client := machines.NewClient("app")
	client.SetBaseURL("http://localhost:30555")
	client.SetToken("api_token")
	return client
}

func fixture(path string) string {
	data, err := os.ReadFile(filepath.Join("testdata", path+".json"))
	if err != nil {
		panic(err)
	}
	return string(data)
}

func testServer() *http.Server {
	srv := gin.New()
	srv.Use(func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
	})

	requireMachine := func(c *gin.Context) {
		if c.Param("id") == "foo" {
			c.AbortWithStatusJSON(404, gin.H{"error": "machine does not exist"})
		}
	}

	api := srv.Group("/v1/apps/" + client.GetAppName())
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

			if input.Name == "fatal" {
				c.AbortWithStatusJSON(500, gin.H{"error": "something went wrong"})
				return
			}

			c.String(200, fixture("get"))
		})

		api.GET("/machines/:id/wait", requireMachine, func(c *gin.Context) {
			c.JSON(200, gin.H{"ok": true})
		})

		api.DELETE("/machines/:id", requireMachine, func(c *gin.Context) {
			c.JSON(200, gin.H{"ok": true})
		})
	}

	return &http.Server{
		Addr:    ":30555",
		Handler: srv.Handler(),
	}
}
