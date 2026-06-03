# Contributing to envii

Thanks for your interest! Contributions of all kinds are welcome.

## Getting started

```sh
git clone https://github.com/trickylab/envii
cd envii
go build ./...
go test ./...
```

## Guidelines

- Run `go test ./...` and `go vet ./...` before opening a PR.
- Keep changes focused; one logical change per PR.
- Follow standard Go formatting (`gofmt` / `goimports`).
- For features, open an issue first to discuss the design.

## Reporting bugs

Open an issue with:
- What you expected vs. what happened
- Steps to reproduce
- OS, terminal, and `envii --version`

## Security

Do **not** open public issues for security vulnerabilities. Email the
maintainer instead (see repository profile).
