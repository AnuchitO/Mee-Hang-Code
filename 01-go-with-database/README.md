# Refactoring Go with Database

## Instuctions
1. copy `.env.example` to `.env` and fill in the values
1. if you don't have [direnv](https://direnv.net/) installed then `source .env` to load the environment variables;
1. `make start-db` to start the database
1. `make run` to run the server
1. `make get-skill` to test get skill `go` from the server

## e2e
1. `make e2e` to run the e2e tests if returned `true` then the tests passed. this is a simple e2e to check just the get skill by key endpoint

## Endpoints
- `GET /skills/:key` - get a skill by key
- `GET /skills` - get all skills
- `POST /skills` - create a new skill
