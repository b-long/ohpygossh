# Development tasks for ohpygossh.
# Install `just`: https://github.com/casey/just
#
# Run `just` (or `just --list`) to see all available recipes.

set shell := ["bash", "-euo", "pipefail", "-c"]

# gopy's version is pinned in go.mod (as a Go 1.24+ tool dependency) and
# managed there by Dependabot; every consumer derives it from go.mod so
# there is exactly one place to bump.
gopy_version := `go list -m -f '{{.Version}}' github.com/go-python/gopy`

# Show available recipes
default:
    @just --list

# Install all tooling needed for local development (macOS and Linux, via Homebrew)
setup:
    #!/usr/bin/env bash
    set -euo pipefail
    if ! command -v brew >/dev/null 2>&1; then
        echo "error: Homebrew is required. Install it from https://brew.sh, then re-run 'just setup'." >&2
        exit 1
    fi

    echo "==> Installing Go, Python 3.11, poetry, pre-commit, golangci-lint"
    brew install go python@3.11 poetry pre-commit golangci-lint

    echo "==> Installing gopy and goimports"
    go install "github.com/go-python/gopy@{{ gopy_version }}"
    go install golang.org/x/tools/cmd/goimports@latest

    GOBIN="$(go env GOPATH)/bin"
    case ":$PATH:" in
        *":$GOBIN:"*) ;;
        *) echo "warning: $GOBIN is not on your PATH; add it to use gopy/goimports directly" >&2 ;;
    esac

    echo "==> Configuring the Python virtual environment"
    poetry env use "$(brew --prefix python@3.11)/bin/python3.11"
    poetry install --no-root

    echo "==> Installing git hooks"
    pre-commit install

    echo "Setup complete."

# Run Go and Python linters (same checks as CI)
lint:
    go get .
    pre-commit run --all-files

# Run the Go test suite
test:
    go test ./...

# Build and validate the Python wheel (produces dist/ohpygossh-*.whl)
#
# CGO_ENABLED=1 is required for gopy's c-shared build and isn't always the
# default (e.g. on some Linux/arm64 toolchains); see build-golang-macos.yaml.
build:
    CGO_ENABLED=1 ./make_and_validate_script.sh

# Re-validate an already-built wheel without rebuilding it
validate:
    ./only_validate.sh

# Remove build artifacts and the virtual environment
clean:
    rm -rf .venv dist build __pycache__ ohpygossh.egg-info ohpygossh myssh

# Run everything CI runs: setup, lint, then test
ci: setup lint test

# Validate that gopy's pin is only sourced from go.mod (no stray hardcoded copies)
gopy-version-check:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "gopy is pinned to: {{ gopy_version }}"

    if [[ "{{ gopy_version }}" =~ ^v[0-9]+\.[0-9]+\.[0-9]+-(0\.)?[0-9]{14}-[0-9a-f]{12}(\+incompatible)?$ ]]; then
        echo "note: gopy is still pinned to an untagged commit ({{ gopy_version }})."
        echo "      Check https://github.com/go-python/gopy/releases -- if a tagged"
        echo "      release now exists, run 'just gopy-version-bump' to switch."
    fi

    if grep -RInE 'gopy@[0-9a-f]{40}' --include='*.yaml' --include='*.yml' --include='*.sh' .github build-scripts Justfile; then
        echo "error: found a hardcoded gopy commit pin outside go.mod; derive it from go.mod instead" >&2
        exit 1
    fi
    echo "no stray hardcoded gopy pins found"

# Bump gopy to the tip of its default branch, updating go.mod/go.sum (the single source of truth).
# Uses @master rather than @latest because gopy's tags lag its default branch;
# @latest would silently downgrade past commits this project relies on.
gopy-version-bump:
    go get -tool github.com/go-python/gopy@master
    go mod tidy
