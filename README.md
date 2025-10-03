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

**Clone Required Repositories:**

Before starting, you need to clone the following repositories into your local environment:

```shell
# Create a directory for Walletera repositories if it doesn't exist
mkdir -p /path/to/walletera

# Clone the required repositories
cd /path/to/walletera
git clone https://github.com/walletera/payments.git
git clone https://github.com/walletera/dinopay-gateway.git
git clone https://github.com/walletera/payments-read-model.git
git clone https://github.com/walletera/local-testing
```

**Configure Environment:**

`local-testing` is the main directory where you'll run all local development and testing commands. 
Make sure you're in this directory when executing the following setup steps.

```shell
cd /path/to/walletera/local-testing
```

Set up the `.env` file in the local-testing directory:

```shell
echo "WALLETERA_DIR=/path/to/walletera" > .env
```

**Start All Services:**

```shell
make start
```

This command will:
- Bring up all Docker Compose services
- Configure Barong (authentication service) database and seed it with initial data

**Create a Payment**

First, we need to authenticate and get a session cookie.

> [!NOTE]
> Basic auth is fine for a local environment, but in production use [api keys](https://github.com/walletera/barong/blob/2-6-stable/docs/api/barong_user_api_v2.md#apiv2barongresourceservice_accountsapi_keys) instead.

```shell
curl -v --location 'http://127.0.0.1:3099/api/v1/auth/identity/sessions' \
--form 'email="superadmin@walletera.dev"' \
--form 'password="aBB^kBg4"'
```
Search for the set-cookie header in the response. It should look something like this:

```shell
 set-cookie: _walletera=2c6242401119f8857b2280e62908369d; path=/; expires=Sat, 04 Oct 2025 14:56:35 -0000; HttpOnly
```

Grab the first part `_walletera=2c6242401119f8857b2280e62908369d`.

Create a payment using the following script.
Note that you need to put the cookie you grabbed from the previous step in the `Cookie` header.

```shell
curl --location 'http://127.0.0.1:3099/api/v1/payments' \
--header 'Content-Type: application/json' \
--header 'Cookie: _walletera=2c6242401119f8857b2280e62908369d' \
--data '{
    "id": "8de1968e-a288-4ba4-9612-2cca1bfaa23a",
    "amount": 100,
    "currency": "USD",
    "gateway": "dinopay",
    "debtor": {
        "institutionName": "dinopay",
        "institutionId": "dinopay",
        "currency": "USD",
        "accountDetails": {
            "accountType": "dinopay",
            "accountHolder": "Richard Roe",
            "accountNumber": "1200079635"
        }
    },
    "beneficiary": {
        "institutionName": "dinopay",
        "institutionId": "dinopay",
        "currency": "USD",
        "accountDetails": {
            "accountType": "dinopay",
            "accountHolder": "Richard Roe",
            "accountNumber": "1200079635"
        }
    }
}'
```

You should receive as a response a JSON similar to the one below:

```json
{
  "id": "8de1968e-a288-4ba4-9612-2cca1bfaa23a",
  "amount": 100,
  "currency": "USD",
  "gateway": "dinopay",
  "debtor": {
    "institutionName": "dinopay",
    "institutionId": "dinopay",
    "currency": "USD",
    "accountDetails": {
      "accountType": "dinopay",
      "accountHolder": "Richard Roe",
      "accountNumber": "1200079635"
    }
  },
  "beneficiary": {
    "institutionName": "dinopay",
    "institutionId": "dinopay",
    "currency": "USD",
    "accountDetails": {
      "accountType": "dinopay",
      "accountHolder": "Richard Roe",
      "accountNumber": "1200079635"
    }
  },
  "direction": "outbound",
  "customerId": "00000000-0000-0000-0000-000000000000",
  "status": "pending",
  "createdAt": "0001-01-01T00:00:00Z",
  "updatedAt": "0001-01-01T00:00:00Z"
}
```

Grab the payment id from the `id` field and use it to retrieve the payment (from the payments-read-model).

```shell
curl --location 'http://127.0.0.1:3099/api/v1/payments/8de1968e-a288-4ba4-9612-2cca1bfaa23a' \
--header 'Cookie: _walletera=2c6242401119f8857b2280e62908369d'
```

Note that now the payment status is `confirmed`.

To understand how walletera processed this payment, I suggest taking a look at the logs
of the components involved in the processing.

```shell
docker-compose logs payments dinopay-gateway payments-read-model
```

**Run end-to-end tests:**

```shell
export BASIC_AUTH_USERNAME=superadmin@walletera.dev
export BASIC_AUTH_PASSWORD=aBB^kBg4
make run-tests
```

This will:
- Build the Go test runner
- Execute end-to-end tests against the running services

**Stop Services:**

```shell script
make stop
```
This stops and removes all containers but keeps Docker images and volumes for faster restarts.

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
