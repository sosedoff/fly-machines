package machines_test

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	machines "github.com/sosedoff/fly-machines"
	"github.com/sosedoff/fly-machines/testdata"
)

func TestConfiguration(t *testing.T) {
	client := machines.NewClient("app")
	require.Equal(t, "app", client.GetAppName())
	require.Equal(t, machines.PublicBaseURL, client.GetBaseURL())

	client.SetAppName("foo")
	client.SetBaseURL("http://hostname")
	require.Equal(t, "foo", client.GetAppName())
	require.Equal(t, "http://hostname", client.GetBaseURL())
}

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

	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()
	_, err = client.CreateContext(ctx, &machines.CreateInput{Name: "timeout", Config: &machines.Config{}})
	require.ErrorContains(t, err, "context deadline exceeded")
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

	err = client.WaitContext(context.Background(), &machines.WaitInput{ID: "1", State: machines.StateStopped})
	require.NoError(t, err)
}

func TestWaitHelpers(t *testing.T) {
	err := client.WaitStarted(context.Background(), &machines.Machine{ID: "1"})
	require.NoError(t, err)

	err = client.WaitStopped(context.Background(), &machines.Machine{ID: "1"})
	require.NoError(t, err)

	err = client.WaitDestroyed(context.Background(), &machines.Machine{ID: "1"})
	require.NoError(t, err)
}

func TestStopContext(t *testing.T) {
	err := client.StopContext(context.Background(), nil)
	require.Equal(t, machines.ErrInputRequired, err)

	err = client.StopContext(context.Background(), &machines.StopInput{})
	require.Equal(t, machines.ErrMachineIDRequired, err)

	err = client.StopContext(context.Background(), &machines.StopInput{ID: "foo"})
	require.Equal(t, "machine does not exist", err.Error())

	err = client.StopContext(context.Background(), &machines.StopInput{ID: "1"})
	require.NoError(t, err)
}

var (
	client *machines.Client
	server *httptest.Server
)

func init() {
	server = testdata.Server("app")
	client = testClient(server.URL)
}

func testClient(url string) *machines.Client {
	client := machines.NewClient("app")
	client.SetBaseURL(url)
	client.SetToken("api_token")
	return client
}
