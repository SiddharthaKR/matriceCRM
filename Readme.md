## Setup and Installation

### Prerequisites

- [Go](https://golang.org/dl/) (version 1.18+)
- [MongoDB](https://www.mongodb.com/try/download/community) (version 4.4+)
- [Docker](https://docs.docker.com/get-docker/) (optional for containerized setup)

### Installation Steps

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/SiddharthaKR/matriceCRM.git
   cd MatriceCRM
Set Up Environment Variables: Create a .env file in the root directory and populate it with the following variables:


# CRM Backend System

The CRM Backend System is designed to manage customers and users efficiently, with a focus on robust user management and a scalable system architecture.

## User Management

The system provides comprehensive CRUD operations and detailed management of users and customers.

### Features

- **CRUD Operations**: 
  The system supports full Create, Read, Update, and Delete functionalities for both customers and users.

- **User/Customer Details**: 
  Each user or customer is associated with the following details:
  - Name
  - Contact Information
  - Company (if applicable)
  - Status
  - Notes

## System Design

The backend architecture is designed as a RESTful API, ensuring scalability and flexibility.

### Components

- **Gin Framework**: 
  Used for handling HTTP requests and routing.
  
- **MongoDB**: 
  A NoSQL database for storing customer and company data.

- **JWT Authentication**: 
  Provides secure token-based authentication for API access.

- **Docker**: 
  Facilitates containerization for consistent development and deployment environments.

# Interaction Tracking

The CRM system includes robust interaction tracking features to manage customer relationships effectively.

## Features

- **Ticket Management**:
  - Customers can raise a ticket for issues or inquiries.
  - Tickets can be marked as resolved by the customer.

- **Meeting Scheduling**:
  - Users can schedule interactions with customers, such as meetings.

- **Interaction Linking**:
  - Ensure that all interactions are linked to the appropriate customer for accurate tracking.

- **Interaction History**:
  - Users can view the complete history of interactions for a specific customer, allowing for better follow-up and service.
# Analytics and Reporting

The CRM system provides advanced analytics and reporting features to help users gain insights into customer and lead interactions.

## Features

- **Report Generation**:
  - Endpoints are provided to generate reports on customer/lead interactions, conversion rates, and other key metrics.

- **Visual Representations**:
  - Reports can include visual elements such as charts and graphs to better illustrate data.

# Email Integration
 The CRM system includes built-in email integration to streamline communication and track engagement.

## Features

- **Email Sending**:
  - Users can send emails directly from the CRM interface.

- **Email Tracking**:
  - The system tracks email open rates and responses, providing valuable insights into customer engagement.



# CRM Application Backend

## Table of Contents
- [Introduction](#introduction)
- [Database Schema Design](#database-schema-design)
- [System Design](#system-design)
- [Setup and Running Locally](#setup-and-running-locally)
- [API Documentation](#api-documentation)

## Introduction
This backend application serves as a CRM (Customer Relationship Management) system for managing users, customers, interactions, and more. The application is built using Go with the Gin framework and uses MongoDB for data storage.

## Database Schema Design

### Customer Collection
```json
{
  "_id": "ObjectId",
  "first_name": "string",
  "last_name": "string",
  "email": "string",
  "phone": "string",
  "company": "string",
  "status": "string",
  "notes": "string",
  "created_at": "DateTime",
  "updated_at": "DateTime",
  "customer_id": "string",
  "last_interaction": "DateTime",
  "company_id": "ObjectId",
  "password_hash": "string",
  "token": "string",
  "refresh_token": "string"
}
Company Collection
json
Copy code
{
  "_id": "ObjectId",
  "name": "string",
  "address": "string",
  "phone": "string",
  "created_at": "DateTime",
  "updated_at": "DateTime",
  "company_id": "string"
}
System Design
Architecture
The backend is designed as a RESTful API with the following components:

Gin Framework: For handling HTTP requests and routing.
MongoDB: NoSQL database for storing customer and company data.
JWT Authentication: Secure token-based authentication for API access.
Docker: Containerization for consistent development and deployment environments.
High-Level Architecture
arduino
Copy code
Client -> Gin API Server -> MongoDB
Setup and Running Locally
Prerequisites
Go (version 1.18+)
MongoDB (version 4.4+)
Docker (optional for containerized setup)
Steps
Clone the Repository

bash
Copy code
git clone https://github.com/your-repo/crm-backend.git
cd crm-backend
Set Up Environment Variables Create a .env file in the root directory with the following content:

plaintext
Copy code
MONGO_URI=<your_mongo_uri>
JWT_SECRET=<your_jwt_secret>
Run the Application

bash
Copy code
go run main.go
Docker Setup (Optional) Build and run the application in a Docker container:

bash
Copy code
docker build -t crm-backend .
docker run -p 8080:8080 crm-backend
API Documentation
Interaction Tracking
Raise a Ticket
URL: /api/tickets
Method: POST
Request Body:
json
Copy code
{
  "customer_id": "customer_id",
  "issue": "Issue description"
}
Response:
json
Copy code
{
  "message": "Ticket raised successfully"
}
Mark Ticket as Resolved
URL: /api/tickets/:ticketId/resolve
Method: PATCH
Response:
json
Copy code
{
  "message": "Ticket marked as resolved"
}
Schedule Interaction
URL: /api/interactions
Method: POST
Request Body:
json
Copy code
{
  "customer_id": "customer_id",
  "date_time": "2024-08-01T14:00:00Z",
  "description": "Meeting description"
}
Response:
json
Copy code
{
  "message": "Interaction scheduled"
}
View Interaction History
URL: /api/customers/:customerId/interactions
Method: GET
Response:
json
Copy code
{
  "interactions": [
    {
      "date_time": "2024-08-01T14:00:00Z",
      "description": "Meeting description"
    }
  ]
}
Analytics and Reporting (Optional)
Generate Reports
URL: /api/reports
Method: GET
Response:
json
Copy code
{
  "total_interactions": 100,
  "conversion_rate": "50%",
  "metrics": {
    "chart": "URL to chart image"
  }
}
Email Integration (Optional)
Send Email
URL: /api/emails
Method: POST
Request Body:
json
Copy code
{
  "to": "recipient@example.com",
  "subject": "Subject",
  "body": "Email body"
}
Response:
json
Copy code
{
  "message": "Email sent successfully"
}
Track Email Engagement
URL: /api/emails/:emailId/engagement
Method: GET
Response:
json
Copy code
{
  "open_rate": "70%",
  "responses": [
    {
      "response_time": "2024-08-01T15:00:00Z",
      "response": "Thank you for your email"
    }
  ]
}
Conclusion
This documentation provides an overview of the backend setup, including the database schema design, system architecture, steps to run the project locally, and detailed API documentation. Follow the provided instructions to set up and run the project, and refer to the API documentation for interacting with the backend endpoints.