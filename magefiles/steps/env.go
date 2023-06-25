package steps

import (
	"fmt"
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var lookupEnv = func(envVar, fallback string) string {
	value, ok := os.LookupEnv(envVar)
	if !ok {
		return fallback
	}
	return value
}

var (
	ControllerGenVersion = lookupEnv("PROTOC_VERSION", "v0.12.0")

	GolangLint        = lookupEnv("GOLANGCI_LINT", getGoBin()+"/golangci-lint")
	GolangLintPath    = "github.com/golangci/golangci-lint/cmd/golangci-lint"
	GolangLintOpts    = lookupEnv("GOLANGCI_LINT_OPTS", "-v --out-format colored-line-number:stdout,junit-xml:report-lint.xml")
	GolangLintVersion = lookupEnv("GOLANGCI_LINT_VERSION", "v1.53.3")

	ProtocGenGo        = lookupEnv("PROTOC_GENGO", getGoBin()+"/protoc-gen-go")
	ProtocGenGoPath    = "google.golang.org/protobuf/cmd/protoc-gen-go"
	ProtocGenGoVersion = lookupEnv("PROTOC_GENGO_VERSION", "v1.28")

	ProtocGenGoGRPC        = lookupEnv("PROTOC_GENGO_GRPC", getGoBin()+"/protoc-gen-go-grpc")
	ProtocGenGoGRPCPath    = "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	ProtocGenGoGRPCVersion = lookupEnv("PROTOC_GENGO_GRPC_VERSION", "v1.2")
)

func installPackage(path string, version string, binary string) error {
	packagePath := fmt.Sprintf("%s@%s", path, version)
	if _, err := os.Stat(binary); err != nil {
		if os.IsNotExist(err) {
			packagePath := fmt.Sprintf("%s@%s", path, version)
			fmt.Printf("installing %s to %s\n", packagePath, binary)
			return sh.Run("go", "install", packagePath)
		}
		return err
	}
	fmt.Printf("%s already installed in %s, skipping install\n", packagePath, binary)
	return nil
}

var getGoBin = func() string {
	path, err := sh.Output("go", "env", "GOPATH")
	if err != nil {
		mg.Fatalf(1, "error getting go path: %s", err) //nolint:errcheck
	}
	return fmt.Sprintf("%s/bin", path)
}
