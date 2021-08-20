/*
Package drive is the Deta Drive service package.

The following is a simple Put operation example.

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
package drive
