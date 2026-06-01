# Architecture

CnzAMnt is a small mobile-first web app with a Vue frontend, Go backend, and SQLite database.

The architecture should stay boring and easy to understand. The MVP should optimize for a working feedback loop, not infrastructure depth.

## App Shape

The app has three main parts:

- `frontend`: Vue 3 + TypeScript mobile-first web app.
- `backend`: Go HTTP JSON API.
- `database`: SQLite.

Docker Compose should run the frontend and backend together for local development. The backend owns the SQLite database.

## Suggested Repository Shape

```text
.
├── backend/
│   ├── cmd/server/
│   └── internal/
│       ├── api/
│       ├── db/
│       ├── models/
│       ├── seed/
│       └── service/
├── frontend/
│   └── src/
│       ├── components/
│       ├── services/
│       ├── types/
│       └── views/
├── docs/
├── docker-compose.yml
├── Makefile
└── README.md
```

This mirrors the simple parts of the Rookery reference project without bringing over its extra workflow systems.

## MVP Data Model

The first database model can stay small:

- `users`
- `artworks`
- `comments`
- `ratings` or `feedback`
- `cnz_transactions`

The MVP can combine comments and ratings into one `feedback` table if that keeps the code simpler.

Suggested concepts:

- A user has a CNZ balance.
- An artwork belongs to one artist.
- A feedback record belongs to one artwork and one reviewer.
- A feedback record includes comment text and a CNZ rating from 1 to 5.
- A transaction records CNZ spent by the reviewer and CNZ earned by the artist.

## Backend

The backend should expose a small JSON API.

Likely first endpoints:

- `GET /api/health`
- `GET /api/me`
- `GET /api/artworks`
- `POST /api/artworks`
- `POST /api/artworks/{id}/feedback`

Because there is no real login yet, the backend can use a fake current user. A request header such as `X-Cnz-User-ID` may be useful later, but the first version can keep this even simpler if needed.

## Frontend

The frontend should be mobile-first from the start.

Core screens:

- Main artwork feed.
- Artwork detail with comments and ratings.
- Simple post-artwork form.
- Current user CNZ balance display.

The UI should feel like the app itself, not a marketing page. The first screen should help the user view art, post art, or give feedback.

## Docker And Makefile

The first local setup should support:

```bash
make dev
```

Future useful targets:

```bash
make backend
make frontend
make test
make test-backend
make test-frontend
make down
```

Docker Compose should run:

- backend service
- frontend service
- SQLite storage through a local file or Docker volume

## Constraints

Keep out of the MVP:

- payments
- real login
- AI API calls
- native apps
- multi-tenant architecture
- complex permissions
- background workers
- notification systems

Add those only after the core art feedback loop works.
