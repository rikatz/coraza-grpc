package steps

import (
	"fmt"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Validate mg.Namespace

// Unused verifies if there are any unused go modules
func (Validate) Unused() error {
	mg.Deps(Go{}.Tidy)
	fmt.Printf("checking differences on go.mod file")
	return sh.RunV("git", "diff", "--exit-code", "--", "go.sum", "go.mod")
}

// Style verifies if gofmt was executed properly
func (Validate) Style() error {
	fmt.Printf("checking if formatting was not executed\n")
	mg.Deps(Go{}.Fmt)
	return sh.RunV("git", "diff", "--exit-code")
}

// Lint runs the golangci-lint on the application
func (Validate) Lint() error {
	mg.Deps(installLinter)
	// From Makefile.common:
	// 'go list' needs to be executed before golangci-lint to prepopulate the modules cache.
	// Otherwise staticcheck code in golangci-lint might fail randomly for some reason not yet explained.
	args := "list -e -compiled -test=true -export=false -deps=true -find=false -tags= -- ./..."
	fmt.Printf("running go %s\n", args)
	if err := sh.RunV("go", strings.Split(args, " ")...); err != nil {
		return err
	}

	ciSha := lookupEnv("CI_MERGE_REQUEST_DIFF_BASE_SHA", "")
	if ciSha != "" {
		GolangLintOpts += " --new-from-rev=" + ciSha
	}
	argsCi := "run " + GolangLintOpts + " ./..."
	fmt.Printf("Running %s %s\n", GolangLint, argsCi)
	return sh.RunV(GolangLint, strings.Split(argsCi, " ")...)
}

func installLinter() error {
	return installPackage(GolangLintPath, GolangLintVersion, GolangLint)
}
