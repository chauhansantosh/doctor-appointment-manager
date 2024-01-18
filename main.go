// main.go

package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/chauhansantosh/doctor-appointment-manager/db"
	"github.com/chauhansantosh/doctor-appointment-manager/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// Initialize the PostgreSQL database connection
	db.InitDB()

	// Create a new Gin router
	router := gin.Default()

	// Read configuration from environment variables
	allowedOrigins := os.Getenv("CORS_ALLOW_ORIGINS")
	allowedMethods := os.Getenv("CORS_ALLOW_METHODS")
	allowedHeaders := os.Getenv("CORS_ALLOW_HEADERS")

	// Use CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = strings.Split(allowedOrigins, ",")
	config.AllowMethods = strings.Split(allowedMethods, ",")
	config.AllowHeaders = strings.Split(allowedHeaders, ",")
	router.Use(cors.New(config))

	// Define routes
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to Doctor Appointment Manger"})
	})
	// Authentication routes
	router.POST("/login", handler.HandleLogin)
	router.GET("/doctors", handler.HandleSearchDoctors)
	// Define the appointment booking endpoint
	router.POST("/book-appointment", handler.HandleAppointmentBooking)
	// Define the appointment cancellation endpoint
	router.POST("/cancel-appointment", handler.HandleAppointmentCancellation)

	// Define the GET endpoint for fetching all center locations
	router.GET("/center-locations", handler.HandleGetAllCenterLocations)
	// Define the GET endpoint for fetching all specializations
	router.GET("/specializations", handler.HandleGetAllSpecializations)
	router.GET("/available-time-slots", handler.HandleGetAvailableSlots)
	router.POST("/update-patient", handler.HandleUpdatePatientDetails)
	router.GET("/booked-appointments", handler.HandleGetBookedAppointments)

	// Run the server
	router.Run(":8080")
}
