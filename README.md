# user-service

# Overview

This project uses the following technologies:

- [postgres](https://www.postgresql.org/) for the databases.
- [go](https://go.dev/) as the server backend language.

# Badges

[![codecov](https://codecov.io/gh/Knoblauchpilze/user-service/graph/badge.svg?token=94SJ008IPB)](https://codecov.io/gh/Knoblauchpilze/user-service)

[![Build services](https://github.com/Knoblauchpilze/user-service/actions/workflows/build-and-deploy.yml/badge.svg)](https://github.com/Knoblauchpilze/user-service/actions/workflows/build-and-deploy.yml)

[![Database migration tests](https://github.com/Knoblauchpilze/user-service/actions/workflows/database-migration-tests.yml/badge.svg)](https://github.com/Knoblauchpilze/user-service/actions/workflows/database-migration-tests.yml)

# Installation

The tools described below are directly used by the project. It is mandatory to install them in order to build the project locally.

See the following links:

- [golang](https://go.dev/doc/install): this project was developed using go `1.23.2`.
- [golang migrate](https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md): following the instructions there should be enough.
- [postgresql](https://www.postgresql.org/) which can be taken from the packages with `sudo apt-get install postgresql-14` for example.

We also assume that this repository is cloned locally and available to use. To achieve this, just use the following command:

```bash
git clone git@github.com:Knoblauchpilze/user-service.git`
```
