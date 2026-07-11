# ReGres

Lightweight Semi-Opinionated Containerized Go + React + PostgreSQL Web Development Stack Kit

## The Goal

ReGres is a Go + React starter kit built to skip the boilerplate that comes with every new full-stack project. Auth, project structure, basic API setup, it's already there. Developers can go straight to building features instead of doing the same setup work over and over.

It's not locked into one structure either. Swap parts out, restructure things, make it fit the project.

For local development everything runs in a single Docker Compose setup, frontend, backend, and database all built and running together. Later on there will be a build script that splits this apart, taking each piece of the stack and packaging it into its own standalone container that can scale on its own. Development stays in one place. Deployment doesn't have to.

## Default Technologies

- Go
  - sqlc (CLI)
  - migrate (CLI)
  - pgxpool
  - godotenv
  - crypto
- PostgreSQL
- npm
  - Typescript
  - React
  - Vite
  - Tailwind

## Features

- Golang backend with session based auth, middleware, migrations, models, and configs ready to be expanded.
- Only uses built in Go http package, using pgxpool, sqlc, and migrate to a enable basic database access.
- Very basic React + Typescript setup ready to be customized to fit an application's need.
- Very basic CLI to control development environment (and later build tooling).

## Local Setup

1. Clone repo
2. Copy .env.example to .env
3. Run `./rgs.sh init` to build containers and run migrations
4. Run `./rgs.sh run` to start local frontend and backend server
5. Develop your application

## Roadmap

- Transition CLI tool from Shell Script to a Go binary
- Add build/detachment feature
- Add default test support for backend (and maybe frontend)
- Implement any other ideas

## Interested in Contributing?

Check out [CONTRIBUTING.md](CONTRIBUTING.md) for details on how to get started

## License

[GNU GPLv3](LICENSE)
