# Sarah AI Call assistant

A comprehensive backend service for managing automated calling assistant using VapiAI for voice interactions and MongoDB for data persistence, with Clerk authentication and organization-based multi-tenancy.

## Overview

The Sarah AI Call assistant is designed to handle automated calling assistant with flexible scheduling, customer management, and integration with VapiAI's voice AI platform. The system supports various campaign types including weekly, monthly, and yearly recurrences, with sophisticated scheduling based on customer-specific dates like insurance renewals. Additionally, you can create single calls and other behaviors from the VapiAI API.

Given the nature of the project, we are using a MongoDB database to store the data. Additionally, VapiAI's native campaigns do not support the ability to schedule calls for a specific date. Therefore, we are using our own campaign management system to schedule calls. If the campaign is a one-time campaign, we will still use our own campaign management system to schedule the call, as to maintain a single campaign management system.

## Features

- **Campaign Management**: Create and manage automated calling campaigns
- **Flexible Scheduling**: Support for weekly, monthly, yearly, and one-time campaigns
- **Customer Management**: Store and manage customer contact information
- **VapiAI Integration**: Seamless integration with VapiAI for voice interactions
- **Organization-based Architecture**: Multi-tenant design with Clerk authentication and organization isolation
- **MongoDB Persistence**: Scalable data storage with MongoDB
- **Authentication & Authorization**: Secure access control using Clerk

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
                              │
                              ▼
                       ┌─────────────────┐
                       │   Clerk Auth    │
                       │   (Identity)    │
                       └─────────────────┘
```

## Prerequisites

- Go 1.24.4 or higher
- MongoDB instance
- VapiAI account and API key
- Clerk account and API key

## Installation

### Option 1: Docker (Recommended)

1. Clone the repository:
```bash
git clone https://github.com/SAMLA-io/Sarah.git
cd Sarah
```

2. Create a `.env` file with the following variables:
```env
# MongoDB Configuration
MONGO_URI=mongodb://localhost:27017
MONGO_COLLECTION_CAMPAIGNS=campaigns
MONGO_COLLECTION_ASSISTANTS=assistants
MONGO_COLLECTION_CONTACTS=contacts
MONGO_COLLECTION_PHONE_NUMBERS=phone_numbers

# VapiAI Configuration
VAPI_API_KEY=your_vapi_api_key_here

# Clerk Configuration
CLERK_SECRET_KEY=your_clerk_secret_key_here
```


3. Build and run the Docker container:

```bash
docker build -t sarah .
docker run -p 8080:8080 --env-file .env sarah
```

The API will be available at `http://localhost:8080`

### Option 2: Local Development

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

# Clerk Configuration
CLERK_SECRET_KEY=your_clerk_secret_key_here
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

**Headers:**
- `Authorization: Bearer <clerk_jwt_token>` (required)

