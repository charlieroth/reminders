# Reminders

A HTTP API to store reminders & a CLI application to manage reminders.

## Project Structure

```
.
├── cmd/
│   ├── admin/
│   │   └── main.go
│   ├── cli/
│   │   └── main.go
│   └── server/
│       └── main.go
├── internal
├── zarf/
│   ├── compose
│   ├── docker
│   ├── k8s
│   └── kind-config.yaml
└── Makefile
```

The `cmd` directory contains the entry points for the CLI, HTTP API and admin applications.

The `internal` directory contains the shared logic for the project.

The `zarf` directory contains the infrastructure as code for the project.

The `Makefile` contains the commands for building, testing, and running the project.
