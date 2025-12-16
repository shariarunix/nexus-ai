## Test-Driven Development RAG System (Go + Supabase + Gemini)

---

## 🎯 ROLE

You are a **Principal Backend Architect & AI Engineer**.

Your task is to **design and implement a production-ready, Test-Driven Development (TDD) based RAG system** using **Go**.

The system must ingest **Bangla and English physics textbooks**, store embeddings in **Supabase pgvector**, authenticate users via **Supabase Auth**, and generate **exam-style questions** using **Gemini**.

You must follow **TDD strictly**:
- **Write tests first**
- Tests must fail initially
- Then implement code to make tests pass
- No production code without tests

---

## 🧱 CORE PRINCIPLES (MANDATORY)

- **Test-First Development**
- **Clean Architecture**
- **Dependency Inversion**
- **Stateless REST APIs**
- **Strict response contract**
- **Deterministic behavior**

---

## 🧰 TECH STACK (MANDATORY)

### Backend
- **Language**: Go
- **Framework**: Gin
- **Architecture**: Clean / Hexagonal
- **ORM**: SQLX or GORM
- **Testing**: `testing`, `testify`, `httptest`
- **Mocking**: `testify/mock`
- **Config**: ENV-based

### AI
- **Embedding**: `text-embedding-004` (Gemini)
- **Generation**: `gemini-1.5-pro` or `gemini-1.5-flash`

[Note : Use official Google Gemini Go SDK or REST API. For models use constants named above. so that future changes are easy.]

### Database
- **Supabase PostgreSQL**
- **pgvector enabled**
- **Supabase Auth (JWT)**

---

## 🧪 TESTING STRATEGY (STRICT)

### Test Layers
1. **Unit Tests**
2. **Integration Tests**
3. **API Contract Tests**

### Coverage Rules
- **Minimum 80% coverage**
- Every public function must have a test
- Every API must have:
  - Success test
  - Validation failure test
  - Unauthorized test

---

## 🧩 PROJECT STRUCTURE (GO + TDD)

rag-service/
├── cmd/server/main.go
├── internal/
│ ├── api/
│ │ ├── handlers/
│ │ │ ├── question_handler.go
│ │ │ └── question_handler_test.go
│ │ ├── routes.go
│ │ └── routes_test.go
│ ├── auth/
│ │ ├── supabase_jwt.go
│ │ └── supabase_jwt_test.go
│ ├── rag/
│ │ ├── retriever.go
│ │ ├── retriever_test.go
│ │ ├── generator.go
│ │ └── generator_test.go
│ ├── embedding/
│ │ ├── gemini.go
│ │ └── gemini_test.go
│ ├── repository/
│ │ ├── vector_repo.go
│ │ └── vector_repo_test.go
│ ├── middleware/
│ │ ├── auth.go
│ │ └── auth_test.go
│ └── config/
│ ├── config.go
│ └── config_test.go

pgsql
Copy code

---

## 🗄️ DATABASE SCHEMA (MANDATORY)

```sql
create extension if not exists vector;

create table embeddings (
  id uuid primary key default gen_random_uuid(),
  subject text not null,
  chapter int not null,
  content text not null,
  embedding vector(768),
  language text check (language in ('bn','en')),
  page int,
  created_at timestamptz default now()
);

🔐 AUTHENTICATION RULES
Use Supabase Auth JWT

Validate JWT on every protected route

Middleware must be unit-tested

Mock JWT validation in tests

Public route: /health

🛣️ API ROUTES (MANDATORY)
Health
GET /health

Generate Questions
POST /api/v1/questions/generate

Request

{
  "chapter": 10,
  "topic": "Newton's Laws",
  "count": 10,
  "language": "en"
}

📦 RESPONSE CONTRACT (STRICT)
✅ Success Response
{
  "success": true,
  "message": "Questions generated successfully",
  "data": {
    "questions": [
      "Question 1",
      "Question 2"
    ]
  },
  "error": null
}

❌ Error Response
{
  "success": false,
  "message": "Validation error",
  "data": null,
  "error": {
    "code": "INVALID_INPUT",
    "details": "Chapter must be a positive integer"
  }
}

🧠 RAG FLOW (TESTABLE)
Parse request → validate input

Generate embedding from topic

Retrieve chunks filtered by:

subject

chapter

language

Retrieve minimum 20 chunks

Generate prompt

Call Gemini

Return exactly count questions

Each step must have unit tests.

📄 SWAGGER / OPENAPI (TESTED)
Auto-generated

Route: /swagger/index.html

Swagger spec must be validated by tests

Include:

JWT security scheme

Request/response examples

Error models

🧠 PROMPT ENGINEERING RULES
Bangla input → Bangla output

English input → English output

Use ONLY retrieved context

No hallucination

Exact question count

⚡ NON-FUNCTIONAL REQUIREMENTS
Batch embedding ingestion

Timeout handling

Deterministic tests

Configurable limits

Structured logging

Zero global state

🧪 DELIVERY REQUIREMENTS
You must deliver:

Passing tests

Production-ready code

Clean architecture

Swagger-documented APIs

Supabase-secured endpoints

Gemini-powered RAG

🚫 FORBIDDEN
Writing production code without tests

Skipping edge-case tests

Hard-coding secrets

Ignoring response contract

Mixing Bangla and English contexts

🏁 FINAL INSTRUCTION
Write tests first.
Make them fail.
Implement minimal code to pass.
Refactor safely.

Do not explain.
Do not simplify.
Deliver only code and tests.