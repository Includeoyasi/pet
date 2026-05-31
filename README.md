# pet

pet is a small REST API for memory capsules.

The current implementation exposes HTTP endpoints for creating, reading, and deleting capsules. Data is stored in memory for now, which keeps the service easy to run while the API shape is being developed.

## Run

```sh
go run .
```

The service reads `internal/config/config.yml` and listens on `HTTP_PORT`, which is `8080` by default.

## API

Create a capsule:

```sh
curl -X POST http://localhost:8080/api/create \
  -H 'Content-Type: application/json' \
  -d '{"author":"Includeoyasi","text":"First memory capsule","label":"demo"}'
```

Get a capsule:

```sh
curl http://localhost:8080/api/capsule?id=1
```

Delete a capsule:

```sh
curl -X DELETE http://localhost:8080/api/delete?id=1
```

Auth endpoints are present as placeholders:

```sh
POST /auth/sing-up
POST /auth/sing-in
```

## Development

```sh
go test ./...
```
