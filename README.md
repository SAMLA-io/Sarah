# Sarah Campaign Management API

A comprehensive backend service for managing automated calling campaigns using VapiAI for voice interactions and MongoDB for data persistence.

## Overview

The Sarah Campaign Management API is designed to handle automated calling campaigns with flexible scheduling, customer management, and integration with VapiAI's voice AI platform. The system supports various campaign types including weekly, monthly, and yearly recurrences, with sophisticated scheduling based on customer-specific dates like insurance renewals.

Given the nature of the project, we are using a MongoDB database to store the data. Additionally, VapiAI's native campaigns do not support the ability to schedule calls for a specific date. Therefore, we are using our own campaign management system to schedule calls. If the campaign is a one-time campaign, we will still use our own campaign management system to schedule the call, as to maintain a single campaign management system.

## Features

- **Campaign Management**: Create and manage automated calling campaigns
- **Flexible Scheduling**: Support for weekly, monthly, yearly, and one-time campaigns
- **Customer Management**: Store and manage customer contact information
- **VapiAI Integration**: Seamless integration with VapiAI for voice interactions
- **Organization-based Architecture**: Multi-tenant design with organization isolation
- **MongoDB Persistence**: Scalable data storage with MongoDB

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP Client   │───▶│   Sarah API     │───▶│   MongoDB       │
│                 │    │   (Port 8080)   │    │   Database      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │   VapiAI API    │
                       │   (Voice AI)    │
                       └─────────────────┘
```

## Prerequisites

- Go 1.24.4 or higher
- MongoDB instance
- VapiAI account and API key

## Installation

1. Clone the repository:
```bash
git clone https://github.com/SAMLA-io/Sarah.git
cd Sarah
```

2. Install dependencies:
```bash
go mod download
```

3. Create a `.env` file with the following variables:
```env
# MongoDB Configuration
MONGO_URI=mongodb://localhost:27017
MONGO_COLLECTION_CAMPAIGNS=campaigns
MONGO_COLLECTION_ASSISTANTS=assistants
MONGO_COLLECTION_CONTACTS=contacts
MONGO_COLLECTION_PHONE_NUMBERS=phone_numbers

# VapiAI Configuration
VAPI_API_KEY=your_vapi_api_key_here
```

4. Run the application:
```bash
go run main.go
```

The API will be available at `http://localhost:8080`

## API Endpoints

### Health Check

#### GET /test
Simple health check endpoint.

**Response:**
```json
"Hello, World!"
```

### Campaign Management

#### POST /campaigns/create
Create a new campaign.

**Query Parameters:**
- `orgId` (required): Organization ID

**Request Body:**
```json
{
  "campaignCreateRequest": {
    "name": "Weekly Insurance Reminders",
    "assistant_id": "asst_1234567890abcdef",
    "phone_number_id": "phone_0987654321fedcba",
    "schedule_plan": {
      "before_day": 3,
      "after_day": 0,
      "week_days": [1, 3, 5],
      "month_days": [],
      "year_months": []
    },
    "customers": [
      {
        "phone_number": "+1234567890",
        "day_number": 15,
        "month_number": 3,
        "week_day": 1,
        "custom_date": null,
        "expiry_date": "2024-12-31T23:59:59Z"
      }
    ],
    "type": "recurrent_weekly",
    "status": "active",
    "start_date": "2024-01-01T00:00:00Z",
    "end_date": "2024-12-31T23:59:59Z",
    "timezone": "America/New_York"
  }
}
```

**Response:**
```json
{
  "name": "Weekly Insurance Reminders",
  "assistant_id": "asst_1234567890abcdef",
  "phone_number_id": "phone_0987654321fedcba",
  "schedule_plan": { ... },
  "customers": [ ... ],
  "type": "recurrent_weekly",
  "status": "active",
  "start_date": "2024-01-01T00:00:00Z",
  "end_date": "2024-12-31T23:59:59Z",
  "timezone": "America/New_York"
}
```

#### GET /campaigns/org
Retrieve all campaigns for an organization.

**Query Parameters:**
- `orgId` (required): Organization ID

**Response:**
```json
[
  {
    "name": "Weekly Insurance Reminders",
    "assistant_id": "asst_1234567890abcdef",
    "phone_number_id": "phone_0987654321fedcba",
    "schedule_plan": { ... },
    "customers": [ ... ],
    "type": "recurrent_weekly",
    "status": "active",
    "start_date": "2024-01-01T00:00:00Z",
    "end_date": "2024-12-31T23:59:59Z",
    "timezone": "America/New_York"
  }
]
```

### Call Management

#### POST /calls/create
Create a new call using VapiAI.

**Query Parameters:**
- `assistantId` (required): VapiAI assistant ID
- `assistantNumberId` (required): VapiAI phone number ID

**Request Body:**
```json
{
  "phoneNumbers": ["+1234567890", "+1987654321"]
}
```

**Response:**
```json
{
  "id": "call_abc123def456",
  "assistantId": "asst_1234567890abcdef",
  "phoneNumberId": "phone_0987654321fedcba",
  "status": "queued",
  "createdAt": "2024-01-01T12:00:00Z"
}
```

#### GET /calls/call
Retrieve a specific call by ID.

**Query Parameters:**
- `callId` (required): VapiAI call ID

**Response:**
```json
{
  "id": "call_abc123def456",
  "assistantId": "asst_1234567890abcdef",
  "phoneNumberId": "phone_0987654321fedcba",
  "status": "completed",
  "duration": 120,
  "createdAt": "2024-01-01T12:00:00Z",
  "endedAt": "2024-01-01T12:02:00Z"
}
```

