#!/usr/bin/env bash
# Seeds a demo vault for the vhs recording.
# Passphrase: "demo" (low work factor set via env var).
set -euo pipefail

VAULT=/tmp/demo-vault.age
BINARY=/tmp/envii-demo

rm -f "$VAULT"

# Use a tiny Go program to seed the vault programmatically.
go run docs/seed_vault.go "$VAULT"
