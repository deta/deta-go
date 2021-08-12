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
#### Put
```go

import (
	"fmt"

	"github.com/deta/deta-go/deta"
	"github.com/deta/deta-go/service/base"
)

type User struct {
	Key      string   `json:"key"` // json struct tag 'key' used to denote the key
	Username string   `json:"username"`
	Active   bool     `json:"active"`
	Age      int      `json:"age"`
	Likes    []string `json:"likes"`
}

func main() {
	d, err := deta.New(deta.WithProjectKey("project_key"))
	if err != nil {
		fmt.Println("failed to init new Deta instance:", err)
		return
	}

	db, err := base.New(d, "users")
	if err != nil {
		fmt.Println("failed to init new Base instance:", err)
		return
	}

	u := &User{
		Key:      "kasdlj1",
		Username: "jimmy",
		Active:   true,
		Age:      20,
		Likes:    []string{"ramen"},
	}
	key, err := db.Put(u)
	if err != nil {
		fmt.Println("failed to put item:", err)
		return
	}
	fmt.Println("successfully put item with key", key)

	// can also use a map
	um := map[string]interface{}{
		"key":      "kasdlj1",
		"username": "jimmy",
		"active":   true,
		"age":      20,
		"likes":    []string{"ramen"},
	}

	key, err = db.Put(um)
	if err != nil {
		fmt.Println("Failed to put item:", err)
		return
	}
	fmt.Println("Successfully put item with key:", key)
}
```

