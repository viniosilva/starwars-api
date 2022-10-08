# Star Wars API

## Requirements

- [go](https://tip.golang.org/doc/go1.19)
- [mockgen](https://github.com/golang/mock)
- [golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
- [swaggo](https://github.com/swaggo/swag)

## Instalation

```bash
$ go get
```

## Configuration

### Creating .env file

For local environment just create a `.env` like this:

```txt
GIN_MODE=debug
MIGRATE_URL=mysql://luke:xQlpKD95kp20Wa1JAX6O@\(127.0.0.1:3306\)/starwars
MYSQL_PASSWORD=xQlpKD95kp20Wa1JAX6O
```

### Migrate

After running `docker-compose`, it's necessary to wait a few seconds to run the `migrate` command.

```bash
$ docker-compose up -d
$ make migrate
```

## Running

```bash
$ make run
```

See API local documentation at [swagger](http:localhost:8080/api/swagger/index.html)

## Tests

```bash
# Unit tests
$ make test

# Tests coverage
$ make test/cov

# Tests benchmark service module
$ make test/bench
```
