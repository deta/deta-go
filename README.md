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

## Example

```go
package main

import (
	"deta"
	"fmt"
)

type User struct {
	Key string `json:"key"`
	Username string `json:"username"`
	Email string `json:"email"`
}

func main(){
	d, err := deta.New("project_key")
	if err != nil{
		fmt.Println("failed to init a new Deta instance:", err)
		return
	}

	db, err := deta.NewBase("base_name")
	if err != nil{
		fmt.Println("failed to init a new Base instance:", err)
		return
	}

	u := &User{
		Key: "abasd",
		Username: "jimmy",
		Email: "jimmy@deta.sh"
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