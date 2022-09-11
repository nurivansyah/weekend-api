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

## Deployment

- assuming our kubernetes cluster already completed setup along with the database, so the focus on this repo is to deploy this container into the kubernetes

- build image, this process might be triggered at CI stage, in case we need to build manually just run below command
    `make docker-build`

- you can test and run the built container using provided `docker-compose.yml`
    `docker compose --env-file ./.env up`

- jump into the cluster and prepare deploying our api

- the cluster utilize ingress to handle traffic from outside cluster, so we might also need to setup the dns to translate domain into the cluster ip

- create namespace for the app
    `kubectl create ns weekend-api`

- fill the necessary environment variables in `configmap.yml` and `secret.yml`, then apply manifests
    `kubectl -n weekend-api apply -f kubernetes/`

- or, if we decide to deploy using helm chart we can use the provided chart in this repo, of course we need to fill the values with necessary environment variables too
    `cd helm`
    `helm upgrade -i api weekend-api -n weekend-api`

- after the manifest applied or helm chart installed, check all resources status
    `kubectl -n weekend-api get all`

- access the `APP_URL` and to access swagger documentation, go to  `{APP_URL}/swagger/index.html`

- note: there are many ways to deploy the app into kubernetes cluster, either you can push the manifest into the cluster, or we can also setup CD tools like ArgoCD and watch to this repo changes to automate deploy the cluster
