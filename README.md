# Walletera Local Development Environment

This repository provides a local development and testing environment for the Walletera platform. It orchestrates multiple microservices and supporting infrastructure using Docker Compose, making it easy to spin up, configure, and test the system end-to-end.

## Features

- **Rapid Setup & Teardown:** Easily start, stop, and reset all infrastructure and app services with simple Makefile targets.
- **Seeded Test Data:** Automated database seeding for authentication service.
- **Custom Test Runner:** Build and run Go-based end-to-end test scenarios to validate business and infrastructure logic.
- **Flexible Service Configuration:** Each service can be (re)built and configured independently to support rapid development and debugging.

## Getting Started

### Prerequisites

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Go](https://go.dev/) (for running tests)
- GNU Make

### Quick Start

1. **Start All Services:**

```shell script
make start
```

This command will:
- Bring up all Docker Compose services
- Configure Barong (authentication service) database and seed it with initial data

2. **Run end-to-end tests:**

```shell script
export BASIC_AUTH_USERNAME=superadmin@walletera.dev
export BASIC_AUTH_PASSWORD=changeme
make run-tests
```

This will:
- Build the Go test runner
- Execute end-to-end tests against the running services

3. **Stop Services:**

```shell script
make stop
```
This stops and removes all containers, but keeps Docker images and volumes for faster restarts.

## Environment Configuration

Service configuration, credentials, secrets, and environment variables are specified via each service’s respective `compose.yml`. Ensure to review and customize them as needed—especially for credentials, secrets, and bind mounts—before deploying outside a development environment.

## Troubleshooting

- To rebuild a service (for development changes), bring it down and re-run `make start`.
- If you encounter issues with stale data or containers, use `docker-compose down -v` to remove volumes or `docker-compose rm` to clear stopped containers.
- Ports required by the stack (e.g., 8051, 3880, 3881) must be free on your host.

## Contributing

Please submit issues and pull requests for updates, bugfixes, or additional test scenarios.

---

**Note:** This environment is intended solely for local development and integration testing. Do not use default secrets or credentials in production.

---

For questions, contact the project maintainers or check the individual service READMEs for specific configuration and usage details.