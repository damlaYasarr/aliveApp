package utils

import (
	"context"
	"path/filepath"
"fmt"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

// SetupFirebase initializes Firebase App and Messaging client
func SetupFirebase() (*firebase.App, context.Context, *messaging.Client) {
	ctx := context.Background()

	// Path to your service account key JSON file
	serviceAccountKeyFilePath, err := filepath.Abs("./serviceAccountKey.json")
	if err != nil {
		panic("Unable to load serviceAccountKeys.json file")
	}

	// Initialize Firebase app with service account key
	opt := option.WithCredentialsFile(serviceAccountKeyFilePath)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		panic(fmt.Errorf("error initializing Firebase app: %v", err))
	}

	// Retrieve Messaging client
	client, err := app.Messaging(ctx)
	if err != nil {
		panic(fmt.Errorf("error getting Messaging client: %v", err))
	}

	return app, ctx, client
}