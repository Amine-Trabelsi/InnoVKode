# Inno-VKode

## Project overview
Inno-VKode is an alpha version of a Max messenger assistant for universities.  
It combines a FastAPI backend (`server-be`) with a Go-based Max bot (`bot-be`)
to help students, professors, and guests work with academic information:

- Students can view schedules, request certificates, read lecture summaries, take quizzes,
  ask questions via RAG-powered search, and track dorm / visa requests.
- Professors upload lecture material, generate quizzes, and respond to student questions.
- Guests explore programs, admission deadlines, campus events, and other high-level details.

The backend exposes APIs for account onboarding, OTP verification, visa and dorm workflows,
library tooling, and horizontal-scaling friendly academic services.  
The bot constrains interactions to button-driven flows so that every user cohort receives an
accessible, guided experience directly in Max.

## Repository layout
| Path        | Description                                                                        |
|-------------|------------------------------------------------------------------------------------|
| `server-be` | FastAPI application that owns persistence, domain rules, and seeding (`main.py`).  |
| `bot-be`    | Go Max bot that renders the button-based UX and calls the backend APIs.       |
| `test`      | Utilities and fixtures used to validate bot and backend interactions.              |

A Dockerized Postgres instance provides the data layer. All components are orchestrated through
`docker-compose.yaml`, so you can bring the entire stack up with a single command.

## Running the stack with Docker
### 1. Prerequisites
1. Docker Engine 24+ and Docker Compose (V2).  
2. A Max bot token.

> The compose file ships with placeholder environment variables (for example `MAX_BOT_TOKEN`); override
> them before going to production.

### 2. Start services
From the repository root:

```bash
docker compose up --build
```

This will:
1. Build the FastAPI backend (`server-be`) and expose it on `http://localhost:8000`.
2. Build the Go bot (`bot-be`). The bot connects to Max using the token supplied through
   `MAX_BOT_TOKEN`, then talks to the backend through the internal Docker network.
3. Launch Postgres 15 with the credentials from `docker-compose.yaml`. Data is stored inside the
   named volume `postgres_data`.

> The backend waits until Postgres is healthy, then automatically creates tables and seeds initial data.

### 3. Verify
- Visit `http://localhost:8000/docs` to inspect and exercise the FastAPI endpoints.
- Use the Max client that is linked to your bot token to start a conversation and walk through
  the Student / Professor / Guest flows.

### 4. Stopping and cleaning up
- Stop the stack: `docker compose down`
- Stop and remove the persistent volume as well (if you want a clean database):  
  `docker compose down -v`

## Environment configuration
You can override any variable declared in `docker-compose.yaml` either by editing the file or by
creating a `.env` file next to it. Common overrides:

| Variable            | Purpose                                                       |
|---------------------|---------------------------------------------------------------|
| `MAX_BOT_TOKEN`     | Max bot token – never commit your production value.      |
| `BACKEND_BASE_URL`  | Internal URL that the bot uses to call the FastAPI service.   |
| `RESET_DB_ON_STARTUP` | Drop & recreate schema each start; set to `false` in prod.  |

Example `.env` snippet:

```env
MAX_BOT_TOKEN=123456:ABCDEF-your-token
BACKEND_BASE_URL=http://server-be:8000
RESET_DB_ON_STARTUP=false
```

> Compose automatically loads `.env` in the same directory, so no additional CLI flags are needed.

## Troubleshooting tips
- **Database keeps resetting** – ensure `RESET_DB_ON_STARTUP=false` when you need persistent data.
- **Bot cannot reach backend** – confirm both services are running and the compose network name
  matches `BACKEND_BASE_URL`.
- **Token errors** – double-check the `MAX_BOT_TOKEN` value and that the bot is not running elsewhere.
