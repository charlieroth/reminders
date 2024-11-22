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
KIND            := kindest/node:v1.31.2
POSTGRES        := postgres:16.4

KIND_CLUSTER    := reminders-cluster
NAMESPACE       := reminders-system
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
	brew list kind || brew install kind
	brew list kubectl || brew install kubectl
	brew list kustomize || brew install kustomize
	brew list pgcli || brew install pgcli
	brew list watch || brew install watch

dev-docker:
	docker pull $(GOLANG) & \
	docker pull $(ALPINE) & \
	docker pull $(KIND) & \
	docker pull $(POSTGRES) & \
	wait;

# ==============================================================================
# Building containers

build: reminders

reminders:
	docker build \
		-f zarf/docker/dockerfile.reminders \
		-t $(REMINDERS_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.

# ==============================================================================
# Running from within k8s/kind

dev-up:
	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/dev/kind-config.yaml

	kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner

	kind load docker-image $(POSTGRES) --name $(KIND_CLUSTER)

dev-down:
	kind delete cluster --name $(KIND_CLUSTER)

dev-status-all:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

dev-status:
	watch -n 2 kubectl get pods -o wide --all-namespaces

# ------------------------------------------------------------------------------

dev-load:
	kind load docker-image $(REMINDERS_IMAGE) --name $(KIND_CLUSTER)

dev-apply:
	kustomize build zarf/k8s/dev/database | kubectl apply -f -
	kubectl rollout status --namespace=$(NAMESPACE) --watch --timeout=120s sts/database

	kustomize build zarf/k8s/dev/reminders | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=$(AUTH_APP) --timeout=120s --for=condition=Ready

dev-restart:
	kubectl rollout restart deployment $(REMINDERS_APP) --namespace=$(NAMESPACE)

dev-update: build dev-load dev-restart

dev-update-apply: build dev-load dev-apply

dev-logs:
	kubectl logs --namespace=$(NAMESPACE) -l app=$(REMINDERS_APP) --all-containers=true -f --tail=100 --max-log-requests=6

# ------------------------------------------------------------------------------

dev-describe-node:
	kubectl describe node

dev-describe-deployment:
	kubectl describe deployment --namespace=$(NAMESPACE) $(REMINDERS_APP)

dev-describe-reminders:
	kubectl describe pod --namespace=$(NAMESPACE) -l app=$(REMINDERS_APP)

dev-describe-database:
	kubectl describe pod --namespace=$(NAMESPACE) -l app=database

# ------------------------------------------------------------------------------

dev-logs-db:
	kubectl logs --namespace=$(NAMESPACE) -l app=database --all-containers=true -f --tail=100

# ------------------------------------------------------------------------------

dev-services-delete:
	kustomize build zarf/k8s/dev/reminders | kubectl delete -f -
	kustomize build zarf/k8s/dev/database | kubectl delete -f -

dev-describe-replicaset:
	kubectl get rs
	kubectl describe rs --namespace=$(NAMESPACE) -l app=$(REMINDERS_APP)

dev-events:
	kubectl get ev --sort-by metadata.creationTimestamp

dev-events-warn:
	kubectl get ev --field-selector type=Warning --sort-by metadata.creationTimestamp

dev-shell:
	kubectl exec --namespace=$(NAMESPACE) -it $(shell kubectl get pods --namespace=$(NAMESPACE) | grep reminders | cut -c1-26) --container reminders -- /bin/sh

dev-database-restart:
	kubectl rollout restart statefulset database --namespace=$(NAMESPACE)

# ==============================================================================
# Docker Compose

compose-up:
	cd ./zarf/compose/ && docker compose -f docker_compose.yaml -p compose up -d

compose-build-up: build compose-up

compose-down:
	cd ./zarf/compose/ && docker compose -f docker_compose.yaml down

compose-logs:
	cd ./zarf/compose/ && docker compose -f docker_compose.yaml logs

# ==============================================================================
# Administration

migrate:
	export REMINDERS_DB_HOST=localhost; go run cmd/admin/main.go migrate

seed: migrate
	export REMINDERS_DB_HOST=localhost; go run cmd/admin/main.go seed

pgcli:
	pgcli postgresql://postgres:postgres@localhost

liveness:
	curl -il http://localhost:8080/liveness

readiness:
	curl -il http://localhost:8080/readiness

# ==============================================================================
# Running tests within the local computer

test-down:
	docker stop servicetest
	docker rm servicetest -v

test-r:
	CGO_ENABLED=1 go test -race -count=1 ./...

test-only:
	CGO_ENABLED=0 go test -count=1 ./...

lint:
	CGO_ENABLED=0 go vet ./...
	staticcheck -checks=all ./...

vuln-check:
	govulncheck ./...

test: test-only lint vuln-check

test-race: test-r lint vuln-check

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
	@echo "  dev-up                  Start the KIND cluster"
	@echo "  dev-down                Stop the KIND cluster"
	@echo "  dev-status-all          Show the status of the KIND cluster"
	@echo "  dev-status              Show the status of the pods"
	@echo "  dev-load                Load the containers into KIND"
	@echo "  dev-apply               Apply the manifests to KIND"
	@echo "  dev-restart             Restart the deployments"
	@echo "  dev-update              Build, load, and restart the deployments"
	@echo "  dev-update-apply        Build, load, and apply the deployments"
	@echo "  dev-logs                Show the logs for the sales service"
	@echo "  dev-logs-auth           Show the logs for the auth service"
	@echo "  dev-logs-init           Show the logs for the init container"
	@echo "  dev-describe-node       Show the node details"
	@echo "  dev-describe-deployment Show the deployment details"
	@echo "  dev-describe-reminders  Show the reminders pod details"
	@echo "  dev-describe-database   Show the database pod details"
	@echo "  dev-logs-db             Show the logs for the database service"
	@echo "  dev-services-delete     Delete all"