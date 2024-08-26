# CRM Application Backend

## Table of Contents
- [Introduction](#introduction)
- [Project Structure](#project-structure)
- [Setup and Running Locally](#setup-and-running-locally)
- [System Design](#system-design)
  - [User Model](#user-model)
- [API Endpoints](#api-endpoints)
- [Email Integration](#email-integration)
- [Analytics and Reporting](#analytics-and-reporting)
- [Contributing](#contributing)

## Introduction
This backend application serves as the core of a CRM system, managing users, interactions, leads, and more. It is built using Go with the Gin framework and MongoDB for data storage. The application supports user authentication, role based support, email notifications, and interaction tracking.

## Project Structure
```
crm-backend/
├── controllers/ │ 
   ├── userController.go │ 
   ├── interactionController.go │ 
   ├── leadController.go │ 
   └── emailController.go 
├── middleware/ │ 
   ├── authMiddleware.go │ 
   └── ... 
├── models/ │ 
   ├── userModel.go │
   ├── customerModel.go │
   ├── comapnyModel.go │ 
   ├── interactionModel.go │ 
   ├── leadModel.go │ 
   └── ... 
├── routes/ │ 
   ├── userRoutes.go │
   ├── comapnyRoutes.go │
   ├── customerRoutes.go │ 
   ├── interactionRoutes.go │ 
   ├── leadRoutes.go │ 
   ├── authRoutes.go │
   └── emailRoutes.go 
├── helpers/ │ 
   ├── emailHelper.go │ 
   ├── authHelper.go │ 
   ├── tokenHelper.go │ 
   └── ... 
├── .env 
├── main.go 
└── go.mod
```

## Setup and Running Locally
### Prerequisites
- Go and MongoDB instance running

### Steps
1. **Clone the Repository**
   ```bash
   git clone https://github.com/SiddharthaKR/matriceCRM
   cd matriceCRM
   ```

2. **Install Dependencies**
   ```bash
   go mod tidy
   ```

3. **Set Up Environment Variables**
   Create a `.env` file in the root directory with the following content:
   ```plaintext
   MONGO_URI=mongodb://localhost:27017
   SMTP_ADDR=smtp.example.com:587
   FROM_EMAIL=your-email@example.com
   FROM_EMAIL_PASSWORD=your-password
   ```

4. **Run the Application**
   ```bash
   go run main.go
   ```

5. **Backend is now running at** `http://localhost:9000`


# Key Components

## User Authentication
- **Login and Registration:**
  - **Login**: Users can log in to the CRM system using their credentials (email and password).
  - **Registration**: New users can register by providing necessary information, including name, email, and password.

- **Token-Based Authentication:**
  - **JWT (JSON Web Tokens)**: The system uses JWT for secure authentication. Upon successful login, a JWT token is issued and used for subsequent requests to verify user identity.
  - **Session Management**: Handle user sessions with tokens for access to protected routes and functionalities.

- **Password Management:**
  - **Password Hashing**: User passwords are hashed and stored securely to prevent unauthorized access.
  - **Password Reset**: Users can request to reset their password if forgotten, usually via email.

## User Management
- **CRUD Operations for Customers and Users:**
  - **Create**: Add new users and customers to the CRM system.
  - **Read**: Retrieve and view details about existing users and customers.
  - **Update**: Modify information about users and customers.
  - **Delete**: Remove users and customers from the system.
  
- **Details for Each User/Customer:**
  - **Name**: First and last names of the user or customer.
  - **Contact Information**: Includes email, phone number, and other relevant contact details.
  - **Company**: Company associated with the user or customer, if applicable.
  - **Status**: Current status such as active, inactive, or pending.
  - **Notes**: Additional information or comments related to the user or customer.

## Interaction Tracking
- **Ticket Management:**
  - **Raise a Ticket**: Customers can create support tickets for issues or requests.
  - **Mark as Resolved**: Customers or users can update the status of a ticket to resolved.

- **Scheduling Interactions:**
  - **Meetings**: Users can schedule meetings with customers.
  - **Link Interactions**: Ensure each interaction (e.g., meeting) is associated with the correct customer.

- **Interaction History:**
  - **View History**: Users can access and review the interaction history for a specific customer.

## Analytics and Reporting (Optional)
- **Report Generation:**
  - **Customer/Lead Interactions**: Generate reports on interactions between users and customers or leads.
  - **Conversion Rates**: Track and report on conversion rates of leads to customers.

- **Visual Representations:**
  - **Charts and Graphs**: Include visual tools such as charts and graphs to illustrate key metrics and trends.

## Email Integration (Optional)
- **Email Sending:**
  - **Direct Sending**: Allow users to send emails directly from the CRM.

- **Tracking:**
  - **Open Rates**: Monitor and track how often emails are opened.
  - **Responses**: Track replies and interactions with sent emails.

## Activity Notifications
- **Notification System:**
  - **Alerts**: Notify users about upcoming tasks, meetings, and follow-ups to ensure timely actions and management.

## Role-Based Access Control
- **Access Control:**
  - **Roles**: Implement role-based access control to manage permissions and access rights within the CRM.
  - **Role Types**: Define roles such as ADMIN, MANAGER, and USER with varying levels of access and functionality.
  
  - **Authorization**: Ensure users have access only to the features and data that are appropriate for their role.


# API Route Documentation

## 1. User Signup
**Endpoint:** `POST http://localhost:9000/users/signup`

**Request Body:**
```json
{
  "first_name": "John",
  "last_name": "Cena",
  "password": "hellbound",
  "email": "admin@gmail.com",
  "phone": "8699195788",
  "user_type": "ADMIN"
}
```
## 2. User Login
**Endpoint:** `POST http://localhost:9000/users/login`

**Request Body:**
```json
{
  "password": "hellbound",
  "email": "john@gmail.com"
}
```
## 3. Customer Signup
**Endpoint:** `http://localhost:9000/customers/signup`

**Request Body:**
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "johndoe@example.com",
  "phone": "+1234567890",
  "company": "Acme Corporation",
  "status": "LEAD",
  "password": "gitbash",
  "notes": "Interested in our premium products.",
  "created_at": "2024-08-25T10:00:00Z",
  "updated_at": "2024-08-25T10:00:00Z",
  "customer_id": "CUST12345",
  "last_interaction": "2024-08-24T15:00:00Z",
  "company_id": "66cb704a8477922b84d2481f"
}
```
## 4. Customer Login
**Endpoint:** `http://localhost:9000/customers/login`

**Request Body:**
```json
{
  "email": "johndoe@example.com",
  "password": "gitbash"
}
```
## 5. Get All Users
**Endpoint:** `GET http://localhost:9000/users`

**Description:** Retrieves a paginated list of all users. This route requires an `ADMIN` token for access.

**Request Headers:**
- `token: <token>`

**Query Parameters:**
- `recordPerPage` (optional): Number of records per page (default is 10)
- `page` (optional): Page number for pagination (default is 1)

**Response:**
```json
{
  "total_count": 100,
  "user_items": [
    {
      "user_id": "user123",
      "first_name": "John",
      "last_name": "Doe",
      "email": "john.doe@example.com",
      "phone": "1234567890",
      "status": "ACTIVE",
      "user_type": "USER",
      "notes": "Some notes here"
    },
    // more user objects
  ]
}
```
## 6. Get User
**Endpoint:** `GET http://localhost:9000/users/:user_id`

**Description:** Retrieves details of a specific user by their ID. This route requires authentication and authorization to access user details.

**Request Headers:**
- `token: <token>`

**URL Parameters:**
- `user_id`: ID of the user to retrieve

**Response:**
```json
{
  "user_id": "user123",
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "phone": "1234567890",
  "status": "ACTIVE",
  "user_type": "USER",
  "notes": "Some notes here"
}
```
## 7. Update User
**Endpoint:** `PUT http://localhost:9000/users/:user_id`

**Description:** Updates details of a specific user. This route requires authentication and authorization to update user details.

**Request Headers:**
- `token: <token>`

**URL Parameters:**
- `user_id`: ID of the user to update

**Request Body:**
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "phone": "1234567890",
  "status": "ACTIVE",
  "notes": "Updated notes",
  "user_type": "USER"
}
```
##8. Delete User
**Endpoint:** `DELETE http://localhost:9000/users/:user_id`

**Description:** Deletes a specific user by their ID. This route requires an ADMIN token for access.

**Request Headers:**
- `token: <token>`

**URL Parameters:**
- `user_id`: ID of the user to delete

**Response:**
```json
{
  "message": "User deleted successfully"
}
```

## 9. Create Company
**Endpoint:** `POST http://localhost:9000/companies`

**Description:** Creates a new company. This route requires authentication.

**Request Headers:**
- `token: <token>`

**Request Body:**
```json
{
  "name": "Acme Corporation",
  "address": "1234 Elm Street",
  "phone": "+1234567890",
  "email": "info@acme.com",
  "created_at": "2024-08-25T10:00:00Z",
  "updated_at": "2024-08-25T10:00:00Z"
}
```
## Get Company
**Endpoint:** `GET http://localhost:9000/companies/:company_id`

**Description:** Retrieves details of a specific company by its ID. Requires authentication and authorization.

**Request Headers:**
- `token: <token>`

**URL Parameters:**
- `company_id`: ID of the company to retrieve

**Response:**
```json
{
  "id": "60f7e3a4b9f1b2c6d8e4f4b0",
  "name": "Acme Corporation",
  "address": "1234 Elm Street",
  "phone": "+1234567890",
  "email": "info@acme.com",
  "created_at": "2024-08-25T10:00:00Z",
  "updated_at": "2024-08-25T10:00:00Z"
}

```
## Update Company
**Endpoint:** `PUT http://localhost:9000/companies/:company_id`

**Description:** Updates details of a specific company. Requires authentication and authorization.

**Request Headers:**
- `token: <token>`

**URL Parameters:**
- `company_id`: ID of the company to update

**Request Body:**
```json
{
  "name": "Acme Corporation",
  "address": "5678 Oak Street",
  "phone": "+0987654321",
  "email": "contact@acme.com",
  "updated_at": "2024-08-25T10:00:00Z"
}
```

## Delete Company
**Endpoint:** `DELETE http://localhost:9000/companies/:company_id`

**Description:** Deletes a specific company by its ID. Requires an ADMIN token for access.

**Request Headers:**
- `token: <token>`

**URL Parameters:**
- `company_id`: ID of the company to delete

**Response:**
```json
{
  "message": "Company deleted successfully",
  "deleted_count": 1
}
```
## Get All Customers

**Endpoint:** `GET http://localhost:9000/all-customers`

**Description:** Retrieves all customers. Requires ADMIN token for access.

### Request Headers

- **Authorization:** `Bearer <token>`

### Query Parameters

- **recordPerPage**: (optional) Number of records per page (default is 10).
- **page**: (optional) Page number (default is 1).

### Response

**Content-Type:** `application/json`

**Response Body:**

```json
{
  "total_count": 100,
  "customer_items": [
    {
      "first_name": "John",
      "last_name": "Doe",
      "email": "john.doe@example.com",
      "phone": "+1234567890",
      "company": "Acme Corporation",
      "status": "active",
      "notes": "Important customer",
      "customer_id": "60f7e3a4b9f1b2c6d8e4f4b1",
      "last_interaction": "2024-08-25T10:00:00Z",
      "company_id": "60f7e3a4b9f1b2c6d8e4f4b0"
    },
    ...
  ]
}
```

## Get Customers by Company

**Endpoint:** `GET http://localhost:9000/company/:company_id/customers`

**Description:** Retrieves all customers for a specific company. Requires authentication and authorization.

### Request Headers

- **Authorization:** `Bearer <token>`

### URL Parameters

- **company_id:** (required) ID of the company to retrieve customers from.

### Response

**Content-Type:** `application/json`

**Response Body:**

```json
[
  {
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "phone": "+1234567890",
    "company": "Acme Corporation",
    "status": "active",
    "notes": "Important customer",
    "customer_id": "60f7e3a4b9f1b2c6d8e4f4b1",
    "last_interaction": "2024-08-25T10:00:00Z",
    "company_id": "60f7e3a4b9f1b2c6d8e4f4b0"
  },
  ...
]
```

## Get Customer by Company ID and Customer ID

**Endpoint:** `GET http://localhost:9000/company/:company_id/customers/:customer_id`

**Description:** Retrieves details of a specific customer by their ID within a company. Requires authentication and authorization.

### Request Headers

- **Authorization:** `Bearer <token>`

### URL Parameters

- **company_id:** (required) ID of the company.
- **customer_id:** (required) ID of the customer to retrieve.

### Response

**Content-Type:** `application/json`

**Response Body:**

```json
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "phone": "+1234567890",
  "company": "Acme Corporation",
  "status": "active",
  "notes": "Important customer",
  "customer_id": "60f7e3a4b9f1b2c6d8e4f4b1",
  "last_interaction": "2024-08-25T10:00:00Z",
  "company_id": "60f7e3a4b9f1b2c6d8e4f4b0"
}
```
## Update Customer by Company ID and Customer ID

**Endpoint:** `PUT http://localhost:9000/company/:company_id/customers/:customer_id`

**Description:** Updates details of a specific customer within a company. Requires authentication and authorization.

### Request Headers

- **Authorization:** `Bearer <token>`

### URL Parameters

- **company_id:** (required) ID of the company.
- **customer_id:** (required) ID of the customer to update.

### Request Body

**Content-Type:** `application/json`

**Request Body:**

```json
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "phone": "+0987654321",
  "status": "inactive",
  "notes": "Updated notes",
  "company_id": "60f7e3a4b9f1b2c6d8e4f4b0"
}
```

## Delete Customer by Company ID and Customer ID

**Endpoint:** `DELETE http://localhost:9000/company/:company_id/customers/:customer_id`

**Description:** Deletes a specific customer by their ID within a company. Requires authentication and authorization.

### Request Headers

- **Authorization:** `Bearer <token>`

### URL Parameters

- **company_id:** (required) ID of the company.
- **customer_id:** (required) ID of the customer to delete.

### Response

**Content-Type:** `application/json`

**Response Body:**

```json
{
  "message": "Customer deleted successfully"
}
```

## Get Customer by User ID

**Endpoint:** `GET http://localhost:9000/customer/:user_id`

**Description:** Retrieves details of a specific customer by their user ID. Requires authentication.

### Request Headers

- **Authorization:** `Bearer <token>`

### URL Parameters

- **user_id:** (required) ID of the user to retrieve.

### Response

**Content-Type:** `application/json`

**Response Body:**

```json
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "phone": "+1234567890",
  "company": "Acme Corporation",
  "status": "active",
  "notes": "Important customer",
  "customer_id": "60f7e3a4b9f1b2c6d8e4f4b1",
  "last_interaction": "2024-08-25T10:00:00Z",
  "company_id": "60f7e3a4b9f1b2c6d8e4f4b0"
}

```
### Raise Ticket

**Endpoint:** `POST /interactions/:company_id/ticket`

**Description:** Creates a new ticket for a customer. Requires authentication.

**Request Headers:**

- **Authorization:** `Bearer <token>`

**URL Parameters:**

- **company_id:** (required) ID of the company where the ticket is raised.

**Request Body:**

```json
{
  "customer_id": "60f7e3a4b9f1b2c6d8e4f4b1",
  "type": "TICKET",
  "status": "OPEN"
}

```

## Create Meeting

**Endpoint:** `POST /interactions/:company_id/meeting`

**Description:** Creates a new meeting interaction. Requires authentication.

**Request Headers:**

- **Authorization:** `Bearer <token>`

**URL Parameters:**

- **company_id:** (required) ID of the company where the meeting is scheduled.

**Request Body:**

```json
{
  "customer_id": "60f7e3a4b9f1b2c6d8e4f4b1",
  "type": "MEETING",
  "status": "SCHEDULED",
  "scheduled_at": "2024-08-25T15:00:00Z"
}
```
## Update Interaction Status

**Endpoint:** `PUT /interactions/:interaction_id/status`

**Description:** Updates the status of a specific interaction. Requires authentication.

**Request Headers:**

- **Authorization:** `Bearer <token>`

**URL Parameters:**

- **interaction_id:** (required) ID of the interaction to update.

**Request Body:**

```json
{
  "status": "RESOLVED"
}

```

## Get Customer Interactions

**Endpoint:** `GET /customers/:customer_id/interactions`

**Description:** Retrieves all interactions for a specific customer. Requires authentication.

**Request Headers:**

- **Authorization:** `Bearer <token>`

**URL Parameters:**

- **customer_id:** (required) ID of the customer whose interactions are to be retrieved.

**Response:**

```json
[
  {
    "type": "MEETING",
    "status": "SCHEDULED",
    "scheduled_at": "2024-08-25T15:00:00Z",
    "interaction_id": "60f7e3a4b9f1b2c6d8e4f4b3"
  },
  ...
]

```

## Get Interaction Report

**Endpoint:** `GET /reports/interactions`

**Description:** Retrieves a report of interactions based on provided filters. Requires ADMIN token for access.

**Request Headers:**

- **Authorization:** `Bearer <token>`

**Query Parameters:**

- **start_date:** (optional) Start date for the report in ISO format (e.g., `2024-08-01T00:00:00Z`).
- **end_date:** (optional) End date for the report in ISO format (e.g., `2024-08-31T23:59:59Z`).
- **type:** (optional) Type of interactions to include in the report (e.g., `MEETING`, `CALL`).

**Response:**

```json
[
  {
    "type": "MEETING",
    "status": "SCHEDULED",
    "day": "2024-08-25",
    "count": 10
  },
  ...
]

```

## Get Conversion Rate Report

**Endpoint:** `GET /reports/conversion_rate`

**Description:** Retrieves a report on lead conversion rates. Requires ADMIN token for access.

**Request Headers:**

- **Authorization:** `Bearer <token>`

**Query Parameters:**

- **start_date:** (optional) Start date for the report in ISO format (e.g., `2024-08-01T00:00:00Z`).
- **end_date:** (optional) End date for the report in ISO format (e.g., `2024-08-31T23:59:59Z`).

**Response:**

```json
{
  "total_leads": 150,
  "total_customers": 50,
  "conversion_rate": 33.33
}

```

## Create Lead

**Endpoint:** `POST /leads`

**Description:** Creates a new lead. Requires ADMIN token for access.

**Request Headers:**

- **Authorization:** `Bearer <token>`

**Request Body:**

```json
{
  "name": "New Lead",
  "email": "lead@example.com",
  "phone": "+0987654321",
  "notes": "Lead details here"
}

```

## Send Email

**Endpoint:** `POST /email`

**Description:** Sends an email to specified recipients. Requires authentication.

**Request Headers:**

- **Authorization:** `Bearer <token>`

**Request Body:**

```json
{
  "to_addr": "findcant13@gmail.com",
  "subject": "Meeting Reminder",
  "body": "This is a reminder for your upcoming meeting scheduled for tomorrow."
}


```

## Conclusion

The MatriceCRM application is a robust and scalable CRM solution designed to streamline customer relationship management, enhance interaction tracking, and optimize lead management. With its backend built using Go and MongoDB, and integrated with email functionalities, MatriceCRM offers a comprehensive suite of features tailored to meet diverse business needs.

Key Features:
Customer Interaction Tracking: Efficiently manage and track interactions with customers, including meetings, tickets, and other communication activities.
Lead Management: Create and manage leads with ease, ensuring your sales pipeline remains organized and up-to-date.
Reporting and Analytics: Generate detailed reports on interactions and conversion rates to gain valuable insights and make data-driven decisions.
Email Integration: Seamlessly send and manage emails, enhancing communication efficiency.
Role-Based Access Control: Implement security and access controls to ensure data integrity and restrict access based on user roles.
Deployment and Setup:
MatriceCRM is designed for easy deployment using Docker, ensuring a consistent environment across development, testing, and production stages. With straightforward Docker build and run commands, you can quickly set up and start using the application.

By leveraging these features and functionalities, MatriceCRM empowers businesses to effectively manage customer relationships, improve sales processes, and drive overall growth. Whether you're looking to optimize internal workflows or enhance customer engagement, MatriceCRM provides a reliable and efficient solution for your CRM needs.

Thank you for choosing MatriceCRM. We hope it brings great value to your business operations