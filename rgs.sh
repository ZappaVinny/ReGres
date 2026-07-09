#!/usr/bin/env bash

set -e

cmd="$1"
shift || true

compose() {
  docker compose "$@"
}

show_help() {
  echo "ReGres CLI"
  echo ""
  echo "Usage:"
  echo "  ./rgs.sh <command> [args]"
  echo ""
  echo "Core commands:"
  echo "  ./rgs.sh init              Build images and create containers"
  echo "  ./rgs.sh run               Run the web, srv, and db stack"
  echo "  ./rgs.sh up                Run the stack in the background"
  echo "  ./rgs.sh down              Stop and remove containers"
  echo "  ./rgs.sh stop              Stop containers without removing them"
  echo "  ./rgs.sh restart [service] Restart all containers or one service"
  echo "  ./rgs.sh build             Build images"
  echo "  ./rgs.sh logs [service]    Show logs"
  echo ""
  echo "Service commands:"
  echo "  ./rgs.sh web <command>     Run a command in the web container"
  echo "  ./rgs.sh srv <command>     Run a command in the srv container"
  echo "  ./rgs.sh db <command>      Run a command in the db container"
  echo ""
  echo "Examples:"
  echo "  ./rgs.sh web npm install"
  echo "  ./rgs.sh web npm run build"
  echo "  ./rgs.sh srv go test ./..."
  echo "  ./rgs.sh srv go mod tidy"
  echo "  ./rgs.sh db psql -U postgres -d regres"
}

service_exec() {
  service="$1"
  shift

  if [ -z "$1" ]; then
    echo "No command provided for service '$service'."
    echo "Example: ./rgs.sh $service sh"
    exit 1
  fi

  compose exec "$service" "$@"
}

run_migrations_reset() {
  compose run --rm srv sh -c '
    migrate -path db/migrations -database "$DATABASE_URL" down -all || true
    migrate -path db/migrations -database "$DATABASE_URL" up
  '
}

case "$cmd" in
  ""|help|-h|--help)
    show_help
    ;;

  init)
    compose build

    echo "Initializing database+++"
    compose up -d db

    echo "Running migrations+++"
    compose run --rm srv sh -c '
        migrate -path db/migrations -database "$DATABASE_URL" down -all || true
        migrate -path db/migrations -database "$DATABASE_URL" up
    '

    echo "Constructing Models+++"
    compose run --rm srv sqlc generate
    
    compose create web srv
    echo ""
    echo "ReGres initialized."
    echo "Run the stack with:"
    echo "  ./rgs.sh run"
    ;;

  run)
    compose up web srv
    ;;

  up)
    compose up -d web srv
    ;;

  down)
    compose down
    ;;

  stop)
    compose stop
    ;;

  build)
    compose build
    ;;

  restart)
    if [ -z "$1" ]; then
      compose restart
    else
      compose restart "$1"
    fi
    ;;

  logs)
    compose logs -f "$@"
    ;;

    migrate)
    subcmd="$1"
    shift || true

    case "$subcmd" in
      up)
        echo "Running migrations up+++"
        compose up -d db
        compose run --rm srv sh -c '
          migrate -path db/migrations -database "$DATABASE_URL" up
        '
        ;;

      down)
        echo "Running migrations down+++"
        compose up -d db
        compose run --rm srv sh -c '
          migrate -path db/migrations -database "$DATABASE_URL" down -all
        '
        ;;

      reset)
        echo "Resetting database+++"
        run_migrations_reset

        echo "Constructing Models+++"
        compose run --rm srv sqlc generate
        ;;

      *)
        echo "Usage:"
        echo "  ./rgs.sh migrate up"
        echo "  ./rgs.sh migrate down"
        echo "  ./rgs.sh migrate reset"
        exit 1
        ;;
    esac
    ;;


  web|srv|db)
    service_exec "$cmd" "$@"
    ;;

  *)
    echo "Unknown command: $cmd"
    echo ""
    show_help
    exit 1
    ;;
esac