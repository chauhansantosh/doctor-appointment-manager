package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

var PostgresDb *sql.DB

func InitDB() {
	// Read environment variables
	projectID := os.Getenv("PROJECT_ID")
	secretName := os.Getenv("SECRET_NAME_DB")

	credentials, err := GetDBCredentials(projectID, secretName)
	fmt.Println("credentials: ", credentials)
	if err != nil {
		log.Fatalf("Error retrieving database credentials: %v", err)
	}

	// Use the retrieved credentials in your application
	host := credentials["host"]
	dbName := credentials["dbname"]
	dbUsername := credentials["user"]
	dbPassword := credentials["password"]
	sslmode := credentials["sslmode"]

	// Build the connection string
	connStr := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=%s",
		dbUsername, dbPassword, host, dbName, sslmode)

	fmt.Println("connStr ", connStr)

	// Open a connection to the database
	PostgresDb, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Test the database connection
	err = PostgresDb.Ping()
	if err != nil {
		log.Fatalf("Error pinging the database: %v", err)
	}

	log.Println("Connected to the database")
}

func GetDBCredentials(projectID, secretName string) (map[string]string, error) {
	var credentials map[string]string
	// Decode the base64-encoded payload
	decodedPayload, err := GetSecret(projectID, secretName)
	if err != nil {
		return credentials, err
	}
	// Decode the payload (assuming it's in JSON format)
	if err := json.Unmarshal([]byte(decodedPayload), &credentials); err != nil {
		return nil, fmt.Errorf("failed to decode secret payload: %v", err)
	}

	return credentials, nil
}

func GetSecret(projectID string, secretName string) (string, error) {
	ctx := context.Background()

	// Create the Secret Manager client
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create Secret Manager client: %v", err)
	}
	defer client.Close()

	// Build the secret version name
	secretVersionName := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, secretName)

	// Access the secret version
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretVersionName,
	}
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %v", err)
	}

	// Decode the secret payload
	return string(result.Payload.Data), nil
}
