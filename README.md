# fly-machines

Golang client for Fly.io [Machines API](https://fly.io/docs/machines/working-with-machines/)

_Status: Work in progress_

## Installation

```bash
go get github.com/sosedoff/fly-machines
```

## Usage Example

```golang
package main

import(
  machines "github.com/sosedoff/fly-machines"
)

func main() {
  // Initialize machines client for the app, and pull API url and token from env vars
  client := machines.NewClient("myapp")

  // Initialize the clients directly
  client := machines.NewClientWithToken("myapp", "api_token")

  // Extra configuration, if necessary
  client.SetAppName("myapp")
  client.SetToken("api_token")
  client.SetBaseURL(machines.PrivateBaseURL)

  // Methods
  client.List()
  client.Get()
  client.Create()
  client.Stop()
  client.Delete()
  client.Wait()

  // Waiting helpers
  client.WaitStarted()
  client.WaitStopped()
  client.WaitDestroyed()
}
```
