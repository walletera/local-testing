# MockServer

MockServer simulates the DinoPay external payment processor API. It is configured via `initialization.json`, which is loaded on startup.

## Initialization file

`initialization.json` contains an array of expectations. Each expectation has:

- `id` — unique string identifier
- `httpRequest` — request matcher (method, path, body)
- `httpResponseTemplate` — Velocity template producing the response
- `priority` — higher value wins when multiple expectations match the same request
- `timeToLive` / `times` — set both to `{ "unlimited": true }` for persistent expectations

## Body matching

Expectations use `JSON_SCHEMA` body matching to match on numeric amount ranges:

```json
"body": {
  "type": "JSON_SCHEMA",
  "jsonSchema": "{\"type\":\"object\",\"properties\":{\"amount\":{\"type\":\"number\",\"minimum\":200}},\"required\":[\"amount\"]}"
}
```

Supported JSON Schema draft-4 keywords for numeric constraints: `minimum`, `maximum`, `exclusiveMinimum` (boolean).

## Response templates

Templates use Apache Velocity. The `$!uuid` variable is available for generating UUIDs. Dynamic request body access via `$!bodyAsJson` does **not** work — hardcode values instead.

Template format (the whole string is parsed as a response descriptor):

```
{
    "statusCode" : 201,
    "headers" : { "content-type" : [ "application/json" ] },
    "body" : { ... }
}
```

In `initialization.json` the template is a JSON string with `\n` for newlines and `\"` for quotes.

## Current routing logic (`POST /dinopay/payments`)

| Priority | JSON Schema constraint | Response `status` |
|----------|------------------------|-------------------|
| 2 (highest) | `amount >= 200` | `rejected` |
| 1 | `amount > 100` | `pending` |
| 0 | `amount <= 100` | `confirmed` |

Priority resolves overlapping matches: a request with amount 250 matches both priority 2 and priority 1, but priority 2 wins.

## Adding or modifying expectations

1. Edit `initialization.json`.
2. Restart MockServer: `docker compose restart mockserver` from the `local-testing` directory.
3. MockServer loads the initialization file fresh on each startup.

Reference docs:
- Expectations: https://www.mock-server.com/mock_server/creating_expectations.html
- Response templates: https://www.mock-server.com/mock_server/response_templates.html
- Initialization file: https://www.mock-server.com/mock_server/initializing_expectations.html
