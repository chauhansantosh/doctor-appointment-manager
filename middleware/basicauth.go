package middleware

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/chauhansantosh/doctor-appointment-manager/db"
	"github.com/gin-gonic/gin"
)

// BasicAuth is a middleware function that performs basic authentication
func BasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")

		// Check if the header is present
		if authHeader == "" {
			// No Authorization header provided, return 401 Unauthorized
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Check if the Authorization header starts with "Basic "
		if !strings.HasPrefix(authHeader, "Basic ") {
			// Invalid Authorization header, return 401 Unauthorized
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Decode the base64-encoded username and password
		credentialsBase64 := strings.TrimPrefix(authHeader, "Basic ")
		credentialsBytes, err := base64.StdEncoding.DecodeString(credentialsBase64)
		if err != nil {
			// Error decoding credentials, return 500 Internal Server Error
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Split the decoded credentials into username and password
		credentials := strings.Split(string(credentialsBytes), ":")
		if len(credentials) != 2 {
			// Malformed credentials, return 401 Unauthorized
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Check if the provided credentials are valid
		if password, ok := validCredentials()[credentials[0]]; !ok || password != credentials[1] {
			// Invalid credentials, return 401 Unauthorized
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Set the username in the request context for later use
		c.Set("client_user", credentials[0])

		// Continue to the next handler
		c.Next()
	}
}

// Credentials represents a pair of username and password
type Credentials struct {
	Username string `json:"client_user"`
	Password string `json:"client_pass"`
}

// validCredentials reads multiple sets of valid credentials from environment variables
func validCredentials() map[string]string {
	// Replace these with your GCP project ID and secret name
	projectID := os.Getenv("PROJECT_ID")
	secretName := os.Getenv("SECRET_NAME_BASIC_AUTH")

	// Fetch the secret from GCP Secret Manager
	credentialsJSON, err := db.GetSecret(projectID, secretName)
	if err != nil {
		log.Fatal("Error retrieving secret:", err)
	}

	// Parse JSON array into a slice of Credentials
	var credentialList []Credentials
	err = json.Unmarshal([]byte(credentialsJSON), &credentialList)
	if err != nil {
		log.Fatal("Error parsing credentials JSON:", err)
	}

	// Populate the credentials map
	credentials := make(map[string]string)
	for _, cred := range credentialList {
		credentials[cred.Username] = cred.Password
	}

	return credentials
}
