# Avions Club Backend

This is the backend server for the Avions Club website. It provides a RESTful API for managing club members, projects, blogs, and file storage.

## Features

- Authentication and Authorization
- Member Management
- Project Management
- Blog Management
- File Storage (using Supabase)
- Search Functionality

## Tech Stack

- Go (Golang)
- Gin Web Framework
- GORM (with PostgreSQL)
- Supabase Storage
- JWT Authentication

## Prerequisites

- Go 1.19 or higher
- PostgreSQL (or Supabase account)
- Git

## Environment Variables

Copy `.env.example` to `.env` and update the following variables:

```bash
# Server Configuration
PORT=8080
ENV=development

# Database Configuration
DB_HOST=your-supabase-project.supabase.co
DB_USER=postgres
DB_PASSWORD=your-password
DB_NAME=postgres
DB_PORT=5432
DB_SSL_MODE=require

# JWT Configuration
JWT_SECRET=your-jwt-secret
ADMIN_PASSWORD=your-admin-password

# Supabase Configuration
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_KEY=your-supabase-key
SUPABASE_SERVICE_KEY=your-service-key

# Storage Configuration
STORAGE_BUCKET_IMAGES=images
STORAGE_BUCKET_MARKDOWN=markdown
MAX_FILE_SIZE=5242880  # 5MB in bytes

# CORS Configuration
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/avions-club-backend.git
   cd avions-club-backend
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up environment variables:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. Run the server:
   ```bash
   go run main.go
   ```

## API Documentation

### Authentication

- `POST /api/auth/login` - Admin login

### Members

- `GET /api/members` - List all members
- `GET /api/members/:id` - Get a specific member
- `POST /api/members` - Create a member (Admin)
- `PUT /api/members/:id` - Update a member (Admin)
- `DELETE /api/members/:id` - Delete a member (Admin)

### Projects

- `GET /api/projects` - List all projects
- `GET /api/projects/:id` - Get a specific project
- `POST /api/projects` - Create a project (Admin)
- `PUT /api/projects/:id` - Update a project (Admin)
- `DELETE /api/projects/:id` - Delete a project (Admin)

### Blogs

- `GET /api/blogs` - List all blogs
- `GET /api/blogs/:id` - Get a specific blog
- `POST /api/blogs` - Create a blog (Admin)
- `PUT /api/blogs/:id` - Update a blog (Admin)
- `DELETE /api/blogs/:id` - Delete a blog (Admin)

### Storage

- `POST /api/storage/upload` - Upload a file (Admin)
- `DELETE /api/storage/:bucket/:filename` - Delete a file (Admin)

### Search

- `GET /api/search?q=query` - Search across all content

## Deployment

1. Set environment variables for production:
   ```bash
   ENV=production
   ```

2. Build the binary:
   ```bash
   go build -o server
   ```

3. Run the server:
   ```bash
   ./server
   ```

## Testing

Run the test script to verify API functionality:
```bash
./test.sh
```

## License

MIT License

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request 