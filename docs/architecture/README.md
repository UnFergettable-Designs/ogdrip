# System Architecture

OG Drip is designed as a modern, scalable monorepo application with clear separation of concerns
between frontend, backend, and shared components.

## High-Level Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│                 │    │                 │    │                 │
│   Frontend      │    │   Backend       │    │   Database      │
│   (Astro +      │◄──►│   (Go +         │◄──►│   (SQLite)      │
│   Svelte 5)     │    │   ChromeDP)     │    │                 │
│                 │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Static Files  │    │   Generated     │    │   Generation    │
│   & Assets      │    │   Images        │    │   History       │
│                 │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Component Overview

### Frontend (Astro + Svelte 5)

- **Purpose**: User interface for Open Graph image generation
- **Technology**: Astro for static site generation, Svelte 5 for reactive components
- **Location**: `frontend/` directory
- **Port**: 3000 (development), served via reverse proxy in production

### Backend (Go + ChromeDP)

- **Purpose**: API server and image generation engine
- **Technology**: Go with ChromeDP for headless browser automation
- **Location**: `backend/` directory
- **Port**: 8888

### Database (SQLite)

- **Purpose**: Store generation history, metadata, and configuration
- **Technology**: SQLite for simplicity and portability
- **Location**: `backend/data/` directory

### Shared Types

- **Purpose**: Common TypeScript types and utilities
- **Technology**: TypeScript
- **Location**: `shared/` directory

## Data Flow

### Image Generation Process

1. **User Request**: User submits URL via frontend or API
2. **Validation**: Backend validates URL and parameters
3. **Browser Launch**: ChromeDP launches headless browser
4. **Page Load**: Browser navigates to target URL
5. **Metadata Extraction**: Extract title, description, and other metadata
6. **Screenshot**: Capture page screenshot with specified dimensions
7. **Image Processing**: Optimize and save generated image
8. **Database Storage**: Store generation record and metadata
9. **Response**: Return image URL and metadata to client

### Request Flow Diagram

```
Client Request
      │
      ▼
┌─────────────┐
│  Frontend   │
│   (Astro)   │
└─────────────┘
      │
      ▼ HTTP API Call
┌─────────────┐
│  Backend    │
│   (Go API)  │
└─────────────┘
      │
      ▼
┌─────────────┐
│  ChromeDP   │
│  (Browser)  │
└─────────────┘
      │
      ▼
┌─────────────┐
│  Target     │
│  Website    │
└─────────────┘
```

## Technology Stack

### Frontend Stack

- **Astro**: Static site generator with partial hydration
- **Svelte 5**: Reactive UI framework with runes
- **TypeScript**: Type-safe JavaScript
- **Vite**: Build tool and development server
- **Tailwind CSS**: Utility-first CSS framework

### Backend Stack

- **Go**: High-performance system programming language
- **ChromeDP**: Go library for Chrome DevTools Protocol
- **Gin**: HTTP web framework for Go
- **SQLite**: Embedded SQL database
- **Swagger**: API documentation generation

### Development & Build Tools

- **Turborepo**: Monorepo build system
- **pnpm**: Fast, disk space efficient package manager
- **Docker**: Containerization platform
- **ESLint/Prettier**: Code linting and formatting
- **Vitest**: Unit testing framework

## Deployment Architecture

### Docker Deployment

```
┌─────────────────────────────────────────┐
│              Docker Container           │
│                                         │
│  ┌─────────────┐    ┌─────────────┐    │
│  │  Frontend   │    │  Backend    │    │
│  │  (Static)   │    │  (Go API)   │    │
│  └─────────────┘    └─────────────┘    │
│                                         │
│  ┌─────────────────────────────────┐    │
│  │         Nginx Proxy             │    │
│  └─────────────────────────────────┘    │
└─────────────────────────────────────────┘
```

### Coolify Deployment

```
┌─────────────────────────────────────────┐
│              Coolify Platform           │
│                                         │
│  ┌─────────────────────────────────┐    │
│  │         nixpacks Build          │    │
│  └─────────────────────────────────┘    │
│                    │                    │
│                    ▼                    │
│  ┌─────────────────────────────────┐    │
│  │      Application Container      │    │
│  └─────────────────────────────────┘    │
│                    │                    │
│                    ▼                    │
│  ┌─────────────────────────────────┐    │
│  │       Reverse Proxy + SSL       │    │
│  └─────────────────────────────────┘    │
└─────────────────────────────────────────┘
```

## Security Architecture

### Authentication & Authorization

- **Admin API**: Bearer token authentication
- **Public API**: Rate limiting and input validation
- **CORS**: Configured for frontend-backend communication

### Data Security

- **Input Validation**: All user inputs validated and sanitized
- **SQL Injection Prevention**: Prepared statements used throughout
- **File System Security**: Generated files stored in isolated directory
- **Environment Variables**: Sensitive configuration stored securely

## Performance Considerations

### Optimization Strategies

- **Browser Instance Pooling**: Reuse ChromeDP instances when possible
- **Image Optimization**: Compress generated images
- **Caching**: HTTP caching headers for static assets
- **Database Indexing**: Optimize queries with proper indexes

### Scalability

- **Horizontal Scaling**: Stateless design allows multiple instances
- **Load Balancing**: Frontend and backend can be load balanced
- **Database**: SQLite suitable for moderate loads, can migrate to PostgreSQL
- **CDN Integration**: Static assets can be served via CDN

## Monitoring & Observability

### Logging

- **Structured Logging**: JSON format for log aggregation
- **Log Levels**: Debug, Info, Warn, Error
- **Request Tracing**: Track requests across components

### Metrics

- **Application Metrics**: Request count, response time, error rate
- **System Metrics**: CPU, memory, disk usage
- **Business Metrics**: Generation count, user activity

### Health Checks

- **API Health**: `/api/health` endpoint
- **Database Health**: Connection and query validation
- **Browser Health**: ChromeDP instance validation

## Further Reading

- [Frontend Architecture](frontend.md) - Detailed frontend implementation
- [Backend Architecture](backend.md) - Detailed backend implementation
- [Database Schema](database.md) - Database structure and relationships
- [Security Architecture](../security/README.md) - Security implementation details
