package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/chauhansantosh/doctor-appointment-manager/structure"
	"github.com/chauhansantosh/doctor-appointment-manager/util"
	"github.com/gin-gonic/gin"
)

// HandleLogin verifies user credentials
func HandleLogin(c *gin.Context) {
	var loginRequest structure.LoginRequest

	// Bind JSON request body to LoginRequest struct
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userExists, patientId := util.VerifyUserCredentials(loginRequest.Username, loginRequest.Password)
	// Verify user credentials
	if userExists {
		// Return a success response
		loginResponse := structure.LoginResponse{
			Success:   true,
			PatientId: patientId,
			Message:   "Authentication successful",
		}
		c.JSON(http.StatusOK, loginResponse)
	} else {
		// Return an error response indicating invalid credentials
		loginResponse := structure.LoginResponse{
			Success: false,
			Message: "Invalid credentials",
		}
		c.JSON(http.StatusUnauthorized, loginResponse)
	}
}

// handles the search doctor by location/specialization
func HandleSearchDoctors(c *gin.Context) {
	// Retrieve query parameters
	centerIdStr := c.Query("center_id")
	specializationIDStr := c.Query("specialization_id")
	centerId := 0
	specializationID := 0
	var err error

	if centerIdStr != "" {
		// Convert center_id to int
		centerId, err = strconv.Atoi(centerIdStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid center_id"})
			return
		}
	}

	if specializationIDStr != "" {
		// Convert specializationID to int
		specializationID, err = strconv.Atoi(specializationIDStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid specialization_id"})
			return
		}
	}

	// Perform doctor search based on location and specializationID
	searchResult, err := util.SearchDoctors(centerId, specializationID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	// Return the search result
	c.JSON(200, searchResult)
}

// HandleAppointmentBooking handles the appointment booking logic
func HandleAppointmentBooking(c *gin.Context) {
	var appointmentRequest structure.AppointmentRequest

	// Bind JSON request body to AppointmentRequest struct
	if err := c.ShouldBindJSON(&appointmentRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check the availability of the selected doctor at the specified date and time
	if util.IsAppointmentAvailable(appointmentRequest.DoctorID, appointmentRequest.AppointmentDate, appointmentRequest.AppointmentTime) {

		/* // Fetch the device token for the patient
		deviceToken, err := util.FetchDeviceToken(appointmentRequest.PatientID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch device token"})
			return
		} */
		// If available, create a new record in the BookedAppointments table
		bookedAppointment, err := util.BookAppointment(appointmentRequest)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to book appointment"})
			return
		}
		/* log.Printf("deviceToken: %v", deviceToken)
		// Send FCM notification to the patient
		err = util.SendAppointmentFCMNotification(deviceToken, "Appointment Booking", "Your appointment is booked!")
		if err != nil {
			// Handle FCM notification error (log or return an error)
			log.Printf("Error sending FCM notification for appointment booking: %v", err)
		} */
		// Return the booked appointment details
		c.JSON(http.StatusOK, bookedAppointment)
	} else {
		// Return an error response indicating unavailability
		c.JSON(http.StatusConflict, gin.H{"error": "Appointment slot not available"})
	}
}

// HandleAppointmentCancellation handles the appointment cancellation logic
func HandleAppointmentCancellation(c *gin.Context) {
	var cancellationRequest structure.CancellationRequest

	// Bind JSON request body to CancellationRequest struct
	if err := c.ShouldBindJSON(&cancellationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate cancellation request
	if err := util.ValidateCancellationRequest(cancellationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	/* // Fetch the device token for the patient
	deviceToken, err := util.FetchDeviceToken(cancellationRequest.PatientId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch device token"})
		return
	}
 */
	// Cancel the appointment
	if err := util.CancelAppointment(cancellationRequest.AppointmentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel appointment"})
		return
	}

	/* // Send FCM notification to the patient
	err = util.SendAppointmentFCMNotification(deviceToken, "Appointment Cancellation", "Your appointment is cancelled!")
	if err != nil {
		// Handle FCM notification error (log or return an error)
		log.Printf("Error sending FCM notification for appointment booking: %v", err)
	} */

	// Return a success message
	c.JSON(http.StatusOK, gin.H{"message": "Appointment cancellation successful"})
}

// HandleGetAllCenterLocations handles the GET request to fetch all center locations
func HandleGetAllCenterLocations(c *gin.Context) {
	centerLocations, err := util.GetAllCenterLocationsFromDB()
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, centerLocations)
}

// HandleGetAllSpecializations handles the GET request to fetch all Specializations
func HandleGetAllSpecializations(c *gin.Context) {
	specializations, err := util.GetAllSpecializationsFromDB()
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, specializations)
}

func HandleGetAvailableSlots(c *gin.Context) {
	// Retrieve query parameters
	doctorIdStr := c.Query("doctor_id")
	doctorId := 0
	var err error

	if doctorIdStr != "" {
		// Convert doctor_id to int
		doctorId, err = strconv.Atoi(doctorIdStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid doctor_id"})
			return
		}
	}

	// Get all available slots of a doctor based on doctorId
	availableSlots, err := util.GetAvailableSlotsFromDB(doctorId)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(200, availableSlots)
}

// HandleUpdatePatientDetails handles the update patient details endpoint
func HandleUpdatePatientDetails(c *gin.Context) {
	var request structure.PatientUpdateRequest

	patientIdStr := c.Query("patient_id")
	if patientIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "patient_id (query param) is required."})
		return
	}
	// Extract patient ID from the query parameters
	patientID, err := strconv.Atoi(c.Query("patient_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID"})
		return
	}

	// Bind JSON request body to PatientUpdateRequest struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert the request body to a JSON string
	requestJSON, err := json.Marshal(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal JSON"})
		return
	}

	// Unmarshal the JSON string to create the fieldsToUpdate map
	var fieldsToUpdate map[string]interface{}
	err = json.Unmarshal(requestJSON, &fieldsToUpdate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal JSON"})
		return
	}

	// Update patient details in the database
	err = util.UpdatePatientDetails(patientID, fieldsToUpdate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Failed to update patient details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Patient details updated successfully"})
}

// @Get /booked-appointments
func HandleGetBookedAppointments(c *gin.Context) {
	// Retrieve query parameters
	patientIdStr := c.Query("patient_id")
	patientId := 0
	var err error

	if patientIdStr != "" {
		// Convert patient_id to int
		patientId, err = strconv.Atoi(patientIdStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid patient_id"})
			return
		}
	}

	// Get all the booked appointments for a patient based on patient_id
	bookedAppointments, err := util.GetBookedAppointmentsFromDB(patientId)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(200, bookedAppointments)
}