#### GET /calls/list
List calls based on criteria.

**Request Body (optional):**
```json
{
  "callListRequest": {
    "assistantId": "asst_1234567890abcdef",
    "limit": 10,
    "offset": 0,
    "status": "completed"
  }
}
```

#### GET /calls/org
Retrieve all calls for an organization.

**Query Parameters:**
- `orgId` (required): Organization ID

### Organization Resources

#### GET /assistants/org
Retrieve all assistants for an organization.

**Query Parameters:**
- `orgId` (required): Organization ID

**Response:**
```json
[
  {
    "id": "asst_1234567890abcdef",
    "name": "Insurance Reminder Assistant",
    "vapiAssistantId": "asst_1234567890abcdef",
    "organizationId": "org_1234567890abcdef"
  }
]
```

#### GET /contacts/org
Retrieve all contacts for an organization.

**Query Parameters:**
- `orgId` (required): Organization ID

**Response:**
```json
[
  {
    "id": "contact_1234567890abcdef",
    "name": "John Doe",
    "phoneNumber": "+1234567890",
    "email": "john.doe@example.com",
    "organizationId": "org_1234567890abcdef"
  }
]
```

#### GET /phone_numbers/org
Retrieve all phone numbers for an organization.

**Query Parameters:**
- `orgId` (required): Organization ID

**Response:**
```json
[
  {
    "id": "phone_0987654321fedcba",
    "phoneNumber": "+1987654321",
    "vapiPhoneNumberId": "phone_0987654321fedcba",
    "organizationId": "org_1234567890abcdef"
  }
]
```

## Data Models

### Campaign
```go
type Campaign struct {
    Name          string         // Human-readable campaign name
    AssistantId   string         // VapiAI assistant ID
    PhoneNumberId string         // VapiAI phone number ID
    SchedulePlan  *SchedulePlan  // Scheduling configuration
    Customers     []Customer     // List of customers to contact
    Type          CampaignType   // Campaign recurrence type
    Status        CampaignStatus // Current campaign status
    StartDate     *time.Time     // Campaign start date
    EndDate       *time.Time     // Campaign end date
    TimeZone      string         // Timezone for date calculations
}
```

### SchedulePlan
```go
type SchedulePlan struct {
    BeforeDay  int   // Days before customer's relevant date
    AfterDay   int   // Days after customer's relevant date
    WeekDays   []int // Allowed days of week (0=Sunday, 1=Monday, etc.)
    MonthDays  []int // Allowed days of month (1-31)
    YearMonths []int // Allowed months (1-12)
}
```

### Customer
```go
type Customer struct {
    PhoneNumber string     // Customer's phone number (E.164 format)
    DayNumber   int        // Day of month for scheduling
    MonthNumber int        // Month for scheduling (1-12)
    WeekDay     int        // Day of week for scheduling (0-6)
    CustomDate  *time.Time // One-time specific date
    ExpiryDate  *time.Time // Insurance/subscription expiry date
}
```

## Campaign Types

- `recurrent_weekly`: Runs on specified days of the week
- `recurrent_monthly`: Runs on specified days of the month
- `recurrent_yearly`: Runs on specified dates annually
- `one_time`: Runs only once on a specific date

## Campaign Statuses

- `active`: Campaign is currently running
- `paused`: Campaign is temporarily stopped
- `completed`: Campaign has finished all scheduled calls
- `cancelled`: Campaign has been permanently stopped

## Error Handling

The API returns appropriate HTTP status codes:

- `200 OK`: Successful operation
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request data
- `405 Method Not Allowed`: Incorrect HTTP method
- `500 Internal Server Error`: Server-side error

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `MONGO_URI` | MongoDB connection string | Yes |
| `MONGO_COLLECTION_CAMPAIGNS` | Campaigns collection name | Yes |
| `MONGO_COLLECTION_ASSISTANTS` | Assistants collection name | Yes |
| `MONGO_COLLECTION_CONTACTS` | Contacts collection name | Yes |
| `MONGO_COLLECTION_PHONE_NUMBERS` | Phone numbers collection name | Yes |
| `VAPI_API_KEY` | VapiAI API key | Yes |

## Development

### Project Structure
```
Sarah/
├── api/                    # HTTP handlers and API endpoints
│   ├── campaigns.go        # Campaign management endpoints
│   ├── calls.go           # Call management endpoints
│   ├── assistants.go      # Assistant management endpoints
│   ├── contacts.go        # Contact management endpoints
│   ├── phone_numbers.go   # Phone number management endpoints
│   └── utils.go           # Shared utility functions
├── mongodb/               # Database operations
│   ├── campaigns.go       # Campaign database operations
│   ├── assistants.go      # Assistant database operations
│   ├── contacts.go        # Contact database operations
│   └── phone_numbers.go   # Phone number database operations
├── types/                 # Data type definitions
│   └── mongodb/           # MongoDB-specific types
│       ├── campaigns.go   # Campaign data structures
│       ├── assistants.go  # Assistant data structures
│       ├── contact.go     # Contact data structures
│       └── phone_numbers.go # Phone number data structures
├── main.go               # Application entry point
├── go.mod               # Go module file
├── go.sum               # Go module checksums
├── README.md            # This file
└── sample_campaigns.json # Example campaign data
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## Support

For support and questions, please contact the development team or create an issue in the repository.
