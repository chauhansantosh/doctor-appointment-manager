package util

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chauhansantosh/doctor-appointment-manager/db"
	"github.com/chauhansantosh/doctor-appointment-manager/structure"
)

// VerifyUserCredentials queries the database to verify user credentials
func VerifyUserCredentials(username, password string) (bool, int) {
	// Query the database to verify user credentials
	query := "SELECT patientid FROM patients WHERE EmployeeCode = $1 AND password = $2"
	var patientID int

	// Execute the query
	err := db.PostgresDb.QueryRow(query, username, password).Scan(&patientID)
	if err != nil {
		log.Printf("Error querying the database: %v", err)
		return false, 0
	}

	// If patientID is greater than 0, credentials are valid
	return patientID > 0, patientID
}

// SearchDoctors performs a doctor search based on location and specialization
func SearchDoctors(location, specializationID int) ([]structure.Doctor, error) {
	result := []structure.Doctor{}

	// Query the database for doctors based on location and specialization
	query := `
		SELECT D.DoctorID, D.FirstName, D.LastName, D.Email, D.PhoneNumber, D.Experience, D.Qualifications, D.SpecializationID, D.CenterID
		FROM Doctors D
		WHERE (D.CenterID = $1 OR $1 = 0) AND (D.SpecializationID = $2 OR $2 = 0)
	`

	rows, err := db.PostgresDb.Query(query, location, specializationID)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	// Iterate through the result set and populate the SearchResult
	for rows.Next() {
		var doctor structure.Doctor
		err := rows.Scan(&doctor.DoctorID, &doctor.FirstName, &doctor.LastName, &doctor.Email,
			&doctor.PhoneNumber, &doctor.Experience, &doctor.Qualifications, &doctor.SpecializationID, &doctor.CenterID)
		if err != nil {
			return result, err
		}
		result = append(result, doctor)
	}

	return result, nil
}

// IsAppointmentAvailable checks the availability of the selected doctor at the specified date and time
func IsAppointmentAvailable(doctorID int, appointmentDate string, appointmentTime string) bool {
	// Query the AvailableDatesTime table to identify potential slots
	availableQuery := `
		SELECT COUNT(*) 
		FROM AvailableDatesTime 
		WHERE DoctorID = $1 
			AND Date = $2 
			AND TimeSlot = $3
	`

	// Query the BookedAppointments table to check if the slot is already booked
	bookedQuery := `
		SELECT COUNT(*) 
		FROM BookedAppointments
		WHERE DoctorID = $1 
			AND AppointmentDate = $2 
			AND AppointmentTime = $3
	`

	var availableCount, bookedCount int

	// Check potential slot availability in AvailableDatesTime table
	err := db.PostgresDb.QueryRow(availableQuery, doctorID, appointmentDate, appointmentTime).Scan(&availableCount)
	if err != nil {
		// Handle the error (e.g., log or return false)
		fmt.Println("Error checking potential slot availability:", err)
		return false
	}

	// Check if the slot is already booked in BookedAppointments table
	err = db.PostgresDb.QueryRow(bookedQuery, doctorID, appointmentDate, appointmentTime).Scan(&bookedCount)
	if err != nil {
		// Handle the error (e.g., log or return false)
		fmt.Println("Error checking booked slot availability:", err)
		return false
	}

	// If there are potential slots and none of them are booked, then the slot is available
	return availableCount > 0 && bookedCount == 0
}

// BookAppointment creates a new record in the BookedAppointments table
func BookAppointment(appointmentRequest structure.AppointmentRequest) (structure.BookedAppointment, error) {
	// Prepare the SQL statement for inserting into BookedAppointments
	query := `
		INSERT INTO BookedAppointments (PatientID, DoctorID, AppointmentDate, AppointmentTime)
		VALUES ($1, $2, $3, $4)
		RETURNING AppointmentID, AppointmentDate
	`

	// Execute the SQL statement and retrieve the generated AppointmentID and AppointmentDate
	var appointmentID int
	var appointmentDate time.Time
	err := db.PostgresDb.QueryRow(query, appointmentRequest.PatientID, appointmentRequest.DoctorID, appointmentRequest.AppointmentDate, appointmentRequest.AppointmentTime).
		Scan(&appointmentID, &appointmentDate)
	if err != nil {
		// Handle the error (e.g., log or return an error)
		return structure.BookedAppointment{}, fmt.Errorf("failed to book appointment: %v", err)
	}

	// Create a BookedAppointment struct with the generated AppointmentID and AppointmentDate
	bookedAppointment := structure.BookedAppointment{
		AppointmentID:   appointmentID,
		PatientID:       appointmentRequest.PatientID,
		DoctorID:        appointmentRequest.DoctorID,
		AppointmentDate: appointmentDate,
		AppointmentTime: appointmentRequest.AppointmentTime,
	}

	return bookedAppointment, nil
}

