package main

import (
	"fmt"
	"io"
	"log"

	"github.com/orlangure/errors"
)

func main() {
	if err := internalFunctionCallsAnotherFunction("joseph"); err != nil {
		log.Println("The original error is created outside of our app:")
		log.Println(err)
		log.Println(errors.Trace(err))
	}

	fmt.Println()

	if err := entryPointOfAnnotatedError(42); err != nil {
		log.Println("The original error is created by us, without any reason:")
		log.Println(err)
		log.Println(errors.Trace(err))
	}
}

// the original error created outside of our world
func internalFunctionCallsAnotherFunction(user string) error {
	if err := internalFunctionReturnsError(); err != nil {
		return errors.WithField("user", "joseph").Wrap(err, "internal function call failed")
	}

	return nil
}
func internalFunctionReturnsError() error {
	if err := externalPackageReturnsError(); err != nil {
		return errors.WithField("foo", "bar").Wrap(err, "external package call failed")
	}

	return nil
}
func externalPackageReturnsError() error {
	return io.EOF
}

// the original error is created by us
func entryPointOfAnnotatedError(id int) error {
	if err := firstCallToAnnotatedErrorFunction(); err != nil {
		return errors.WithField("id", id).Wrap(err, "action on an id failed")
	}

	return nil
}
func firstCallToAnnotatedErrorFunction() error {
	if err := internalFunctionReturnsAnnotatedError(); err != nil {
		return errors.WithField("resource", "container").Wrap(err, "operation on a resource failed")
	}

	return nil
}
func internalFunctionReturnsAnnotatedError() error {
	return errors.WithField("root", "cause").New("this is expected")
}
