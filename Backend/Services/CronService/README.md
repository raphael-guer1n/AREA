# Cron Service

Cron/Timer action scheduler service for AREA. It manages timer-based actions that trigger at specified intervals, storing action configurations in PostgreSQL and using the robfig/cron library for scheduling.

## Quick Start

```bash
# Build and start containers
docker-compose up -d
```

The API should be accessed through the gateway at `http://localhost:8080/area_cron_api`.

## API Endpoints

- **GET** `/health` - Check service health
- **POST** `/actions` - Create cron/timer actions (internal only)
- **POST** `/activate/{actionId}` - Activate a cron action (internal only)
- **POST** `/deactivate/{actionId}` - Deactivate a cron action (internal only)
- **DELETE** `/actions/{actionId}` - Delete a cron action (internal only)

## Configuration

Environment variables can be configured in `.env`:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=cron_service_db
SERVER_PORT=8086
INTERNAL_SECRET=secret
AREA_SERVICE_URL=http://gateway:8080/area_area_api
LOG_ALL_REQUESTS=false
```

## How It Works

1. **Action Creation**: When a timer action is created via `POST /actions`, the service:
   - Stores the action in PostgreSQL
   - Parses the `delay` input field (in seconds)
   - If `active: true`, schedules a cron job using `@every Ns` syntax

2. **Timer Execution**: When a scheduled timer fires:
   - The service calls `POST /triggerArea` on the AreaService
   - Sends the `action_id` and `output_fields` (containing the delay value)

3. **Action Management**:
   - `/activate/{actionId}` - Starts the timer for an inactive action
   - `/deactivate/{actionId}` - Stops the timer for an active action
   - `/actions/{actionId}` (DELETE) - Removes the action and stops its timer

## Action Types

### 1. Delay Action (timer_delay / delay_action)

Triggers every N seconds.

```json
{
  "actions": [
    {
      "active": true,
      "action_id": 1,
      "type": "cron",
      "provider": "",
      "service": "timer",
      "title": "delay_action",
      "input": [
        {
          "name": "delay",
          "value": "10"
        }
      ]
    }
  ]
}
```

The `delay` value is in seconds. For example, `"10"` means the action will trigger every 10 seconds.

### 2. Daily Action

Triggers every day at a specific time (hour and minute).

```json
{
  "actions": [
    {
      "active": true,
      "action_id": 2,
      "type": "cron",
      "provider": "",
      "service": "timer",
      "title": "daily_action",
      "input": [
        {
          "name": "hour",
          "value": "9"
        },
        {
          "name": "minute",
          "value": "30"
        }
      ]
    }
  ]
}
```

This example triggers every day at 9:30 AM.

### 3. Weekly Action

Triggers once a week on a specific day and time.

```json
{
  "actions": [
    {
      "active": true,
      "action_id": 3,
      "type": "cron",
      "provider": "",
      "service": "timer",
      "title": "weekly_action",
      "input": [
        {
          "name": "day_of_week",
          "value": "1"
        },
        {
          "name": "hour",
          "value": "14"
        },
        {
          "name": "minute",
          "value": "0"
        }
      ]
    }
  ]
}
```

Day of week: 0=Sunday, 1=Monday, ..., 6=Saturday. This example triggers every Monday at 2:00 PM.

### 4. Monthly Action

Triggers once a month on a specific day and time.

```json
{
  "actions": [
    {
      "active": true,
      "action_id": 4,
      "type": "cron",
      "provider": "",
      "service": "timer",
      "title": "monthly_action",
      "input": [
        {
          "name": "day_of_month",
          "value": "1"
        },
        {
          "name": "hour",
          "value": "0"
        },
        {
          "name": "minute",
          "value": "0"
        }
      ]
    }
  ]
}
```

This example triggers on the 1st of every month at midnight.

## Trigger Format

When a timer fires, the CronService sends this to AreaService:

```json
{
  "action_id": 1,
  "output_fields": [
    {
      "name": "delay",
      "value": "10"
    }
  ]
}
```

## Database Schema

The service uses a single table `cron_actions`:

- `action_id` (PRIMARY KEY) - Unique identifier for the action
- `active` (BOOLEAN) - Whether the action is currently scheduled
- `type` (VARCHAR) - Action type (always "cron")
- `provider` (VARCHAR) - Provider name (usually empty for timers)
- `service` (VARCHAR) - Service name (e.g., "timer")
- `title` (VARCHAR) - Action title
- `input` (JSONB) - Input fields including delay configuration
- `cron_job_id` (INTEGER) - Internal cron scheduler entry ID
- `created_at`, `updated_at` (TIMESTAMP) - Timestamps

## Persistence & Restart Behavior

When the service restarts:
- It reconnects to the database
- Loads all active actions from `cron_actions` table
- Re-schedules all active timers automatically

This ensures timers continue running across service restarts.
