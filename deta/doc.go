/*
Package deta is the core Deta SDK package.

Configuring credentials

You will require you project key when using the Deta SDK.

By default, the SDK looks for the environment variable DETA_PROJECT_KEY for the project key.

	// Create a new Deta instance taking the project key from the environment by default
	d, err := deta.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create new deta instance: %v\n", err)
	}

You can use the WithProjectKey option when creating a Deta instance to provide the project key explicitly.

	// Create a new Deta instance with explicit project key
	d, err := deta.New(deta.WithProjectKey("project_key"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create new deta instance: %v\n", err)
	}
*/
package deta
