# CnzAMnt

CnzAMnt is a small mobile-first web app for art feedback.

A user posts one piece of art. Other users can leave comments and spend CNZ to rate the artwork from 1 to 5. CNZ is the app's internal points/money system. Each new user starts with 5000 CNZ, and users cannot spend more CNZ than they have.

The goal is simple: help artists get useful comments while giving feedback a little weight.

## MVP

The first version should include:

- A mobile-first Vue 3 frontend.
- A Go backend API.
- SQLite persistence.
- A simple fake-user flow instead of real login.
- Seed users that each start with 5000 CNZ.
- One artwork post per user.
- A feed or main screen showing posted artwork.
- Comments on artwork.
- A 1 to 5 CNZ rating/spend action.
- Balance checks so users cannot spend CNZ they do not have.
- Artist earnings equal to 10% of CNZ spent on their artwork.
- Docker and Makefile commands for local development.

## Not In The MVP

The first version should not include:

- Real login or account recovery.
- Payments or buying CNZ.
- Native iOS or Android apps.
- Multiple artwork posts per user.
- AI-generated comments.
- AI moderation.
- Social following, direct messages, or notifications.
- Public deployment complexity.
- Admin dashboards.

These may come later, but the first build should prove the core feedback loop.

## CNZ Basics

CNZ is an internal app balance, not real money in the MVP.

- Every new user starts with 5000 CNZ.
- Users spend CNZ when they rate an artwork.
- A rating is from 1 to 5 CNZ.
- A user must have enough CNZ before rating.
- If the user does not have enough CNZ, the rating is rejected.
- The artist earns 10% of the CNZ spent on their artwork.

See [docs/cnz-rules.md](docs/cnz-rules.md) for the detailed rules.

## Comments And Ratings

A feedback action should include:

- The artwork being reviewed.
- The commenting user.
- A text comment.
- A CNZ rating from 1 to 5.

The comment gives the artist useful feedback. The CNZ rating gives the feedback weight and transfers a small artist earning.

## AI Feedback Later

AI feedback is intentionally later.

The likely future shape is:

- A user asks for AI feedback on their artwork.
- The app sends the artwork context and existing comments to an AI service.
- AI returns a helpful critique or summary.
- The artist can compare human feedback and AI feedback.

For now, CnzAMnt should only store human comments and ratings. No AI API calls should be added in the MVP.

## Local Development

Once the app is built, local development should work with:

```bash
make dev
```

Expected local services:

- Frontend: Vue 3 + TypeScript, likely at `http://localhost:5173`
- Backend: Go API, likely at `http://localhost:8080`
- Database: SQLite file stored locally or in a Docker volume

Useful future commands:

```bash
make backend
make frontend
make test
make down
```

Backend foundation commands now available:

```bash
make backend-dev
make backend-test
make backend-build
```

Frontend and Docker commands now available:

```bash
make frontend-dev
make frontend-build
make dev
make up
make down
make logs
make logs-backend
make logs-frontend
```

## Project Docs

- [Architecture](docs/architecture.md)
- [CNZ Rules](docs/cnz-rules.md)
