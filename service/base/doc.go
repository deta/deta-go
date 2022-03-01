/*
Package base is the Deta Base service package.

The following is a simple Put operation example.

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

More examples and complete documentation on https://docs.deta.sh/docs/base/sdk



*/
package base
