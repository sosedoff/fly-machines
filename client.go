package machines

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	PublicBaseURL  = "https://api.machines.dev"
	PrivateBaseURL = "http://_api.internal:4280"
	DefaultBaseURL = PublicBaseURL
)

type Client struct {
	client   *http.Client
	baseURL  string
	apiToken string
	appName  string
}

func NewClient(appName string) *Client {
	return &Client{
		appName:  appName,
		client:   http.DefaultClient,
		baseURL:  envVarWithDefault("FLY_API_HOSTNAME", DefaultBaseURL),
		apiToken: envVarWithDefault("FLY_API_TOKEN", ""),
	}
}

func NewClientWithToken(appName, token string) *Client {
	client := NewClient(appName)
	client.SetToken(token)

	return client
}

func (c *Client) SetAppName(name string) {
	c.appName = name
}

func (c *Client) SetToken(token string) {
	c.apiToken = token
}

func (c *Client) SetBaseURL(url string) {
	c.baseURL = url
}

func (c *Client) List(input *ListInput) ([]Machine, error) {
	return c.ListContext(context.Background(), input)
}

func (c *Client) ListContext(ctx context.Context, input *ListInput) ([]Machine, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/machines", nil)
	if err != nil {
		return nil, err
	}
	if input == nil {
		input = &ListInput{}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, handleError(resp)
	}

	machines := []Machine{}
	err = json.NewDecoder(resp.Body).Decode(&machines)

	return machines, err
}

func (c *Client) Create(input *CreateInput) (*Machine, error) {
	return c.CreateContext(context.Background(), input)
}

func (c *Client) CreateContext(ctx context.Context, input *CreateInput) (*Machine, error) {
	if input == nil {
		return nil, ErrInputRequired
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/machines", input)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, handleError(resp)
	}

	machine := &Machine{}
	if err := json.NewDecoder(resp.Body).Decode(machine); err != nil {
		return nil, err
	}
	return machine, nil
}

func (c *Client) Get(input *GetInput) (*Machine, error) {
	return c.GetContext(context.Background(), input)
}

func (c *Client) GetContext(ctx context.Context, input *GetInput) (*Machine, error) {
	if input == nil {
		return nil, ErrInputRequired
	}
	if input.ID == "" {
		return nil, ErrMachineIDRequired
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/machines/"+input.ID, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	machine := &Machine{}
	err = json.NewDecoder(resp.Body).Decode(machine)

	return machine, err
}

func (c *Client) Stop(ctx context.Context, input *StopInput) error {
	if input == nil {
		return ErrInputRequired
	}
	if err := input.Validate(); err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/machines/"+input.ID+"/stop", nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return handleError(resp)
	}

	_, err = io.Copy(io.Discard, resp.Body)
	return err
}

func (c *Client) Delete(input *DeleteInput) error {
	return c.DeleteContext(context.Background(), input)
}

func (c *Client) DeleteContext(ctx context.Context, input *DeleteInput) error {
	if input == nil {
		return ErrInputRequired
	}
	if err := input.Validate(); err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodDelete, "/machines/"+input.ID, nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, resp.Body)

	if resp.StatusCode != 200 {
		return fmt.Errorf("received error: %v", resp.StatusCode)
	}

	return nil
}

func (c *Client) Wait(input *WaitInput) error {
	return c.WaitContext(context.Background(), input)
}

func (c *Client) WaitContext(ctx context.Context, input *WaitInput) error {
	if input == nil {
		return ErrInputRequired
	}
	if err := input.Validate(); err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/machines/"+input.ID+"/wait", nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("state", string(input.State))
	if input.InstanceID != "" {
		q.Add("instance_id", input.InstanceID)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, resp.Body)

	if resp.StatusCode != 200 {
		return fmt.Errorf("received error: %v", resp.StatusCode)
	}

	return nil
}

func (c *Client) WaitStarted(ctx context.Context, machine *Machine) error {
	return c.WaitContext(ctx, &WaitInput{ID: machine.ID, State: StateStarted})
}

func (c *Client) WaitStopped(ctx context.Context, machine *Machine) error {
	return c.WaitContext(ctx, &WaitInput{
		ID:         machine.ID,
		InstanceID: machine.InstanceID,
		State:      StateStopped,
	})
}

func (c *Client) WaitDestroyed(ctx context.Context, machine *Machine) error {
	return c.WaitContext(ctx, &WaitInput{
		ID: machine.ID,
		//InstanceID: machine.InstanceID,
		State: StateDestroyed,
	})
}

func (c *Client) urlForPath(path string) string {
	return fmt.Sprintf("%s/v1/apps/%s%s", c.baseURL, c.appName, path)
}

func (c *Client) newRequest(ctx context.Context, method string, path string, body any) (*http.Request, error) {
	if c.appName == "" {
		return nil, errors.New("app name must be set")
	}
	if c.apiToken == "" {
		return nil, errors.New("api token must be set")
	}

	var bodyReader io.Reader
	if body != nil {
		buf := bytes.NewBuffer(nil)
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
		bodyReader = buf
	}

	req, err := http.NewRequestWithContext(ctx, method, c.urlForPath(path), bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.apiToken)
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

func handleError(resp *http.Response) error {
	err := APIError{}

	if err := json.NewDecoder(resp.Body).Decode(&err); err != nil {
		return err
	}
	err.StatusCode = resp.StatusCode

	return err
}