// CancelAppointment deletes the corresponding record from the BookedAppointments table
func CancelAppointment(appointmentID int) error {
	query := `
		DELETE FROM BookedAppointments
		WHERE AppointmentID = $1
	`

	// Execute the SQL statement to delete the appointment and get the affected rows
	result, err := db.PostgresDb.Exec(query, appointmentID)
	if err != nil {
		fmt.Println("failed to cancel appointment:", err.Error())
		// Handle the error (e.g., log or return an error)
		return fmt.Errorf("failed to cancel appointment: %v", err.Error())
	}

	// Check the number of affected rows
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// Handle the error (e.g., log or return an error)
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	// If no rows were affected, return an error
	if rowsAffected == 0 {
		return fmt.Errorf("no appointment found with ID %d", appointmentID)
	}

	return nil
}

// GetAllCenterLocationsFromDB fetches all center locations from the database
func GetAllCenterLocationsFromDB() ([]structure.CenterLocation, error) {
	// Query to fetch all center locations
	rows, err := db.PostgresDb.Query("SELECT * FROM CenterLocations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	centerLocations := []structure.CenterLocation{}

	// Iterate through the result set and populate centerLocations
	for rows.Next() {
		var centerLocation structure.CenterLocation
		if err := rows.Scan(&centerLocation.CenterID, &centerLocation.CenterName, &centerLocation.Address,
			&centerLocation.City, &centerLocation.State, &centerLocation.ZipCode, &centerLocation.ContactNumber); err != nil {
			return nil, err
		}
		centerLocations = append(centerLocations, centerLocation)
	}

	return centerLocations, nil
}

// GetAllSpecializationsFromDB fetches all specializations from the database
func GetAllSpecializationsFromDB() ([]structure.Specialization, error) {
	// Query to fetch all specializations
	rows, err := db.PostgresDb.Query("SELECT * FROM specializations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	specializations := []structure.Specialization{}

	// Iterate through the result set and populate specializations
	for rows.Next() {
		var specialization structure.Specialization
		if err := rows.Scan(&specialization.SpecializationID, &specialization.SpecializationName); err != nil {
			return nil, err
		}
		specializations = append(specializations, specialization)
	}

	return specializations, nil
}

func GetAvailableSlotsFromDB(doctorId int) ([]structure.AvailableDatesTime, error) {
	result := []structure.AvailableDatesTime{}

	query := `
		SELECT availabilityid, doctorid, date, timeslot
		FROM availabledatestime
		WHERE (doctorid = $1 OR $1 = 0)
	`

	rows, err := db.PostgresDb.Query(query, doctorId)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	// Iterate through the result set and populate the SearchResult
	for rows.Next() {
		var availableSlot structure.AvailableDatesTime
		err := rows.Scan(&availableSlot.AvailabilityId, &availableSlot.DoctorId, &availableSlot.Date, &availableSlot.TimeSlot)
		if err != nil {
			return result, err
		}
		result = append(result, availableSlot)
	}

	return result, nil
}

// Function to update patient details in the database
func UpdatePatientDetails(patientId int, fieldsToUpdate map[string]interface{}) error {
	// Build the dynamic SQL query to update patient details
	var queryBuilder strings.Builder
	queryBuilder.WriteString("UPDATE patients SET ")

	// Prepare SET clause based on provided fields
	var params []interface{}
	var setClauses []string

	params = append(params, patientId)
	count := 1
	for field, value := range fieldsToUpdate {
		fmt.Println(field, value)
		if value != nil {
			fmt.Println("value!=nil", field, value)
			count = count + 1
			setClauses = append(setClauses, fmt.Sprintf("%s = $%d", field, count))
			params = append(params, value)
		}
	}

	queryBuilder.WriteString(strings.Join(setClauses, ", "))
	queryBuilder.WriteString(" WHERE patientid = $1")

	query := queryBuilder.String()

	fmt.Println("query:", query, "params:", params)

	// Execute the SQL query
	_, err := db.PostgresDb.Exec(query, params...)

	return err
}

// FetchDeviceToken retrieves the device token for a given patient ID
func FetchDeviceToken(patientID int) (string, error) {
	var deviceToken string
	err := db.PostgresDb.QueryRow("SELECT device_token FROM patients WHERE patientid = $1", patientID).Scan(&deviceToken)
	if err != nil {
		return "", err
	}
	return deviceToken, nil
}

func GetBookedAppointmentsFromDB(doctorId int) ([]structure.BookedAppointment, error) {
	result := []structure.BookedAppointment{}

	query := `
		SELECT appointmentid, patientid, doctorid, appointmentdate, appointmenttime
		FROM bookedappointments
		WHERE (patientid = $1 OR $1 = 0)
	`

	rows, err := db.PostgresDb.Query(query, doctorId)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	// Iterate through the result set and populate the SearchResult
	for rows.Next() {
		var bookedAppointments structure.BookedAppointment
		err := rows.Scan(&bookedAppointments.AppointmentID, &bookedAppointments.PatientID, &bookedAppointments.DoctorID,
			&bookedAppointments.AppointmentDate, &bookedAppointments.AppointmentTime)
		if err != nil {
			return result, err
		}
		result = append(result, bookedAppointments)
	}

	return result, nil
}
