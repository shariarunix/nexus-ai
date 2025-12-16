## NEXUS-AI

## Features
- **Gin Framework**: Fast, lightweight, and expressive HTTP routing.
- **Supabase Integration**:
  - Direct PostgreSQL access via `DATABASE_URL` (with pgvector enabled).
  - Admin operations using Supabase Service Role Key (bypasses RLS).
  - JWT verification using Supabase JWT secret.
- **Google Gemini API**: Access to powerful multimodal generative models (text, vision, etc.).
- **Environment-based Configuration**: Clean `.env` setup.
- **Production Ready**: Structured logging, graceful shutdown, and easy deployment.

## Tech Stack
- Go (1.22+ recommended)
- Gin Gonic (web framework)
- Supabase (Postgres + pgvector)
- Google Gemini API
- godotenv (for .env loading)
- pgx (PostgreSQL driver)

## Getting Started

### Prerequisites
- Go 1.22 or higher
- A Supabase project with pgvector enabled [](https://supabase.com)
- Google AI Studio account for Gemini API key [](https://ai.google.dev)

### Installation
1. Clone the repository:
   ```bash
   [git clone https://github.com/yourusername/your-go-gin-supabase-gemini-project.git
   cd your-go-gin-supabase-gemini-project](https://github.com/shariarunix/nexus-ai)
   ```

Initialize Go module (if not already done): ```go mod tidy```
Create a .env file in the root directory:

# Env
PORT=8080
ENV=development

# Supabase
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your-anon-public-key
SUPABASE_SERVICE_ROLE_KEY=your-service-role-secret-key
SUPABASE_JWT_SECRET=your-jwt-secret

# Database (pgvector enabled)
DATABASE_URL=postgres://postgres.your-project:[YOUR-PASSWORD]@db.your-project.supabase.co:5432/postgres?sslmode=require

# Gemini
GEMINI_API_KEY=your-gemini-api-keyImportant Security Notes:
