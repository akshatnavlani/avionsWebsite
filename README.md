# Avions Club Website

A modern website for the Avions Club, built with Next.js and Go.

## Features

- Modern, responsive design
- Project showcase
- Blog system
- Member profiles
- Admin dashboard
- File storage integration
- Authentication system

## Tech Stack

### Frontend
- Next.js 14
- TypeScript
- Tailwind CSS
- Shadcn UI
- React Server Components

### Backend
- Go
- Gin web framework
- Supabase for storage
- PostgreSQL database

## Prerequisites

- Node.js 18+
- Go 1.21+
- PostgreSQL
- Supabase account

## Environment Variables

### Frontend (.env.local)
```env
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_SUPABASE_URL=your_supabase_url
NEXT_PUBLIC_SUPABASE_ANON_KEY=your_supabase_anon_key
```

### Backend (.env)
```env
PORT=8080
DATABASE_URL=your_database_url
SUPABASE_URL=your_supabase_url
SUPABASE_SERVICE_KEY=your_supabase_service_key
JWT_SECRET=your_jwt_secret
ALLOWED_ORIGINS=http://localhost:3000
```

## Development

1. Clone the repository:
```bash
git clone https://github.com/yourusername/avions-club-website.git
cd avions-club-website
```

2. Install frontend dependencies:
```bash
cd frontend
pnpm install
```

3. Install backend dependencies:
```bash
cd ../backend
go mod download
```

4. Start the development servers:

Frontend:
```bash
cd frontend
pnpm dev
```

Backend:
```bash
cd backend
go run main.go
```

The website will be available at http://localhost:3000

## Production Deployment

### Frontend Deployment (Vercel)

1. Push your code to GitHub
2. Connect your repository to Vercel
3. Configure environment variables in Vercel
4. Deploy

### Backend Deployment

1. Build the Go binary:
```bash
cd backend
GOOS=linux GOARCH=amd64 go build -o avions-club-backend
```

2. Deploy to your preferred hosting service (e.g., DigitalOcean, AWS, etc.)

3. Set up environment variables on your hosting service

4. Run the backend service:
```bash
./avions-club-backend
```

### Database Setup

1. Create a PostgreSQL database
2. Run migrations:
```bash
cd backend
go run scripts/migrate.go
```

### Storage Setup

1. Create a Supabase project
2. Create storage buckets for:
   - projects
   - blogs
   - members
3. Configure CORS policies in Supabase

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 