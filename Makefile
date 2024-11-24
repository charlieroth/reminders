# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

# Deploy First Mentality

# ==============================================================================
# Go Installation
#
#	You need to have Go version 1.23 to run this code.
#
#	https://go.dev/dl/
#
#	If you are not allowed to update your Go frontend, you can install
#	and use a 1.23 frontend.
#
#	$ go install golang.org/dl/go1.23@latest
#	$ go1.23 download
#
#	This means you need to use `go1.23` instead of `go` for any command
#	using the Go frontend tooling from the makefile.

# ==============================================================================
# Brew Installation
#
# 	Install GCC:
#	$ brew install gcc

# ==============================================================================
# Install Tooling and Dependencies
#
#	This project uses Docker and it is expected to be installed. Please provide
#	Docker at least 4 CPUs. To use Podman instead please alias Docker CLI to
#	Podman CLI or symlink the Docker socket to the Podman socket. More
#	information on migrating from Docker to Podman can be found at
#	https://podman-desktop.io/docs/migrating-from-docker.
#
#	Run these commands to install everything needed.
#	$ make dev-brew
#	$ make dev-docker
#	$ make dev-gotooling

# ==============================================================================
# Running Test
#
#	Running the tests is a good way to verify you have installed most of the
#	dependencies properly.
#
#	$ make test

# ==============================================================================
# Running The Project
#
#	$ make dev-up
#	$ make dev-update-apply
#	$ make token
#	$ export TOKEN=<token>
#	$ make users
#
#	You can use `make dev-status` to look at the status of your KIND cluster.

# ==============================================================================
# Define dependencies

GOLANG          := golang:1.23
ALPINE          := alpine:3.20
POSTGRES        := postgres:16.4
GENI            := ghcr.io/emilpriver/geni:v1.1.4

REMINDERS_APP   := reminders
BASE_IMAGE_NAME := localhost/charlieroth
VERSION         := 0.0.1
REMINDERS_IMAGE := $(BASE_IMAGE_NAME)/$(REMINDERS_APP):$(VERSION)

# ==============================================================================
# Install dependencies

dev-gotooling:
	go install github.com/divan/expvarmon@latest
	go install github.com/rakyll/hey@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install golang.org/x/tools/cmd/goimports@latest

dev-brew:
	brew update
	brew list pgcli || brew install pgcli
	brew list watch || brew install watch
	brew list geni || brew install geni

dev-docker:
	docker pull $(GOLANG) & \
	docker pull $(ALPINE) & \
	docker pull $(POSTGRES) & \
	docker pull $(GENI) & \
	wait;

# ==============================================================================
# Docker Compose

compose-dev-up:
	docker compose --profile dev up

compose-db-up:
	docker compose --profile db up

compose-dev-down:
	docker compose --profile dev down

compose-db-down:
	docker compose --profile db down

# ==============================================================================
# Administration

pgcli:
	pgcli postgresql://postgres:postgres@localhost:5432/reminders_dev

liveness:
	curl -il http://localhost:8080/liveness

readiness:
	curl -il http://localhost:8080/readiness

# ==============================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-list:
	go list -m -u -mod=readonly all

deps-upgrade:
	go get -u -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

list:
	go list -mod=mod all

# ==============================================================================
# Class Stuff

run:
	go run cmd/server/main.go

run-help:
	go run cmd/server/main.go --help

admin:
	go run cmd/admin/main.go

ready:
	curl -il http://localhost:8080/readiness

live:
	curl -il http://localhost:8080/liveness

# ==============================================================================
# Help command
help:
	@echo "Usage: make <command>"
	@echo ""
	@echo "Commands:"
	@echo "  dev-gotooling           Install Go tooling"
	@echo "  dev-brew                Install brew dependencies"
	@echo "  dev-docker              Pull Docker images"
	@echo "  build                   Build all the containers"
	@echo "  reminders               Build the reminders container"
	@echo "  compose-dev-up          Start the Docker Compose cluster in dev mode"
	@echo "  compose-db-up           Start the Docker Compose cluster in db mode"
	@echo "  compose-dev-down        Stop the Docker Compose cluster in dev mode"
	@echo "  compose-db-down         Stop the Docker Compose cluster in db mode"
	@echo "  run                     Run the Reminders server"
	@echo "  ready                   Make GET:/readiness request"
	@echo "  live                    Make GET:/livness request"
	@echo "  pgcli                   Connect to the database"
