# Econominhas - Authentication

Authentication microservice, responsible for controlling users account creation and sign ins.

## About

This project use lot's of tools to be as efficient as possible, here's the list with the links that you need to learn more about them.

### Help Development

- [make](https://makefiletutorial.com/) to run commands easily
- [localstack](https://www.localstack.cloud/) to simulate AWS environment locally
- [docker](https://www.docker.com/) & [docker-compose](https://docs.docker.com/compose/) to orchestrate (to "run") the api, database, localstack and all the heavy-external tools that we need to make the project work
- [github actions](https://docs.github.com/en/actions) to run pipelines to deploy and validate things
- [golang migrate](https://github.com/golang-migrate/migrate) to manage migrations
- [editorconfig](https://editorconfig.org/) to help with linting

### Documentation

- [dbdocs](https://dbdocs.io/) to host the database docs written in DBML
- [openapi](https://www.openapis.org/) to document the API routes
  - We don't document the API using the code to don't bind us to any library or framework, this way we can be more tool agnostic and use the default way to document APIs: OpenAPI

### Hosted Docs

- [API](https://econominhas.readme.io/reference/)
- [Database](https://dbdocs.io/henriqueleite42/Econominhas?view=relationships)

## How to

### Run the API for the first time

1. Copy and paste `.env.example`, rename the copy to `.env.docker` and use the values provided to you by your team
2. Run `make start`
3. Open another console tab and run `make migrate`
4. The API will be available at http://localhost:3000/

### Create a migration

```sh
make gen-migration NAME=CreateAccountsTable
```