#### Get
```go
import (
	"fmt"

	"github.com/deta/deta-go/deta"
	"github.com/deta/deta-go/service/base"
)

type User struct {
	Key      string   `json:"key"` // json struct tag 'key' used to denote the key
	Username string   `json:"username"`
	Active   bool     `json:"active"`
	Age      int      `json:"age"`
	Likes    []string `json:"likes"`
}

func main() {
	d, err := deta.New(deta.WithProjectKey("project_key"))
	if err != nil {
		fmt.Println("failed to init new Deta instance:", err)
		return
	}

	db, err := base.New(d, "users")
	if err != nil {
		fmt.Println("failed to init new Base instance:", err)
	}

	// a variable to store the result
	var u User

	// get item
	// returns ErrNotFound if no item was found
	err = db.Get("kasdlj1", &u)
	if err != nil {
		fmt.Println("Failed to get item:", err)
	}
}
```
#### Insert
```go

import (
	"fmt"

	"github.com/deta/deta-go/deta"
	"github.com/deta/deta-go/service/base"
)

type User struct {
	Key      string   `json:"key"` // json struct tag 'key' used to denote the key
	Username string   `json:"username"`
	Active   bool     `json:"active"`
	Age      int      `json:"age"`
	Likes    []string `json:"likes"`
}

func main() {
	d, err := deta.New(deta.WithProjectKey("project_key"))
	if err != nil {
		fmt.Println("failed to init new Deta instance:", err)
		return
	}

	db, err := base.New(d, "users")
	if err != nil {
		fmt.Println("failed to init new Base instance:", err)
	}

	u := &User{
		Key:      "kasdlj1",
		Username: "jimmy",
		Active:   true,
		Age:      20,
		Likes:    []string{"ramen"},
	}

	// insert item in the database
	key, err := db.Insert(u)
	if err != nil {
		fmt.Println("Failed to insert item:", err)
		return
	}
	fmt.Println("Successfully inserted item with key:", key)
}

```
#### Delete
```go
import (
	"fmt"

	"github.com/deta/deta-go/deta"
	"github.com/deta/deta-go/service/base"
)

func main() {
	d, err := deta.New(deta.WithProjectKey("project_key"))
	if err != nil {
		fmt.Println("failed to init new Deta instance:", err)
		return
	}

	db, err := base.New(d, "users")
	if err != nil {
		fmt.Println("failed to init new Base instance:", err)
	}

	// delete item
	// returns a nil error if item was not found
	err = db.Delete("dakjkfa")
	if err != nil {
		fmt.Println("Failed to delete item:", err)
	}
}

```
#### Put Many
```go
import (
    "fmt"

    "github.com/deta/deta-go/deta"
    "github.com/deta/deta-go/service/base"
)

type User struct {
    Key      string   `json:"key"` // json struct tag 'key' used to denote the key
    Username string   `json:"username"`
    Active   bool     `json:"active"`
    Age      int      `json:"age"`
    Likes    []string `json:"likes"`
}

func main() {
    d, err := deta.New(deta.WithProjectKey("project_key"))
    if err != nil {
        fmt.Println("failed to init new Deta instance:", err)
        return
    }

    db, err := base.New(d, "users")
    if err != nil {
        fmt.Println("failed to init new Base instance:", err)
    }

    // users
    u1 := &User{
        Key:      "kasdlj1",
        Username: "jimmy",
        Active:   true,
        Age:      20,
        Likes:    []string{"ramen"},
    }
    u2 := &User{
        Key:      "askdjf",
        Username: "joel",
        Active:   true,
        Age:      23,
        Likes:    []string{"coffee"},
    }
    users := []*User{u1, u2}

    // put items in the database
    keys, err := db.PutMany(users)
    if err != nil {
        fmt.Println("Failed to put items:", err)
        return
    }
    fmt.Println("Successfully put item with keys:", keys)
}
```
#### Update
```go

import (
	"fmt"

	"github.com/deta/deta-go/deta"
	"github.com/deta/deta-go/service/base"
)

type User struct {
	Key      string   `json:"key"` // json struct tag 'key' used to denote the key
	Username string   `json:"username"`
	Active   bool     `json:"active"`
	Age      int      `json:"age"`
	Likes    []string `json:"likes"`
}

func main() {
	d, err := deta.New(deta.WithProjectKey("project_key"))
	if err != nil {
		fmt.Println("failed to init new Deta instance:", err)
		return
	}

	db, err := base.New(d, "users")
	if err != nil {
		fmt.Println("failed to init new Base instance:", err)
	}

	// define the updates
	updates := base.Updates{
		"age": 33, // set profile.age to 33
	}
	// update
	err = db.Update("kasdlj1", updates)
	if err != nil {
		fmt.Println("failed to update")
		return
	}
}
```
#### Fetch 
```go
import (
	"fmt"

	"github.com/deta/deta-go/deta"
	"github.com/deta/deta-go/service/base"
)

type User struct {
	Key      string   `json:"key"` // json struct tag 'key' used to denote the key
	Username string   `json:"username"`
	Active   bool     `json:"active"`
	Age      int      `json:"age"`
	Likes    []string `json:"likes"`
}

func main() {
	d, err := deta.New(deta.WithProjectKey("project_key"))
	if err != nil {
		fmt.Println("failed to init new Deta instance:", err)
		return
	}

	db, err := base.New(d, "users")
	if err != nil {
		fmt.Println("failed to init new Base instance:", err)
	}

	// query to get users with age less than 30
	query := base.Query{
		{"age?lt": 30},
	}

	// variabe to store the results
	var results []*User

	// fetch items
	_, err = db.Fetch(&base.FetchInput{
		Q:    query,
		Dest: &results,
	})
	if err != nil {
		fmt.Println("failed to fetch items:", err)
	}
}
```
#### Fetch Paginated
```go
import (
	"fmt"

	"github.com/deta/deta-go/deta"
	"github.com/deta/deta-go/service/base"
)

type User struct {
	Key      string   `json:"key"` // json struct tag 'key' used to denote the key
	Username string   `json:"username"`
	Active   bool     `json:"active"`
	Age      int      `json:"age"`
	Likes    []string `json:"likes"`
}

func main() {
	d, err := deta.New(deta.WithProjectKey("project_key"))
	if err != nil {
		fmt.Println("failed to init new Deta instance:", err)
		return
	}

	db, err := base.New(d, "users")
	if err != nil {
		fmt.Println("failed to init new Base instance:", err)
	}

	// query to get users with age less than 30
	query := base.Query{
		{"age?lt": 30},
	}

	// variabe to store the results
	var results []*User

	// variable to store the page
	var page []*User

	// fetch input
	i := &base.FetchInput{
		Q:     query,
		Dest:  &page,
		Limit: 1, // limit provided so each page will only have one item
	}

	// fetch items
	lastKey, err := db.Fetch(i)
	if err != nil {
		fmt.Println("failed to fetch items:", err)
		return
	}

	// append page items to results
	results = append(results, page...)

	// get all pages
	for lastKey != "" {
		// provide the last key in the fetch input
		i.LastKey = lastKey

		// fetch
		lastKey, err = db.Fetch(i)
		if err != nil {
			fmt.Println("failed to fetch items:", err)
			return
		}

		// append page items to all results
		results = append(results, page...)
	}
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
	drawings, err := drive.New(d, "drawings")
	if err != nil {
		fmt.Println("failed to init new Drive instance:", err)
		return 
	}
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
	drawings, err := drive.New(d, "drawings")
	if err != nil {
		fmt.Println("failed to init new Drive instance:", err)
		return 
	}

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
	drawings, err := drive.New(d, "drawings")
	if err != nil {
		fmt.Println("failed to init new Drive instance:", err)
		return 
	}

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
	drawings, err := drive.New(d, "drawings")
	if err != nil {
		fmt.Println("failed to init new Drive instance:", err)
		return 
	}

	// LIST
	lr, err := drawings.List(1000, "", "")
	if err != nil {
		fmt.Println("Failed to list names from drive with err:", err)
	}
	fmt.Println("names:", lr.Names)
}
```
More examples and complete documentation on https://docs.deta.sh/docs/drive/sdk/
