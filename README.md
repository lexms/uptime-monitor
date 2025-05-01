# Uptime Monitor

A SaaS uptime monitoring service to track the availability of your websites and APIs. Built with Go backend and Next.js TypeScript frontend.

## Project Structure

```
.
├── api/                # Go backend
│   ├── cmd/            # Command line applications
│   │   └── server/     # Main server application
│   └── pkg/            # Reusable packages
│       ├── handler/    # HTTP handlers
│       ├── middleware/ # HTTP middleware
│       ├── models/     # Data models
│       └── service/    # Business logic
│
└── web/                # Next.js frontend
```

## Features

- Monitor the uptime of websites and APIs
- Configure monitoring intervals and expected responses
- Get notifications when services are down
- View historical uptime data and analytics
- User authentication and multi-user support

## Tech Stack

### Backend
- Go
- Gin web framework
- JWT for authentication

### Frontend
- Next.js
- TypeScript
- TailwindCSS

## Getting Started

### Prerequisites

- Go 1.21+
- Node.js 18+
- pnpm

### Running the Backend

```bash
cd api
go mod tidy
go run cmd/server/main.go
```

### Running the Frontend

```bash
cd web
pnpm install
pnpm dev
```

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Log in a user

### Endpoints (Protected)
- `GET /api/endpoints` - Get all endpoints for current user
- `GET /api/endpoints/:id` - Get a specific endpoint
- `POST /api/endpoints` - Create a new endpoint to monitor
- `PUT /api/endpoints/:id` - Update an existing endpoint
- `DELETE /api/endpoints/:id` - Delete an endpoint

## Deployment

The frontend can be deployed to Vercel, and the backend can be deployed to any platform that supports Go applications.

## License

This project is licensed under the MIT License. 