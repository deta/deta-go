/*
Package sdk is the official Deta SDK for Go.

SDK Packages

The SDK constitutes of two main components, the core package and service packages.

* `deta`: The core SDK package, provides shared functionalities to the service packages. All the `errors` are also exported from this package.

* `service`: The service packages, the services supported by the SDK.
	* `base`: Deta Base service package
	* `drive`: Deta Drive service package

Configuring credentials

When using the SDK you will require you project key. The project key can be provided explicitly or is taken from the environement variable `DETA_PROJECT_KEY`.

Default

By default, the SDK looks for the environment variable `DETA_PROJECT_KEY` for the project key

	// Create a new Deta instance taking the project key from the environment by default
	d, err := deta.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create new deta instance: %v\n", err)
	}

Provide the project key explicitly

You can use the `WithProjectKey` option when creating a `Deta` instance to provide the project key explicitly.

	// Create a new Deta instance with explicit project key
	d, err := deta.New(deta.WithProjectKey("project_key"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create new deta instance: %v\n", err)
	}

Examples

Base

The following is a simple `Put` operation example.

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

Drive

The following is a simple `Put` operation example.

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

More examples and complete documentation on https://docs.deta.sh/docs/drive/sdk/
*/

package sdk
