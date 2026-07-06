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
  echo "  rgs <command> [args]"
  echo ""
  echo "Core commands:"
  echo "  rgs init              Build images and create containers"
  echo "  rgs run               Run the web, srv, and db stack"
  echo "  rgs up                Run the stack in the background"
  echo "  rgs down              Stop and remove containers"
  echo "  rgs stop              Stop containers without removing them"
  echo "  rgs restart [service] Restart all containers or one service"
  echo "  rgs build             Build images"
  echo "  rgs logs [service]    Show logs"
  echo ""
  echo "Service commands:"
  echo "  rgs web <command>     Run a command in the web container"
  echo "  rgs srv <command>     Run a command in the srv container"
  echo "  rgs db <command>      Run a command in the db container"
  echo ""
  echo "Examples:"
  echo "  rgs web npm install"
  echo "  rgs web npm run build"
  echo "  rgs srv go test ./..."
  echo "  rgs srv go mod tidy"
  echo "  rgs db psql -U postgres -d regres"
}

service_exec() {
  service="$1"
  shift

  if [ -z "$1" ]; then
    echo "No command provided for service '$service'."
    echo "Example: rgs $service sh"
    exit 1
  fi

  compose exec "$service" "$@"
}

case "$cmd" in
  ""|help|-h|--help)
    show_help
    ;;

  init)
    compose build
    compose create
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