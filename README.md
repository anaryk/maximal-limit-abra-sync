# Maximal Limit Abra Sync

## Project Description
This project is a Go application that synchronizes data between an internal system and the Abra economic software. It implements API integration with Abra, allowing management of invoices, contacts, and other entities. The application is containerized and supports automation through cron jobs.

## Key Features
- **Automatic data synchronization** between the internal system and Abra API.
- **Management of contacts and invoices** via REST API.
- **Scheduled cron jobs** for periodic execution of tasks.
- **Docker support** for easy deployment.
- **CI/CD pipeline** using GitHub Actions for automated builds and deployments.
- **Database integration** for persistent storage of synchronized data.

---

## Project Structure

```
maximal-limit-abra-sync/
├── .github/workflows/       # CI/CD pipeline configuration
│   ├── docker-image.yml     # Workflow for building Docker image
│
├── docker/                  # Docker-related files
│   ├── Dockerfile           # Container definition
│
├── pkg/abra/                # Module for interacting with Abra API
│   ├── connector.go         # API connection management
│   ├── contacts.go          # Contact operations
│   ├── invoice.go           # Invoice operations
│   ├── model.go             # Data models
│
├── pkg/cron/                # Scheduled tasks (cron jobs)
│   ├── emailSendCron.go     # Automated email sending
│   ├── invoiceSyncer.go     # Synchronization of invoices
│
├── pkg/db/                  # Database integration
│   ├── connector.go         # Database connection handling
│   ├── internal.go          # Internal database utilities
│   ├── maxadmin_adapter.go  # Adapter for MaxAdmin
│   ├── model.go             # Database models
│
├── pkg/email/               # Email handling
│   ├── invoiceSender.go     # Sending invoices via email
│   ├── model.go             # Email-related models
│
├── pkg/internal/            # Internal utilities
│   ├── constants.go         # Application constants
│
├── pkg/utils/               # Utility functions
│   ├── utils.go             # General helper functions
│
├── main.go                  # Main application entry point
├── go.mod                   # Go module dependencies
├── go.sum                   # Dependency integrity check
```

---

## Installation and Running

### Requirements
- **Go 1.22+**
- **Docker (optional for containerized deployment)**
- **PostgreSQL or MySQL database**

### Running Locally

1. Clone the repository:
   ```sh
   git clone https://github.com/your-repo/maximal-limit-abra-sync.git
   cd maximal-limit-abra-sync
   ```
2. Install dependencies:
   ```sh
   go mod tidy
   ```
3. Set up the database and apply migrations (if applicable).
4. Start the application:
   ```sh
   go run main.go
   ```

### Running in Docker
1. Build the Docker image:
   ```sh
   docker build -t abra-sync .
   ```
2. Run the container:
   ```sh
   docker run -d --name abra-sync abra-sync
   ```

---

## Configuration

The application uses a `.env` file for storing necessary API keys and settings:

```
ABRA_API_URL=https://api.abra.cz
ABRA_API_KEY=your_api_key
DB_MAXADMIN_HOST=your_maxadmin_host
DB_MAXADMIN_NAME=your_maxadmin_db_name
DB_INTERNAL_HOST=your_internaldb_host
DB_INTERNAL_NAME=your_internal_db_name
ABRA_USER=your_abra_user
ENABLE_EMAIL_CRON=true
POSTAL_URL=https://your-postal-service.com
EMAIL_SMTP_SERVER=smtp.example.com
EMAIL_USERNAME=user@example.com
EMAIL_PASSWORD=your_password
```

---

## Features

### 1. Contact Synchronization
- **File:** `pkg/abra/contacts.go`
- **Description:** Fetches contact data from Abra API and stores it in the database.
- **Usage:**
  ```go
  contacts := abra.FetchContacts()
  fmt.Println(contacts)
  ```

### 2. Invoice Management
- **File:** `pkg/abra/invoice.go`
- **Description:** Creates and updates invoices in Abra.
- **Usage:**
  ```go
  invoice := abra.CreateInvoice(data)
  fmt.Println(invoice)
  ```

### 3. Scheduled Invoice Synchronization
- **File:** `pkg/cron/invoiceSyncer.go`
- **Description:** Periodically syncs invoices between the internal database and Abra.
- **Usage:** Runs automatically at scheduled intervals.

### 4. Email Invoice Sending
- **File:** `pkg/email/invoiceSender.go`
- **Description:** Sends invoices to customers via email.
- **Usage:**
  ```go
  email.SendInvoice(invoiceID)
  ```

---

## Database Integration
The application supports PostgreSQL and MySQL as a storage backend.

### Database Models
- **File:** `pkg/db/model.go`
- **Entities:** Contacts, Invoices, SyncLogs
- **Example Contact Model:**
  ```go
  type Contact struct {
      ID        int    `json:"id"`
      Name      string `json:"name"`
      Email     string `json:"email"`
      Phone     string `json:"phone"`
  }
  ```

### Database Connection
- **File:** `pkg/db/connector.go`
- **Handles:**
  - Connection pooling
  - Query execution
  - Error handling

---

## CI/CD Pipeline
This project uses **GitHub Actions** for automated building and deployment.

- **`docker-image.yml`**: Automatically builds a Docker image on push.

---

## Contribution
If you want to contribute to the project:
1. Fork this repository.
2. Create a new branch (`feature/new-feature`).
3. Make changes and commit them.
4. Submit a **Pull Request**.

---

## License
This project is available under the **MIT** license.

