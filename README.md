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
package main

import (
	"fmt"
	"github.com/deta/deta-go"
)

type User struct {
	Key string `json:"key"` // json struct tag key to denote the key
	Username string `json:"username"`
	Email string `json:"email"`
}

func main(){
	d, err := deta.New("project_key")
	if err != nil{
		fmt.Println("failed to init a new Deta instance:", err)
		return
	}

	db, err := d.NewBase("base_name")
	if err != nil{
		fmt.Println("failed to init a new Base instance:", err)
		return
	}

	u := &User{
		Key: "abasd",
		Username: "jimmy",
		Email: "jimmy@deta.sh",
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
```go
package main

import (
	"bufio"
	"fmt"
	"github.com/deta/deta-go"
	"io/ioutil"
	"os"
)

func main() {
	// initialize with project key
	// returns ErrBadProjectKey if project key is invalid
	d, err := deta.New("project_key")
	if err != nil {
		fmt.Println("failed to init new Deta instance:", err)
		return
	}

	// initialize with drive name
	// returns ErrBadDriveName if drive name is invalid
	drive, err := d.NewDrive("drive_name")
	if err != nil {
		fmt.Println("failed to init new Drive instance:", err)
		return
	}

	// PUT
	// reading from a local file
	file, err := os.Open("./art.svg")
	defer file.Close()

	name, err := drive.Put(&deta.PutInput{
		Name:        "art.svg",
		Body:        bufio.NewReader(file),
		ContentType: "image/svg+xml",
	})
	if err != nil {
		fmt.Println("Failed to put file:", err)
		return
	}
	fmt.Println("Successfully put file with name:", name)

	// GET
	name = "art.svg"
	f, err := drive.Get(name)
	if err != nil {
		fmt.Println("Failed to get file with name:", name, err)
		return
	}
	defer f.Close()

	c, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("Failed read file content with err:", err)
		return
	}
	fmt.Println("file content:", string(c))

	// DELETE
	name, err = drive.Delete("hello.txt")
	if err != nil {
		fmt.Println("Failed to delete file with name:", name)
		return
	}
	fmt.Println("Successfully deleted file with name:", name)

	// LIST
	lr, err := drive.List(1000, "", "")
	if err != nil {
		fmt.Println("Failed to list names from drive with err:", err)
	}
	fmt.Println("names:", lr.Names)
}
```
More examples and complete documentation on https://docs.deta.sh/docs/drive/sdk/
