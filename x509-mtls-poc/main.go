package main

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
)

func main() {
	ctx := context.Background()

	// 1. Explicitly Define the Project ID
	// Unlike static service account JSON keys, external WIF configurations
	// do not contain the GCP Project ID. We must explicitly inject it into the SDK.
	config := &firebase.Config{ProjectID: "vault-rotation-sandbox"}

	// 2. Initialize the Firebase Application via Application Default Credentials (ADC)
	// The SDK automatically reads the GOOGLE_APPLICATION_CREDENTIALS environment variable.
	// It parses wif-config.json, locates the workload certificate on disk,
	// and performs the mTLS handshake with the Google Security Token Service (STS).
	app, err := firebase.NewApp(ctx, config)
	if err != nil {
		log.Fatalf("Error initializing app: %v\n", err)
	}

	// 3. Initialize the Firestore Client
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Error initializing Firestore client: %v\n", err)
	}
	defer client.Close()

	// 4. Execute the Database Write
	// This proves that GCP accepted our self-signed certificate, successfully mapped
	// the Common Name (my-local-vm) to the IAM Service Account, and granted write access.
	_, _, err = client.Collection("x509_sim_poc").Add(ctx, map[string]interface{}{
		"message": "Hello from the X.509 mTLS Certificate simulation!",
		"auth":    "Secretless WIF - Sectigo Simulation",
	})
	if err != nil {
		log.Fatalf("Failed adding document: %v\n", err)
	}

	fmt.Println("Success! Document written to Firestore using an X.509 client certificate.")
}
