# user-service

# Overview

This project uses the following technologies:

- [postgres](https://www.postgresql.org/) for the databases.
- [go](https://go.dev/) as the server backend language.

# Badges

[![codecov](https://codecov.io/gh/Knoblauchpilze/user-service/graph/badge.svg?token=94SJ008IPB)](https://codecov.io/gh/Knoblauchpilze/user-service)

[![Build services](https://github.com/Knoblauchpilze/user-service/actions/workflows/build-and-push.yml/badge.svg)](https://github.com/Knoblauchpilze/user-service/actions/workflows/build-and-push.yml)

[![Database migration tests](https://github.com/Knoblauchpilze/user-service/actions/workflows/database-migration-tests.yml/badge.svg)](https://github.com/Knoblauchpilze/user-service/actions/workflows/database-migration-tests.yml)

# Installation

The tools described below are directly used by the project. It is mandatory to install them in order to build the project locally.

See the following links:

- [golang](https://go.dev/doc/install): this project was developed using go `1.23.2`.
- [golang migrate](https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md): following the instructions there should be enough.
- [postgresql](https://www.postgresql.org/) which can be taken from the packages with `sudo apt-get install postgresql-14` for example.

We also assume that this repository is cloned locally and available to use. To achieve this, just use the following command:

```bash
git clone git@github.com:Knoblauchpilze/user-service.git
```

## Secrets in the CI

The CI workflows define several secrets that are expected to be created for the repository when cloned/forked/used. Each secret should be self-explanatory based on its name. Most of them require to setup an account on one or the other service mentioned in this README.

# How does sign-up/login/logout work?

## General design

The base of our authentication system is the `user-service`. This service describes what a user is and which ones are registered in our system.

Each user is attached a set of credentials, along with some permissions and limits. This information can be returned by the `user-service` through a `auth` endpoint.

Traefik has a [forwardAuth](https://doc.traefik.io/traefik/middlewares/http/forwardauth) middleware which allows (as its name suggests) to forward any request it receives to an external authentication server. Based on the response of this server it either denies or forwars the request.

We leveraged this principle to hook the `user-service` with this middleware so that we can control the access to the API to only authenticated users.

## The session concept

In order to allow users to access information about the service, we provide a session mechanism. It is quite a wide topic and we gathered a few resources on whether this is a RESTful approach or not in the dedicated [PR #7](https://github.com/KnoblauchPilze/galactic-sovereign/pull/7).

Upon calling the `POST /v1/users/sessions` route, the user will be able to obtain a token valid (see [API keys](#api-keys)) for a certain period of time and which can be used to access other endpoints in the cluster. This endpoint, along with the `POST /v1/users` endpoint to create a new user, are the only one which can be called unauthenticated.

The session token is only valid for a certain amount of time and can be revoked early by calling `DELETE /v1/users/sessions/{user-id}`.

## API keys

We use API keys in a similar way as the session keys described in this [Kong article](https://konghq.com/blog/learning-center/what-are-api-keys). Each key is a simple identifier that is required to access our service. It is created upon logging in and deactivated upon logging out.

## The authentication endpoint

The authentication endpoint is a corner stone of the strategy: this takes any http request and look for an API key attached to it as a header:

- if there's no such header the request is denied.
- if there's one but the key is invalid (either expired or unknown) the request is denied.

As all requests are routed towards this endpoint by traefik before they reach the target service, we can guarantees an efficient filtering and only allow authorized users to access our cluster.

# Cheat sheet

## Create new user

```bash
curl -X POST -H "Content-Type: application/json" http://localhost:60001/v1/users -d '{"email":"user-1@mail.com","password":"password-for-user-1"}' | jq
```

## Query existing user

```bash
curl -X GET -H "Content-Type: application/json" -H 'X-Api-Key: 2da3e9ec-7299-473a-be0f-d722d870f51a' http://localhost:60001/v1/users/4f26321f-d0ea-46a3-83dd-6aa1c6053aaf | jq
```

## Query non existing user

```bash
curl -X GET -H 'Content-Type: application/json' -H 'X-Api-Key: 2da3e9ec-7299-473a-be0f-d722d870f51a' http://localhost:60001/v1/users/4f26321f-d0ea-46a3-83dd-6aa1c6053aae | jq
```

## Query without API key

```bash
curl -X GET -H 'Content-Type: application/json' http://localhost:60001/v1/users/4f26321f-d0ea-46a3-83dd-6aa1c6053aae | jq
```

## List users

```bash
curl -X GET -H 'Content-Type: application/json' '-H 'X-Api-Key: 2da3e9ec-7299-473a-be0f-d722d870f51a' http://localhost:60001/v1/users | jq
```

## Patch existing user

```bash
curl -X PATCH -H 'Content-Type: application/json' -H 'X-Api-Key: 2da3e9ec-7299-473a-be0f-d722d870f51a' http://localhost:60001/v1/users/0463ed3d-bfc9-4c10-b6ee-c223bbca0fab -d '{"email":"test-user@real-provider.com","password":"strong-password"}'| jq
```

## Delete user

```bash
curl -X DELETE -H 'Content-Type: application/json' -H 'X-Api-Key: 2da3e9ec-7299-473a-be0f-d722d870f51a' http://localhost:60001/v1/users/0463ed3d-bfc9-4c10-b6ee-c223bbca0fab | jq
```

## Login a user

```bash
curl -X POST -H 'Content-Type: application/json' http://localhost:60001/v1/users/sessions/4f26321f-d0ea-46a3-83dd-6aa1c6053aaf | jq
```

## Login a user by email

```bash
curl -X POST -H "Content-Type: application/json" http://localhost:60001/v1/users/sessions -d '{"email":"test-user@provider.com","password":"strong-password"}' | jq
```

## Login a user by email with wrong credentials

```bash
curl -X POST -H "Content-Type: application/json" http://localhost:60001/v1/users/sessions -d '{"email":"test-user@provider.com","password":"not-the-password"}' | jq
```

## Logout a user

```bash
curl -X DELETE -H 'Content-Type: application/json' -H 'X-Api-Key: 2da3e9ec-7299-473a-be0f-d722d870f51a' http://localhost:60001/v1/users/sessions/4f26321f-d0ea-46a3-83dd-6aa1c6053aaf | jq
```
