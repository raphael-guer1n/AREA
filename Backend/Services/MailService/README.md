# MailService

Internal SMTP mail sender for AREA reactions. Exposes a simple HTTP API consumed by AreaService via the gateway.

## Endpoints

- **GET** `/health` — Health check
- **POST** `/send` — Send an email

Request body:
```json
{
  "to": "user@example.com",
  "subject": "Notification from AREA",
  "body": "Hello from AREA"
}
```

`to` can also be an array of recipients.

## Configuration

Set these in `.env` (not committed):

```
SERVER_PORT=8088
SMTP_HOST=ssl0.ovh.net
SMTP_PORT=465
SMTP_USERNAME=no-reply@area-connect.cloud
SMTP_PASSWORD=...
SMTP_FROM=no-reply@area-connect.cloud
SMTP_FROM_NAME=AREA Connect
SMTP_SECURITY=ssl
```

## HTML template

The email layout lives at:
`app/templates/email.html`

## Local run

```bash
make docker-up
```
