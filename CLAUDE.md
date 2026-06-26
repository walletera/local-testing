# local-testing

Docker Compose environment for running the full Walletera stack locally. See the parent `CLAUDE.md` for platform-wide architecture and the `README.md` for setup instructions.

## Starting and stopping

```bash
make start              # Start all services and configure Barong
make stop               # Stop all services (keeps volumes)
docker compose down -v  # Stop and remove volumes (full reset)
```

## Key ports

| Service | Port |
|---------|------|
| Envoy (API Gateway) | 3099 |
| payments | 3880 |
| Barong (auth) | 9090 |

## Running end-to-end tests

```bash
make run-dinopay-outbound-succeed-tests
make run-dinopay-inbound-succeed-tests
```

Requires the stack to be running and `BASIC_AUTH_USERNAME` / `BASIC_AUTH_PASSWORD` env vars set.

## Services

@mockserver/CLAUDE.md
