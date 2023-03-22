package machines

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func (c *Client) GetAppName() string {
	return c.appName
}

func (c *Client) SetToken(token string) {
	c.apiToken = token
}

func (c *Client) SetBaseURL(url string) {
	c.baseURL = url
}

func (c *Client) GetBaseURL() string {
	return c.baseURL
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
		input = &ListInput{} //nolint:ineffassign
	}

	var machines []Machine
	err = c.execute(req, &machines)
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

	var machine Machine
	err = c.execute(req, &machine)
	return &machine, err
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

	var machine Machine
	err = c.execute(req, &machine)
	return &machine, err
}

func (c *Client) Stop(input *StopInput) error {
	return c.StopContext(context.Background(), input)
}

func (c *Client) StopContext(ctx context.Context, input *StopInput) error {
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

	return c.execute(req, nil)
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

	return c.execute(req, nil)
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
	if input.Timeout > 0 {
		q.Add("timeout", input.Timeout.String())
	}

	req.URL.RawQuery = q.Encode()

	return c.execute(req, nil)
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
		ID:    machine.ID,
		State: StateDestroyed,
	})
}

func (c *Client) Lease(input *LeaseInput) (*Lease, error) {
	return c.LeaseContext(context.Background(), input)
}

func (c *Client) LeaseContext(ctx context.Context, input *LeaseInput) (*Lease, error) {
	if input == nil {
		return nil, ErrInputRequired
	}
	if err := input.Validate(); err != nil {
		return nil, err
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/machines/"+input.ID+"/lease", input)
	if err != nil {
		return nil, err
	}

	var lease Lease
	err = c.execute(req, &lease)
	return &lease, err
}

func (c *Client) ReleaseLeaseContext(ctx context.Context, input *LeaseInput) error {
	if input == nil {
		return ErrInputRequired
	}
	if err := input.Validate(); err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodDelete, "/machines/"+input.ID+"/lease", nil)
	if err != nil {
		return err
	}
	req.Header.Add("fly-machine-lease-nonce", input.Nonce)

	return c.execute(req, nil)
}

func (c *Client) urlForPath(path string) string {
	return fmt.Sprintf("%s/v1/apps/%s%s", c.baseURL, c.appName, path)
}

func (c *Client) newRequest(ctx context.Context, method string, path string, body any) (*http.Request, error) {
	if c.appName == "" {
		return nil, ErrAppNameRequired
	}
	if c.apiToken == "" {
		return nil, ErrAuthRequired
	}

	bodyReader, err := c.jsonBody(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, c.urlForPath(path), bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", ClientVersion())
	req.Header.Add("Authorization", "Bearer "+c.apiToken)
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

func (c *Client) execute(req *http.Request, out any) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		if out == nil {
			_, err := io.Copy(io.Discard, resp.Body)
			return err
		}
		return json.NewDecoder(resp.Body).Decode(out)
	}

	return apiErrorFromResponse(resp)
}

func (c *Client) jsonBody(body any) (io.Reader, error) {
	var bodyReader io.Reader

	if body != nil {
		buf := bytes.NewBuffer(nil)
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
		bodyReader = buf
	}

	return bodyReader, nil
}
