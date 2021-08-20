/*
Package base is the Deta Base service package.

The following is a simple Put operation example.

	import (
		"fmt"
		"os"

		"github.com/deta/deta-go/deta"
		"github.com/deta/deta-go/service/base"
	)

	// User an example user struct
	type User struct {
		Key      string   `json:"key"` // json struct tag 'key' used to denote the key
		Username string   `json:"username"`
		Active   bool     `json:"active"`
		Age      int      `json:"age"`
		Likes    []string `json:"likes"`
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
		key, err := db.Put(&User{
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
		key, err = db.Put(jimmy)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to put item: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("successfully put item with key: %s\n", key)
	}

More examples and complete documentation on https://docs.deta.sh/docs/base/sdk



*/
package base
