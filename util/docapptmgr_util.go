package util

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	messaging "firebase.google.com/go/v4/messaging"
	"github.com/chauhansantosh/doctor-appointment-manager/structure"
)

// ValidateCancellationRequest validates the cancellation request
func ValidateCancellationRequest(request structure.CancellationRequest) error {
	//appointment ID must be greater than 0
	if request.AppointmentID <= 0 {
		return fmt.Errorf("invalid appointment ID=%d. It should be greater than 0", request.AppointmentID)
	}

	return nil
}

// SendAppointmentFCMNotification sends an FCM notification to the specified patient
func SendAppointmentFCMNotification(deviceToken string, title, message string) error {
	/* // Get the path to the service account file from environment variable
	credentialsPath := os.Getenv("FIREBASE_CREDENTIALS_PATH")
	fmt.Println("FIREBASE_CREDENTIALS_PATH", credentialsPath)
	if credentialsPath == "" {
		log.Fatal("FIREBASE_CREDENTIALS_PATH environment variable is not set.")
	}
	// Initialize the Firebase Admin SDK with the specified service account file
	opt := option.WithCredentialsFile(credentialsPath) */
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Printf("Error initializing Firebase app: %v", err)
		return err
	}

	// Get a messaging client from the Firebase app
	client, err := app.Messaging(context.Background())
	if err != nil {
		log.Printf("Error creating FCM client: %v", err)
		return err
	}

	// Create a new FCM message
	msg := &messaging.Message{
		Data: map[string]string{
			"title":   title,
			"message": message,
		},
		Token: deviceToken,
	}

	// Send the FCM message
	response, err := client.Send(context.Background(), msg)
	if err != nil {
		// Handle FCM sending error (log or return an error)
		log.Printf("Error sending FCM notification: %v", err)
		return err
	}

	// Log the FCM response
	log.Printf("FCM response: %+v\n", response)

	return nil
}
