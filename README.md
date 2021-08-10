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
go get github.com/deta/deta-go@v0.0.1
```

To get the latest SDK repository change use `@latest`.
```
go get github.com/aws/deta-go@latest
```

## Examples

### Base

```go

import (
	"fmt"

	"github.com/deta/deta-go/deta"
	"github.com/deta/deta-go/service/base"
)

type User struct {
	Key      string `json:"key"` // json struct tag key to denote the key
	Username string `json:"username"`
	Email    string `json:"email"`
}

func main() {
	d, err := deta.New(deta.WithProjectKey("project_key"))
	if err != nil {
		fmt.Println("failed to init new Deta instance:", err)
		return
	}

	db := base.New(d, "users")

	u := &User{
		Key:      "abasd",
		Username: "jimmy",
		Email:    "jimmy@deta.sh",
	}
	key, err := db.Put(u)
	if err != nil {
		fmt.Println("failed to put item:", err)
		return
	}
	fmt.Println("successfully put item with key", key)
}
```

More examples and complete documentation on https://docs.deta.sh/docs/base/sdk

### Drive

#### Put
```go
import (
	"bufio"
	"fmt"
	"os"

	"github.com/deta/deta-go/deta"
	"github.com/deta/deta-go/service/drive"
)

func main() {

	// initialize with project key
	// returns ErrBadProjectKey if project key is invalid
	d, err := deta.New(deta.WithProjectKey("project_key"))
	if err != nil {
		fmt.Println("failed to init new Deta instance:", err)
		return
	}

	// initialize with drive name
	// returns ErrBadDriveName if drive name is invalid
	drawings := drive.New(d, "drawings")

	// PUT
	// reading from a local file
	file, err := os.Open("./art.svg")
	defer file.Close()

	name, err := drawings.Put(&drive.PutInput{
		Name:        "art.svg",
		Body:        bufio.NewReader(file),
		ContentType: "image/svg+xml",
	})
	if err != nil {
		fmt.Println("Failed to put file:", err)
		return
	}
	fmt.Println("Successfully put file with name:", name)
}
```

#### Get
```go
import (
	"fmt"
	"io/ioutil"

	"github.com/deta/deta-go/deta"
	"github.com/deta/deta-go/service/drive"
)

func main() {

	// initialize with project key
	// returns ErrBadProjectKey if project key is invalid
	d, err := deta.New(deta.WithProjectKey("project_key"))
	if err != nil {
		fmt.Println("failed to init new Deta instance:", err)
		return
	}

	// initialize with drive name
	// returns ErrBadDriveName if drive name is invalid
	drawings := drive.New(d, "drawings")

	// GET
	name := "art.svg"
	f, err := drawings.Get(name)
	if err != nil {
		fmt.Println("Failed to get file with name:", name)
		return
	}
	defer f.Close()

	c, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("Failed read file content with err:", err)
		return
	}
	fmt.Println("file content:", string(c))
}
```

#### Delete
```go
import (
	"fmt"

	"github.com/deta/deta-go/deta"
	"github.com/deta/deta-go/service/drive"
)

func main() {

	// initialize with project key
	// returns ErrBadProjectKey if project key is invalid
	d, err := deta.New(deta.WithProjectKey("project_key"))
	if err != nil {
		fmt.Println("failed to init new Deta instance:", err)
		return
	}

	// initialize with drive name
	// returns ErrBadDriveName if drive name is invalid
	drawings := drive.New(d, "drawings")

	// DELETE
	name, err := drawings.Delete("art.svg")
	if err != nil {
		fmt.Println("Failed to delete file with name:", name)
		return
	}
	fmt.Println("Successfully deleted file with name:", name)
}
```

#### List
```go
import (
	"fmt"

	"github.com/deta/deta-go/deta"
	"github.com/deta/deta-go/service/drive"
)

func main() {

	// initialize with project key
	// returns ErrBadProjectKey if project key is invalid
	d, err := deta.New(deta.WithProjectKey("project_key"))
	if err != nil {
		fmt.Println("failed to init new Deta instance:", err)
		return
	}

	// initialize with drive name
	// returns ErrBadDriveName if drive name is invalid
	drawings := drive.New(d, "drawings")

	// LIST
	lr, err := drawings.List(1000, "", "")
	if err != nil {
		fmt.Println("Failed to list names from drive with err:", err)
	}
	fmt.Println("names:", lr.Names)
}
```
More examples and complete documentation on https://docs.deta.sh/docs/drive/sdk/
