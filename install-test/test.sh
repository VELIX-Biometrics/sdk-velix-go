#!/usr/bin/env bash
# Simula um consumidor externo com `go get`: usa `go mod edit -replace`
# apontando pro checkout local (mesmo efeito de resolução de módulo que
# `go get github.com/VELIX-Biometrics/sdk-velix-go` teria, sem precisar
# de rede/tag publicada) — não usa o pacote via caminho relativo direto.
set -e
REPO_DIR="$(cd "$(dirname "$0")/.." && pwd)"

rm -rf /tmp/velix-install-test-go
mkdir -p /tmp/velix-install-test-go
cd /tmp/velix-install-test-go

go mod init velix-install-test >/dev/null
cat > main.go <<'EOF'
package main

import (
	"fmt"
	velix "github.com/VELIX-Biometrics/sdk-velix-go"
)

func main() {
	client := velix.NewClient(velix.Config{APIURL: "http://localhost", APIKey: "test"})
	_ = client
	fmt.Println("INSTALL_TEST:go:PASS: módulo resolvido via go.mod replace, client construído")
}
EOF

go mod edit -require=github.com/VELIX-Biometrics/sdk-velix-go@v0.0.0
go mod edit -replace github.com/VELIX-Biometrics/sdk-velix-go="$REPO_DIR"
go mod tidy
go run main.go
