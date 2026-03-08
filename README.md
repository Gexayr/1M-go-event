# Risk Monitoring Event API

This API allows external systems to send security-related events to the Risk Monitoring service.  
Each event is analyzed by the risk scoring engine and may trigger alerts or appear in the dashboard.

---

# Base URL


http://localhost:8080


Production example:


https://api.yourdomain.com


---

# Authentication

Requests must include an API key in the header.


X-API-KEY: your_api_key


Example:


X-API-KEY: client_test_key


---

# Endpoint

## Register Event


POST /events


Registers a new event in the system.

The event will be:

- stored in the database
- analyzed by the risk scoring engine
- potentially trigger alerts
- visible in the dashboard

---

# Request Headers

| Header | Required | Description |
|------|------|------|
| Content-Type | Yes | application/json |
| X-API-KEY | Yes | Client API key |

Example:


Content-Type: application/json
X-API-KEY: client_test_key


---
# Events API

## Example Payload

```json
{
  "client_id": "casino_123",
  "event_type": "failed_login",
  "timestamp": "2026-03-05T18:00:00Z",
  "metadata": {
    "user_id": "u7781",
    "ip": "192.168.1.1"
  }
}
```

---

## Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `client_id` | string | Yes | Identifier of the client system |
| `event_type` | string | Yes | Type of event |
| `timestamp` | ISO8601 | No | Event time (server time used if missing) |
| `metadata` | object | No | Additional event information |

---

## Example Event Types

Common event types include:

- `login`
- `failed_login`
- `withdrawal`
- `deposit`
- `account_update`
- `password_change`
- `suspicious_ip`
- `multiple_login_attempts`

> Custom event types are also supported.

---

## Example Request

**cURL**

```bash
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -H "X-API-KEY: client_test_key" \
  -d '{
    "client_id": "casino_123",
    "event_type": "failed_login",
    "metadata": {
      "user_id": "u7781",
      "ip": "192.168.1.1"
    }
  }'
```

---

## Responses

### Successful Response

**`200 OK`**

```json
{
  "status": "success",
  "event_id": 1042,
  "timestamp": "2026-03-05T18:00:01Z"
}
```

### Error Responses

**Invalid Request — `400 Bad Request`**

```json
{
  "error": "invalid request body"
}
```

**Unauthorized — `401 Unauthorized`**

```json
{
  "error": "invalid api key"
}
```

**Internal Error — `500 Internal Server Error`**

```json
{
  "error": "internal server error"
}
```

---

## Event Processing Flow

```
Client System
      ↓
POST /events
      ↓
Event Stored in Database
      ↓
Risk Scoring Engine
      ↓
High Risk Detection
      ↓
Telegram Alert (if score >= threshold)
      ↓
Dashboard + AI Reports
```

---

## Example Metadata

Metadata can contain any additional information relevant to the event.

```json
{
  "user_id": "user_9281",
  "ip": "192.168.0.12",
  "device": "mobile",
  "country": "DE"
}
```

---

## Rate Limits

Recommended limit: **100 events / second** per client.

If higher throughput is required, contact the system administrator.

---

## Support

For integration support, contact the platform administrator.