# Weekend API

This contains short documentation about how to use this project/repo (an API built over the weekend)

## Development
- assuming our workstation already completed setup for development in Go with configured database

- install dependencies and tools for this API
    `go mod download`

- install swaggo (swagger for go)
    `go install github.com/swaggo/swag/cmd/swag@latest`

- auto-generate API documentation with swagger
    `swag init`

- copy the env file and fill the necessary environment variables and database config
    `cp .env.example .env`

- run database migration to setup database schema and create user admin
   `CREATE_ADMIN_USER=true go run db/migration.go`

- run the project
    `go run main.go`

- to access swagger documentation, go to  `{APP_URL}/swagger/index.html`

- default credential for Administrator role is
    `Username: admin`
    `Password: admin`
