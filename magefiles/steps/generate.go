package steps

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Generate mg.Namespace

// Generates the protobuffers
func (Generate) Protobuf() error {
	mg.SerialDeps(validateProtocExecutable, installProtoc, Go{}.Tidy)

	args := []string{
		"--go_out=.",
		"--go_opt=paths=source_relative",
		"--go-grpc_out=.",
		"--go-grpc_opt=paths=source_relative",
		"--plugin=protoc-gen-go=" + ProtocGenGo,
		"--plugin=protoc-gen-go-grpc=" + ProtocGenGoGRPC,
		"apis/nginx/filter.proto",
	}
	return sh.RunV("protoc", args...)
}

func validateProtocExecutable() error {
	if err := sh.Run("which", "protoc"); err != nil {
		return fmt.Errorf("no protoc was found. Please check https://grpc.io/docs/protoc-installation/")
	}
	return nil
}

func installProtoc() error {
	if err := installPackage(ProtocGenGoPath, ProtocGenGoVersion, ProtocGenGo); err != nil {
		return err
	}

	return installPackage(ProtocGenGoGRPCPath, ProtocGenGoGRPCVersion, ProtocGenGoGRPC)
}
