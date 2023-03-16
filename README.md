# artist-db

## Local development

Start the whole stack (attached):

```shell
make start
```

Start the database (attached):

```shell
make start-db
```

Start the  API (attached):

```shell
make start-api
```

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

## Configuration 

### Frontend Configuration

You need a `config.json` at `./frontend/src/assets/data`. You'll find an example file at this location holding the variables that need to be set.