# doctor-appointment-manager
System for Booking Doctor’s Appointment

 This project involves developing the logic for user authentication, doctor search, appointment booking, cancellation, and notifications. Additionally, it implements API endpoints for frontend-backend communication and used Cloud Endpoints wtih Google API Gateway to expose these APIs securely.

### Backend Logic Implementation:

1. **User Authentication:**
   - Validated user credentials (HCL-Employee Code) against the `Patients` table in PostgreSQL database.

2. **Doctor Search:**
   - Created a Golang function to handle doctor searches based on location and specialization.
   - Queries the `Doctors` and `Specializations` tables in PostgreSQL database to retrieve relevant information.

3. **Appointment Booking:**
   - Implemented logic to allow patients to book appointments.
   - Checked the availability of the selected doctor at the specified date and time using `AvailableDatesTime`and `BookedAppointments` tables.
   - If available, create a new record in the `BookedAppointments` table.

4. **Appointment Cancellation:**
   - Allowed patients to cancel their booked appointments.
   - Validated cancellation requests and delete the corresponding record from the `BookedAppointments` table.

5. **Notification Logic:**
   - Implemented logic to send push notifications to patient'd web browser for appointment confirmations and cancellations.
   - Used Firebase Cloud Messaging (FCM) for push notifications.

### API Endpoints:

1. **Created API Endpoints:**
   - Defined API endpoints in Golang application for user authentication, doctor search, appointment booking, and cancellation.
   - Used a gin-gonic package to handle HTTP routing.

2. **Secured Endpoints:**
   - Implemented security measures for API endpoints using API Key.

### Cloud Endpoints:

1. **Created OpenAPI Specification:**
   - Defined an OpenAPI Specification (formerly known as Swagger) for APIs. This document describes the structure of API, including endpoints, request/response formats, and authentication requirements.

2. **Created API config and API Gateway:**
   - Created api config from OpenAPI Specifications using Google Cloud Endpoints.
   - Created API Gateway using created API config.
   - API Gateway allows us to manage and secure the API easily.

### Cloud Functions:

  **Created Cloud Function to send push notification:**
   - Created a cloud function to send notfication to any device with device token and deployed to Google cloud functions.
   - Invoked cloud function on successful booking or cancellation by patient.


### doctor-appointment-manager api endpoints with request/response.

#### End-point: search doctors
##### Method: GET
>```
>/doctors?center_id=1&specialization_id=1
>```
##### Query Params

|Param|value|
|---|---|
|center_id|1|
|specialization_id|1|


##### Response: 200
```json
[
    {
        "doctor_id": 1,
        "first_name": "Michael",
        "last_name": "Johnson",
        "email": "michael.johnson@example.com",
        "phone_number": "1112223333",
        "experience": 12,
        "qualifications": "MD",
        "specialization_id": 1,
        "center_id": 1
    }
]
```


⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

#### End-point: get all center-locations
##### Method: GET
>```
>/center-locations
>```
##### Body (**raw**)

```json

```

##### Response: 200
```json
[
    {
        "center_id": 1,
        "center_name": "Hospital A",
        "address": "123 Main St",
        "city": "City A",
        "state": "State A",
        "zip_code": "12345",
        "contact_number": "1111111111"
    },
    {
        "center_id": 2,
        "center_name": "Clinic B",
        "address": "456 Oak St",
        "city": "City B",
        "state": "State B",
        "zip_code": "67890",
        "contact_number": "2222222222"
    },
    {
        "center_id": 3,
        "center_name": "Medical Center C",
        "address": "789 Pine St",
        "city": "City C",
        "state": "State C",
        "zip_code": "45678",
        "contact_number": "3333333333"
    }
]
```


⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

#### End-point: get all specializations
##### Method: GET
>```
>/specializations
>```
##### Body (**raw**)

```json

```

##### Response: 200
```json
[
    {
        "specialization_id": 1,
        "specialization_name": "Cardiology"
    },
    {
        "specialization_id": 2,
        "specialization_name": "Dermatology"
    },
    {
        "specialization_id": 3,
        "specialization_name": "Pediatrics"
    }
]
```


⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

#### End-point: get booked appointments
##### Method: GET
>```
>/booked-appointments?patient_id=1
>```
##### Query Params

|Param|value|
|---|---|
|patient_id|1|


##### Response: 200
```json
[
    {
        "appointment_id": 20,
        "patient_id": 1,
        "doctor_id": 1,
        "appointment_date": "2023-03-01T00:00:00Z",
        "appointment_time": "Morning"
    }
]
```


⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

#### End-point: get available slots
##### Method: GET
>```
>/available-time-slots?doctor_id=1
>```
##### Query Params

|Param|value|
|---|---|
|doctor_id|1|


##### Response: 200
```json
[
    {
        "availability_id": 1,
        "doctor_id": 1,
        "date": "2023-03-01T00:00:00Z",
        "time_slot": "Morning"
    }
]
```


⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

#### End-point: login
##### Method: POST
>```
>/login
>```
##### Body (**raw**)

```json
{
    "username": "EMP001",
    "password": "password123"
}
```

##### Response: 200
```json
{
    "success": true,
    "patient_id": 1,
    "message": "Authentication successful"
}
```


⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

#### End-point: update patient
##### Method: POST
>```
>/update-patient?patient_id=1
>```
##### Body (**raw**)

```json
{
    "firstname":"John",
    "lastname": "Doe",
    "email":"john.doe@example.com",
    "phonenumber":"+1234567890",
    "employeecode":"EMP001",
    "hclemail":"john.doe@hcl.com",
    "password":"password123"
}
```

##### Query Params

|Param|value|
|---|---|
|patient_id|1|


##### Response: 200
```json
{
    "message": "Patient details updated successfully",
    "success": true
}
```


⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

#### End-point: book appointment error: slot not avaiable
##### Method: POST
>```
>/book-appointment
>```
##### Body (**raw**)

```json
{
    "patient_id": 1,
    "doctor_id": 6,
    "appointment_date": "2023-03-01",
    "appointment_time": "Afternoon"
}
```

##### Response: 409
```json
{
    "error": "Appointment slot not available"
}
```


⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

#### End-point: book appointment
##### Method: POST
>```
>/book-appointment
>```
##### Body (**raw**)

```json
{
    "patient_id": 1,
    "doctor_id": 1,
    "appointment_date": "2023-03-01",
    "appointment_time": "Morning"
}
```

##### Response: 200
```json
{
    "appointment_id": 21,
    "patient_id": 1,
    "doctor_id": 1,
    "appointment_date": "2023-03-01T00:00:00Z",
    "appointment_time": "Morning"
}
```


⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

#### End-point: cancel appointment error: invalid appointment id
##### Method: POST
>```
>/cancel-appointment
>```
##### Body (**raw**)

```json
{
    "patient_id": 1,
    "appointment_id": 1
}
```

##### Response: 400
```json
{
    "error": "invalid appointment ID=0. It should be greater than 0"
}
```


⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

#### End-point: cancel appointment success
##### Method: POST
>```
>/cancel-appointment
>```
##### Body (**raw**)

```json
{
    "patient_id": 1,
    "appointment_id": 21
}
```

##### Response: 200
```json
{
    "message": "Appointment cancellation successful"
}
```


⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

#### End-point: send-message cloud function
##### Method: POST
>```
>https://my-api-gateway-cvai9ad9.an.gateway.dev/send-message
>```
##### Headers

|Content-Type|Value|
|---|---|
|x-api-key|your-ap-key|


##### Body (**raw**)

```json
{
    "device_token": "cnPBcWCil6yGNhL6OL-fPf:APA91bETOWHGNdadJ7mDMXS9iyDachIFKLB4D7nuFrQoq088tU3dLeduoeuMn19C4QcEK3CRE5vGa8n1Z46DEN7WRqOsPvjpv4lOKzljbRva9cCU271kHQLfvL_yf_8LB4e3a3fB_iz6",
    "title": "test cloud function",
    "message": "message sent successfully"
}
```

##### Response: 200
```json
Notification sent successfully

```


⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃
