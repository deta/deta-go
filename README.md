# DETA SDK for Go

[![Go Doc](https://img.shields.io/badge/go-doc-blue)](https://godoc.org/github.com/deta/deta-go)

deta-go is the official Deta SDK for Go. 

## Installing

Use `go get` to retreive the SDK to add it to your `GOPATH` workspace, or project's Go module dependencies.

```
go get github.com/deta/deta-go
```

To update the SDK use `go get -u` to retrieve the latest version of the SDK.

```
go get -u github.com/deta/deta-go
```

If you are using Go modules, your `go get` will default to the latest tagged release version of the SDK. To get a specific release version of the SDK use `@<tag>` in your `go get` command.

```
go get github.com/deta/deta-go@v1.0.0
```

To get the latest SDK repository change use `@latest`.
```
go get github.com/aws/deta-go@latest
```

### SDK Packages

The SDK constitutes of two main components, the core package and service packages.

- `deta`: The core SDK package, provides shared functionalities to the service packages. All the `errors` are also exported from this package.

- `service`: The service packages, the services supported by the SDK.
	- `base`: Deta Base service package
	- `drive`: Deta Drive service package

### Configuring credentials

When using the SDK you will require you project key. The project key can be provided explicitly or is taken from the environement variable `DETA_PROJECT_KEY`.

#### Default

By default, the SDK looks for the environment variable `DETA_PROJECT_KEY` for the project key

```go
// Create a new Deta instance taking the project key from the environment by default
d, err := deta.New()
if err != nil {
	fmt.Fprintf(os.Stderr, "failed to create new deta instance: %v\n", err)
}
```

#### Provide the project key explicitly

You can use the `WithProjectKey` option when creating a `Deta` instance to provide the project key explicitly. 

```go
// Create a new Deta instance with explicit project key
d, err := deta.New(deta.WithProjectKey("project_key"))
if err != nil {
	fmt.Fprintf(os.Stderr, "failed to create new deta instance: %v\n", err)
}
```


## Examples

### Base

The following is a simple `Put` operation example.

```go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/deta/deta-go/deta"
	"github.com/deta/deta-go/service/base"
)

// User an example user struct
type User struct {
	// json struct tag 'key' used to denote the key
	Key      string   `json:"key"` 
	Username string   `json:"username"`
	Active   bool     `json:"active"`
	Age      int      `json:"age"`
	Likes    []string `json:"likes"`
	// json struct tag '__expires' for expiration timestamp
	// 'omitempty' to prevent default 0 value
	Expires  int64    `json:"__expires,omitempty"`
}

func main() {
	// Create a new Deta instance with a project key
	d, err := deta.New(deta.WithProjectKey("project_key"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create new Deta instance: %v\n", err)
		os.Exit(1)
	}

	// Create a new Base instance called "users", provide the previously created Deta instance 
	users, err := base.New(d, "users")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create new Base instance: %v\n", err)
		os.Exit(1)
	}

	// Put "jimmy" to the "users" Base
	key, err := users.Put(&User{
		Key: "jimmy_neutron", 
		Username: "jimmy",
		Active: true,
		Age: 20,
		Likes: []string{"science"},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to put item: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("successfully put item with key %s\n", key)

	// A map can also be used
	jimmy := map[string]interface{}{
		"key":      "jimmy_neutron",
		"username": "jimmy",
		"active":   true,
		"age":      20,
		"likes":    []string{"science"},
	}
	key, err = users.Put(jimmy)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to put item: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("successfully put item with key: %s\n", key)

	// Put with expiration timestamp
	eu := &User{
		Key: "tmp_user_key",
		Username: "test_user",
		Expires: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}
	key, err = users.Put(eu)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to put item: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("successfully put expiring item with key: %s\n", key)

	// Put map with expiration timestamp
	tmp := map[string]interface{}{
		"key": "tmp_user_key",
		"username": "test_user",
		// use `__expires` as the key for expiration timestamp
		"__expires": time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}
	key, err = users.Put(tmp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to put item: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("successfully put expiring item with key: %s\n", key)
}
```

More examples and complete documentation on https://docs.deta.sh/docs/base/sdk

### Drive

The following is a simple `Put` operation example.

```go
package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/deta/deta-go/deta"
	"github.com/deta/deta-go/service/drive"
)

func main() {
	// Create a new Deta instance with a project key
	d, err := deta.New(deta.WithProjectKey("project_key"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create new Deta instance:%v\n", \n)
		os.Exit(1)
	}

	// Create a new Drive instance called "drawings", provide the previously created Deta instance
	drawings, err := drive.New(d, "drawings")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create new Drive instance: %v\n", err)
		os.Exit(1)
	}

	// Open local file "art.svg"
	file, err := os.Open("./art.svg")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Put "art.svg" to "drawings"
	name, err := drawings.Put(&drive.PutInput{
		Name:        "art.svg",
		Body:        bufio.NewReader(file),
		ContentType: "image/svg+xml",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to put file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("successfully put file %s", name)
}
```

More examples and complete documentation on https://docs.deta.sh/docs/drive/sdk/