**Request Body:**
```json
{
  "campaignCreateRequest": {
    "name": "Weekly Insurance Reminders",
    "assistant_id": "asst_1234567890abcdef",
    "phone_number_id": "phone_0987654321fedcba",
    "schedule_plan": {
      "before_day": 3,
      "after_day": 0
    },
    "customers": [
      {
        "phone_number": "+1234567890",
        "day_number": 15,
        "month_number": 3,
        "year_number": 2024
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

**Headers:**
- `Authorization: Bearer <clerk_jwt_token>` (required)

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

**Headers:**
- `Authorization: Bearer <clerk_jwt_token>` (required)

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

**Headers:**
- `Authorization: Bearer <clerk_jwt_token>` (required)

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

**Headers:**
- `Authorization: Bearer <clerk_jwt_token>` (required)

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

**Headers:**
- `Authorization: Bearer <clerk_jwt_token>` (required)

### Organization Resources

#### GET /assistants/org
Retrieve all assistants for an organization.

**Headers:**
- `Authorization: Bearer <clerk_jwt_token>` (required)

**Response:**
```json
[
  {
    "id": "507f1f77bcf86cd799439011",
    "name": "Insurance Reminder Assistant",
    "vapi_assistant_id": "asst_1234567890abcdef",
    "type": "insurance"
  }
]
```

#### POST /assistants/create
Create a new assistant.

**Headers:**
- `Authorization: Bearer <clerk_jwt_token>` (required)

**Request Body:**
```json
{
  "id": "foo",
  "orgId": "foo",
  "createdAt": "foo",
  "updatedAt": "foo",
  "transcriber": {
    "provider": "assembly-ai",
    "language": "en",
    "confidenceThreshold": 0.4,
    "enableUniversalStreamingApi": false,
    "formatTurns": false,
    "endOfTurnConfidenceThreshold": 0.7,
    "minEndOfTurnSilenceWhenConfident": 160,
    "wordFinalizationMaxWaitTime": 160,
    "maxTurnSilence": 400,
    "realtimeUrl": "foo",
    "wordBoost": [
      "foo"
    ],
    "endUtteranceSilenceThreshold": 42,
    "disablePartialTranscripts": true,
    "fallbackPlan": {
      "transcribers": [
        {
          "provider": "assembly-ai",
          "language": "en",
          "confidenceThreshold": 0.4,
          "enableUniversalStreamingApi": false,
          "formatTurns": false,
          "endOfTurnConfidenceThreshold": 0.7,
          "minEndOfTurnSilenceWhenConfident": 160,
          "wordFinalizationMaxWaitTime": 160,
          "maxTurnSilence": 400,
          "realtimeUrl": "foo",
          "wordBoost": [
            "foo"
          ],
          "endUtteranceSilenceThreshold": 42,
          "disablePartialTranscripts": true
        }
      ]
    }
  },
  "model": {
    "messages": [
      {
        "content": "foo",
        "role": "assistant"
      }
    ],
    "tools": [
      {
        "messages": [
          {
            "contents": [
              {
                "type": "text",
                "text": "foo",
                "language": "aa"
              }
            ],
            "type": "request-start",
            "blocking": false,
            "content": "foo",
            "conditions": [
              {
                "operator": "eq",
                "param": "foo",
                "value": "foo"
              }
            ]
          }
        ],
        "type": "apiRequest",
        "method": "POST",
        "timeoutSeconds": 20,
        "name": "foo",
        "description": "foo",
        "url": "foo",
        "body": {
          "type": "string",
          "items": {},
          "properties": {},
          "description": "foo",
          "pattern": "foo",
          "format": "date-time",
          "required": [
            "foo"
          ],
          "enum": [
            "foo"
          ],
          "title": "foo"
        },
        "headers": {
          "type": "string",
          "items": {},
          "properties": {},
          "description": "foo",
          "pattern": "foo",
          "format": "date-time",
          "required": [
            "foo"
          ],
          "enum": [
            "foo"
          ],
          "title": "foo"
        },
        "backoffPlan": {
          "type": "fixed",
          "maxRetries": 0,
          "baseDelaySeconds": 1
        },
        "variableExtractionPlan": {
          "schema": {
            "type": "string",
            "items": {},
            "properties": {},
            "description": "foo",
            "pattern": "foo",
            "format": "date-time",
            "required": [
              "foo"
            ],
            "enum": [
              "foo"
            ],
            "title": "foo"
          },
          "aliases": [
            {
              "key": "foo",
              "value": "foo"
            }
          ]
        }
      }
    ],
    "toolIds": [
      "foo"
    ],
    "knowledgeBase": {
      "provider": "custom-knowledge-base",
      "server": {
        "timeoutSeconds": 20,
        "url": "foo",
        "headers": {},
        "backoffPlan": {
          "type": "fixed",
          "maxRetries": 0,
          "baseDelaySeconds": 1
        }
      }
    },
    "knowledgeBaseId": "foo",
    "model": "claude-3-opus-20240229",
    "provider": "anthropic",
    "thinking": {
      "type": "enabled",
      "budgetTokens": 42
    },
    "temperature": 42,
    "maxTokens": 42,
    "emotionRecognitionEnabled": true,
    "numFastTurns": 42
  },
  "voice": {
    "cachingEnabled": true,
    "provider": "azure",
    "voiceId": "andrew",
    "chunkPlan": {
      "enabled": true,
      "minCharacters": 30,
      "punctuationBoundaries": "。",
      "formatPlan": {
        "enabled": true,
        "numberToDigitsCutoff": 2025,
        "replacements": [
          {
            "type": "exact",
            "replaceAllEnabled": false,
            "key": "foo",
            "value": "foo"
          }
        ],
        "formattersEnabled": "markdown"
      }
    },
    "speed": 42,
    "fallbackPlan": {
      "voices": [
        {
          "cachingEnabled": true,
          "provider": "azure",
          "voiceId": "andrew",
          "speed": 42,
          "chunkPlan": {
            "enabled": true,
            "minCharacters": 30,
            "punctuationBoundaries": "。",
            "formatPlan": {
              "enabled": true,
              "numberToDigitsCutoff": 2025,
              "replacements": [
                {
                  "type": {},
                  "replaceAllEnabled": {},
                  "key": {},
                  "value": {}
                }
              ],
              "formattersEnabled": "markdown"
            }
          },
          "oneOf": null
        }
      ]
    }
  },
  "firstMessage": "Hello! How can I help you today?",
  "firstMessageInterruptionsEnabled": false,
  "firstMessageMode": "assistant-speaks-first",
  "voicemailDetection": {
    "beepMaxAwaitSeconds": 30,
    "provider": "google",
    "backoffPlan": {
      "startAtSeconds": 5,
      "frequencySeconds": 5,
      "maxRetries": 6
    }
  },
  "clientMessages": "conversation-update",
  "serverMessages": "conversation-update",
  "maxDurationSeconds": 600,
  "backgroundSound": "off",
  "modelOutputInMessagesEnabled": false,
  "transportConfigurations": [
    {
      "provider": "twilio",
      "timeout": 60,
      "record": false,
      "recordingChannels": "mono"
    }
  ],
  "observabilityPlan": {
    "provider": "langfuse",
    "tags": [
      "foo"
    ],
    "metadata": {}
  },
  "credentials": [
    {
      "provider": "anthropic",
      "apiKey": "foo",
      "name": "foo"
    }
  ],
  "hooks": [
    {
      "on": "call.ending",
      "do": [
        {
          "type": "tool",
          "tool": {
            "messages": [
              {
                "contents": [
                  {
                    "type": "text",
                    "text": "foo",
                    "language": "aa"
                  }
                ],
                "type": "request-start",
                "blocking": false,
                "content": "foo",
                "conditions": [
                  {
                    "operator": "eq",
                    "param": "foo",
                    "value": "foo"
                  }
                ]
              }
            ],
            "type": "apiRequest",
            "method": "POST",
            "timeoutSeconds": 20,
            "name": "foo",
            "description": "foo",
            "url": "foo",
            "body": {
              "type": "string",
              "items": {},
              "properties": {},
              "description": "foo",
              "pattern": "foo",
              "format": "date-time",
              "required": [
                "foo"
              ],
              "enum": [
                "foo"
              ],
              "title": "foo"
            },
            "headers": {
              "type": "string",
              "items": {},
              "properties": {},
              "description": "foo",
              "pattern": "foo",
              "format": "date-time",
              "required": [
                "foo"
              ],
              "enum": [
                "foo"
              ],
              "title": "foo"
            },
            "backoffPlan": {
              "type": "fixed",
              "maxRetries": 0,
              "baseDelaySeconds": 1
            },
            "variableExtractionPlan": {
              "schema": {
                "type": "string",
                "items": {},
                "properties": {},
                "description": "foo",
                "pattern": "foo",
                "format": "date-time",
                "required": [
                  "foo"
                ],
                "enum": [
                  "foo"
                ],
                "title": "foo"
              },
              "aliases": [
                {
                  "key": "foo",
                  "value": "foo"
                }
              ]
            }
          },
          "toolId": "foo"
        }
      ],
      "filters": [
        {
          "type": "oneOf",
          "key": "foo",
          "oneOf": [
            "foo"
          ]
        }
      ]
    }
  ],
  "name": "foo",
  "voicemailMessage": "foo",
  "endCallMessage": "foo",
  "endCallPhrases": [
    "foo"
  ],
  "compliancePlan": {
    "hipaaEnabled": {
      "hipaaEnabled": false
    },
    "pciEnabled": {
      "pciEnabled": false
    }
  },
  "metadata": {},
  "backgroundSpeechDenoisingPlan": {
    "smartDenoisingPlan": {
      "enabled": false
    },
    "fourierDenoisingPlan": {
      "enabled": false,
      "mediaDetectionEnabled": true,
      "staticThreshold": -35,
      "baselineOffsetDb": -15,
      "windowSizeMs": 3000,
      "baselinePercentile": 85
    }
  },
  "analysisPlan": {
    "minMessagesThreshold": 42,
    "summaryPlan": {
      "messages": [
        {}
      ],
      "enabled": true,
      "timeoutSeconds": 42
    },
    "structuredDataPlan": {
      "messages": [
        {}
      ],
      "enabled": true,
      "schema": {
        "type": "string",
        "items": {},
        "properties": {},
        "description": "foo",
        "pattern": "foo",
        "format": "date-time",
        "required": [
          "foo"
        ],
        "enum": [
          "foo"
        ],
        "title": "foo"
      },
      "timeoutSeconds": 42
    },
    "structuredDataMultiPlan": [
      {
        "key": "foo",
        "plan": {
          "messages": [
            {}
          ],
          "enabled": true,
          "schema": {
            "type": "string",
            "items": {},
            "properties": {},
            "description": "foo",
            "pattern": "foo",
            "format": "date-time",
            "required": [
              "foo"
            ],
            "enum": [
              "foo"
            ],
            "title": "foo"
          },
          "timeoutSeconds": 42
        }
      }
    ],
    "successEvaluationPlan": {
      "rubric": "NumericScale",
      "messages": [
        {}
      ],
      "enabled": true,
      "timeoutSeconds": 42
    }
  },
  "artifactPlan": {
    "recordingEnabled": true,
    "recordingFormat": "wav;l16",
    "videoRecordingEnabled": false,
    "pcapEnabled": true,
    "pcapS3PathPrefix": "/pcaps",
    "transcriptPlan": {
      "enabled": true,
      "assistantName": "foo",
      "userName": "foo"
    },
    "recordingPath": "foo"
  },
  "messagePlan": {
    "idleMessages": [
      "foo"
    ],
    "idleMessageMaxSpokenCount": 42,
    "idleMessageResetCountOnUserSpeechEnabled": true,
    "idleTimeoutSeconds": 42,
    "silenceTimeoutMessage": "foo"
  },
  "startSpeakingPlan": {
    "waitSeconds": 0.4,
    "smartEndpointingPlan": {
      "provider": "vapi"
    },
    "customEndpointingRules": [
      {
        "type": "assistant",
        "regex": "foo",
        "regexOptions": [
          {
            "type": "ignore-case",
            "enabled": true
          }
        ],
        "timeoutSeconds": 42
      }
    ],
    "transcriptionEndpointingPlan": {
      "onPunctuationSeconds": 0.1,
      "onNoPunctuationSeconds": 1.5,
      "onNumberSeconds": 0.5
    },
    "smartEndpointingEnabled": false
  },
  "stopSpeakingPlan": {
    "numWords": 0,
    "voiceSeconds": 0.2,
    "backoffSeconds": 1,
    "acknowledgementPhrases": [
      "i understand",
      "i see",
      "i got it",
      "i hear you",
      "im listening",
      "im with you",
      "right",
      "okay",
      "ok",
      "sure",
      "alright",
      "got it",
      "understood",
      "yeah",
      "yes",
      "uh-huh",
      "mm-hmm",
      "gotcha",
      "mhmm",
      "ah",
      "yeah okay",
      "yeah sure"
    ],
    "interruptionPhrases": [
      "stop",
      "shut",
      "up",
      "enough",
      "quiet",
      "silence",
      "but",
      "dont",
      "not",
      "no",
      "hold",
      "wait",
      "cut",
      "pause",
      "nope",
      "nah",
      "nevermind",
      "never",
      "bad",
      "actually"
    ]
  },
  "monitorPlan": {
    "listenEnabled": false,
    "listenAuthenticationEnabled": false,
    "controlEnabled": false,
    "controlAuthenticationEnabled": false
  },
  "credentialIds": [
    "foo"
  ],
  "server": {
    "timeoutSeconds": 20,
    "url": "foo",
    "headers": {},
    "backoffPlan": {
      "type": "fixed",
      "maxRetries": 0,
      "baseDelaySeconds": 1
    }
  },
  "keypadInputPlan": {
    "enabled": true,
    "timeoutSeconds": 42,
    "delimiters": "#"
  },
  "backgroundDenoisingEnabled": false
}
```

**Response:**

```json
{
  "InsertedID": "507f1f77bcf86cd799439011",
  "Acknowledged": true
}
```

#### GET /contacts/org
Retrieve all contacts for an organization.

**Headers:**
- `Authorization: Bearer <clerk_jwt_token>` (required)

**Response:**
```json
[
  {
    "id": "507f1f77bcf86cd799439011",
    "name": "John Doe",
    "email": "john.doe@example.com",
    "phone_number": "+1234567890",
    "company": "Acme Corp",
    "position": "Manager",
    "address": "123 Main St, City, State",
    "metadata": {
      "preferences": "morning_calls",
      "notes": "Prefers calls before 10 AM"
    }
  }
]
```

#### GET /phone_numbers/org
Retrieve all phone numbers for an organization.

**Headers:**
- `Authorization: Bearer <clerk_jwt_token>` (required)

**Response:**
```json
[
  {
    "id": "507f1f77bcf86cd799439011",
    "name": "Main Office Line",
    "phone_number_id": "phone_0987654321fedcba",
    "phone_number": "+1987654321"
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
}
```

### Customer
```go
type Customer struct {
    PhoneNumber string // Customer's phone number (E.164 format)
    DayNumber   int    // Day of month for scheduling
    MonthNumber int    // Month for scheduling (1-12)
    YearNumber  int    // Year for scheduling
}
```

### Contact
```go
type Contact struct {
    Id          bson.ObjectID            // Unique MongoDB ObjectID
    Name        string                   // Full name of the contact
    Email       string                   // Contact's email address
    PhoneNumber string                   // Contact's phone number (E.164 format)
    Company     string                   // Company name
    Position    string                   // Job title or position
    Address     string                   // Physical address
    Metadata    map[string]interface{}   // Flexible field for additional data
}
```

### Assistant
```go
type Assistant struct {
    Id              bson.ObjectID // Unique MongoDB ObjectID
    Name            string        // Human-readable name for the assistant
    VapiAssistantId string        // Unique identifier in VapiAI
    Type            string        // Category or purpose of the assistant
}
```

### PhoneNumber
```go
type PhoneNumber struct {
    Id            bson.ObjectID // Unique MongoDB ObjectID
    Name          string        // Human-readable name for the phone number
    PhoneNumberId string        // Unique identifier in VapiAI
    PhoneNumber   string        // Actual phone number (E.164 format)
}
```

## Campaign Types

- `recurrent_weekly`: Runs on a weekly basis
- `recurrent_monthly`: Runs on a monthly basis
- `recurrent_yearly`: Runs on a yearly basis
- `one_time`: Runs only once on a specific date

## Campaign Statuses

- `active`: Campaign is currently running
- `paused`: Campaign is temporarily stopped
- `completed`: Campaign has finished all scheduled calls
- `cancelled`: Campaign has been permanently stopped

## Authentication

The API uses Clerk for authentication and authorization. All endpoints (except `/test`) require a valid JWT token in the Authorization header:

```
Authorization: Bearer <clerk_jwt_token>
```

The system automatically extracts the user's organization ID from the JWT token and provides organization-based data isolation.

## Error Handling

The API returns appropriate HTTP status codes:

- `200 OK`: Successful operation
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Missing or invalid authentication token
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
| `CLERK_SECRET_KEY` | Clerk secret key for authentication | Yes |

## Development

### Project Structure
```
Sarah/
├── api/                    # HTTP handlers and API endpoints
│   ├── handlers.go         # Main API handlers for all endpoints
│   └── utils.go            # Shared utility functions
├── auth/                   # Authentication and authorization
│   └── auth.go             # Clerk authentication middleware
├── clerk/                  # Clerk integration
│   └── organizations.go    # Organization management functions
├── sarah/                  # Core business logic
│   ├── campaigns.go        # Campaign management logic
│   ├── calls.go            # Call management logic
│   └── utils.go            # Business logic utilities
├── mongodb/                # Database operations
│   ├── campaigns.go        # Campaign database operations
│   ├── assistants.go       # Assistant database operations
│   ├── contacts.go         # Contact database operations
│   └── phone_numbers.go    # Phone number database operations
├── types/                  # Data type definitions
│   └── mongodb/            # MongoDB-specific types
│       ├── campaigns.go    # Campaign data structures
│       ├── assistants.go   # Assistant data structures
│       ├── contact.go      # Contact data structures
│       └── phone_numbers.go # Phone number data structures
├── main.go                 # Application entry point
├── go.mod                  # Go module file
├── go.sum                  # Go module checksums
├── README.md               # This file
└── LICENSE                 # License file
```

### Key Dependencies

- **VapiAI SDK**: For voice AI integration
- **Clerk SDK**: For authentication and user management
- **MongoDB Driver**: For database operations
- **Godotenv**: For environment variable management

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

## Contributors

- [@jpgtzg](https://github.com/jpgtzg)
