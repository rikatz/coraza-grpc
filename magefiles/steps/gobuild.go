package steps

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Go mg.Namespace

// Test runs the unittests
func (Go) Test() error {
	mg.Deps(Go{}.Tidy)
	fmt.Println("Running unit tests")
	return sh.RunV("go", "test", "-v", "./...")
}

// Build builds the program
func (Go) Build() {
	mg.Deps(Go{}.Server, Go{}.Client)
}

func (Go) Client() error {
	mg.Deps(Go{}.Tidy)
	fmt.Println("Building client")
	return sh.RunV("go", "build", "-o", "build/coraza-client", "cmd/client/main.go")
}

func (Go) Server() error {
	mg.Deps(Go{}.Tidy)
	fmt.Println("Building server")
	return sh.RunV("go", "build", "-o", "build/coraza-grpc", "cmd/server/main.go")
}

// Tidy tidies the go package
func (Go) Tidy() error {
	fmt.Println("Running go mod tidy")
	return sh.Run("go", "mod", "tidy")
}

// Vet runs a govet
func (Go) Vet() error {
	fmt.Println("Running go vet")
	return sh.Run("go", "vet", "./...")
}

// Fmt formats the code
func (Go) Fmt() error {
	fmt.Println("Running go fmt")
	return sh.Run("go", "fmt", "./...")
}
