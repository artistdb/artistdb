# artist-db

artist-db is an open source database for artists and their works. It is 
currently a work in progress.

## Setup

### Configuration

You need to configure `config.json` at `./frontend/src/assets/data`. With the
following structure:

```json
{
  "apiUri": "http://localhost:8080"
}
```

The URI should point to the API server, which for local development is localhost.
You will need to change this for production.

### Local development

Make sure you have the following prerequisites installed:

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Node.js](https://nodejs.org/en/) >= 16.0.0
- [Go](https://golang.org/) >= 1.17

Start the whole stack (attached):

```shell
make start
```

Start the database and API (attached):

```shell
make start
```

Start the frontend (attached):

```shell
make start-frontend
```

The frontend will be available under `http://localhost:4200`.

Alternatively run `make start-full` to bring up every service, including
monitoring and tracing.

Stop the stack:

```shell
make stop
```

With a database up-and-running integration tests can be run. This will require
another terminal window, or detachment from the DB container:

```shell
make test-integration
```

Alternatively you can run the whole test suite locally with:

```shell
make test-local
```

