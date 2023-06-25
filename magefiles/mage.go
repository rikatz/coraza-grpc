//go:build mage

package main

import (
	"github.com/magefile/mage/mg"

	//mage:import
	"github.com/rikatz/coraza-grpc/magefiles/steps"
)

// Validate runs all the validation steps
func Validate() {
	mg.Deps(steps.Validate{}.Unused, steps.Validate{}.Style)
	// Those are resource intensive and should be run Serialized
	mg.SerialDeps(steps.Go{}.Vet, steps.Validate{}.Lint)
}

// Runs the generation programs
func Generate() {
	mg.Deps(steps.Generate{}.Protobuf)
}
