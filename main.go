package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"google.golang.org/api/idtoken"

	"github.com/nickolasgough/cloud-community-iam/internal/api"
	"github.com/nickolasgough/cloud-community-iam/internal/google"
	"github.com/nickolasgough/cloud-community-iam/internal/shared/constants"
)

const (
	PORT = 8000
)

func main() {
	ctx := context.Background()
	mux := http.NewServeMux()

	// Initialize environment data.
	gcpAuthClientID := os.Getenv(constants.GCP_OAUTH_CLIENT_ID_ENV_VAR)
	if gcpAuthClientID == "" {
		fmt.Printf("Failed to load GCP_OAUTH_CLIENT_ID environment variable\n")
		os.Exit(1)
	}

	// Initialize the Google ID token validator client
	googleIDVerifier, err := idtoken.NewValidator(ctx)
	if err != nil {
		fmt.Printf("Failed to initialize Google ID validator client")
		os.Exit(1)
	}
	googleService := google.NewService(gcpAuthClientID, googleIDVerifier)

	// Register API endpoints.
	mux.HandleFunc("/sign-in/with-google", api.SignInWithGoogle(ctx, googleService))

	fmt.Printf("Server listening on port %d\n", PORT)
	err = http.ListenAndServe(fmt.Sprintf(":%d", PORT), mux)
	if err != nil {
		os.Exit(1)
	}
}
