package structure

import "time"

// LoginRequest represents the request payload for user login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the structure of the login response.
type LoginResponse struct {
	Success   bool   `json:"success"`
	PatientId int    `json:"patient_id"`
	Message   string `json:"message"`
}

// Doctor represents a doctor entity
type Doctor struct {
	DoctorID         int    `json:"doctor_id"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Email            string `json:"email"`
	PhoneNumber      string `json:"phone_number"`
	Experience       int    `json:"experience"`
	Qualifications   string `json:"qualifications"`
	SpecializationID int    `json:"specialization_id"`
	CenterID         int    `json:"center_id"`
}

// SearchParams represents the request payload for doctor search
type SearchParams struct {
	Location         string `json:"location"`
	SpecializationID int    `json:"specialization_id"`
}

// AppointmentRequest represents the request payload for appointment booking
type AppointmentRequest struct {
	PatientID       int    `json:"patient_id" binding:"required"`
	DoctorID        int    `json:"doctor_id" binding:"required"`
	AppointmentDate string `json:"appointment_date" binding:"required"`
	AppointmentTime string `json:"appointment_time" binding:"required"`
}

// BookedAppointment represents the booked appointment
type BookedAppointment struct {
	AppointmentID   int       `json:"appointment_id"`
	PatientID       int       `json:"patient_id"`
	DoctorID        int       `json:"doctor_id"`
	AppointmentDate time.Time `json:"appointment_date"`
	AppointmentTime string    `json:"appointment_time"`
}

// CancellationRequest represents the request payload for appointment cancellation
type CancellationRequest struct {
	PatientId     int `json:"patient_id"`
	AppointmentID int `json:"appointment_id"`
}

// CenterLocation represents the data model for a center location
type CenterLocation struct {
	CenterID      int    `json:"center_id"`
	CenterName    string `json:"center_name"`
	Address       string `json:"address"`
	City          string `json:"city"`
	State         string `json:"state"`
	ZipCode       string `json:"zip_code"`
	ContactNumber string `json:"contact_number"`
}

// Specialization represents the data model for a specialization
type Specialization struct {
	SpecializationID   int    `json:"specialization_id"`
	SpecializationName string `json:"specialization_name"`
}

type AvailableDatesTime struct {
	AvailabilityId int       `json:"availability_id"`
	DoctorId       int       `json:"doctor_id"`
	Date           time.Time `json:"date"`
	TimeSlot       string    `json:"time_slot"`
}

// PatientUpdateRequest represents the JSON payload for updating patient details
type PatientUpdateRequest struct {
	FirstName    *string `json:"firstname"`
	LastName     *string `json:"lastname"`
	Email        *string `json:"email"`
	PhoneNumber  *string `json:"phonenumber"`
	EmployeeCode *string `json:"employeecode"`
	HCLEmail     *string `json:"hclemail"`
	Password     *string `json:"password"`
	DeviceToken  *string `json:"device_token"`
}
