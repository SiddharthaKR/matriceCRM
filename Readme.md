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

makefile
Copy code
MONGO_URI=<your_mongo_uri>
JWT_SECRET=<your_jwt_secret>
Run the Application:

bash
Copy code
go run main.go
Docker Setup (Optional): Build and run the application in a Docker container:

bash
Copy code
docker build -t matricecrm .
docker run -p 8080:8080 matricecrm
Database Schema Design
Customer Collection
json
Copy code
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
The backend is designed as a RESTful API with the following components:

Gin Framework: For handling HTTP requests and routing.
MongoDB: NoSQL database for storing customer and company data.
JWT Authentication: Secure token-based authentication for API access.
Docker: Containerization for consistent development and deployment environments.
High-Level Architecture
arduino
Copy code
Client -> Gin API Server -> MongoDB
API Documentation
Authentication
Login

Endpoint: /login
Method: POST
Request Body:
json
Copy code
{
  "email": "user@example.com",
  "password": "password123"
}
Response:
json
Copy code
{
  "token": "jwt_token",
  "refresh_token": "refresh_token"
}
Signup

Endpoint: /signup
Method: POST
Request Body:
json
Copy code
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "phone": "1234567890",
  "company_id": "ObjectId",
  "password_hash": "password123"
}
Response:
json
Copy code
{
  "message": "Signup successful"
}
Customers
Get Customers

Endpoint: /customers
Method: GET
Response:
json
Copy code
{
  "total_count": 100,
  "customer_items": [...]
}
Get Customer by ID

Endpoint: /customers/:id
Method: GET
Response:
json
Copy code
{
  "customer": {...}
}
Companies
Create Company
Endpoint: /companies
Method: POST
Request Body:
json
Copy code
{
  "name": "Company Name",
  "address": "Company Address",
  "phone": "1234567890"
}
Response:
json
Copy code
{
  "message": "Company created successfully"
}