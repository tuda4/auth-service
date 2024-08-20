# mb-backend project

## Tech stack

- Go-lang
- Echo
- sqlc

## Getting started

```sh
// Create .env 
cp .env.example app.env
```

### Run app on docker

```sh
// Docker up
docker compose up -d
```

### Run app on terminal

Install air

```sh
go install github.com/air-verse/air@latest
```

Run app via air

```sh
air -c .air.toml
```

Access to api at http://localhost:8080

