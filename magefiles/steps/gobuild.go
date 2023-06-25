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
func (Go) Build() error {
	mg.Deps(Go{}.Tidy)
	fmt.Println("Building command")
	return sh.RunV("go", "build", "-o", "build/coraza-grpc", "cmd/main.go")
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
