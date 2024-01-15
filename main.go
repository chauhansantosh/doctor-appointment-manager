// main.go

package main

import (
	"net/http"

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

	// Use CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} // Replace with your frontend URL
	config.AllowMethods = []string{"GET", "POST"}           // Adjust methods as needed
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
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
	router.POST("/patients/:patient_id", handler.HandleUpdatePatientDetails)
	router.GET("/booked-appointments", handler.HandleGetBookedAppointments)

	// Run the server
	router.Run(":8080")
}
