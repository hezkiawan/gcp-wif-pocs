package main

import (
	"context"
	"fmt"
	"log"
	"time"

	firebase "firebase.google.com/go/v4"
)

func main() {
	ctx := context.Background()

	// 1. Configure the exact GCP Project ID where Firestore lives.
	// Unlike static service account keys, WIF configs do not contain the Project ID,
	// so it must be explicitly declared in the application configuration.
	config := &firebase.Config{ProjectID: "vault-rotation-sandbox"}

	// 2. Initialize Firebase Admin SDK
	// The SDK heavily relies on standard environment variables. By setting
	// GOOGLE_APPLICATION_CREDENTIALS to point to our wif-config.json, the SDK
	// transparently handles reading the short-lived SPIFFE token.jwt and
	// exchanging it with the Google Security Token Service (STS) for temporary GCP access.
	app, err := firebase.NewApp(ctx, config)
	if err != nil {
		log.Fatalf("Error initializing firebase app: %v\n", err)
	}

	// 3. Initialize Firestore Client
	// If the WIF authentication flow above succeeded,
	// every Firestore request is automatically authorized
	// using the impersonated Service Account.
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Error initializing firestore client: %v\n", err)
	}
	defer client.Close()

	// 4. Perform a CRUD operation (Write a document)
	// Note that the application logic is completely agnostic to SPIFFE/SPIRE.
	// Zero hardcoded secrets are maintained in the codebase.
	fmt.Println("Attempting to write to Firestore using SPIFFE identity...")

	docRef, _, err := client.Collection("spire_poc").Add(ctx, map[string]interface{}{
		"message":   "Hello from a bare-metal Go app authenticated via SPIRE!",
		"timestamp": time.Now(),
		"spiffe_id": "spiffe://poc.local/go-app",
	})

	if err != nil {
		log.Fatalf("Failed adding document: %v\n", err)
	}

	fmt.Printf("✅ Success! Document written with ID: %v\n", docRef.ID)
	fmt.Println("🎉 The PoC is complete! Zero static secrets were used.")
}
