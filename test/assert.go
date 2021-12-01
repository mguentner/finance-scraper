// +build unit

package test

import (
	"fmt"
	"reflect"
	"testing"
)

// taken from https://github.com/benbjohnson/testing

// Assert fails the test if the condition is false.
func Assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	tb.Helper()

	if !condition {
		tb.Fatalf(msg, v...)
	}
}

// Failed fails the test if an err is nil.
func Failed(tb testing.TB, err error) {
	tb.Helper()

	if err == nil {
		tb.Fatalf("expected error, but was nil")
	}
}

// Ok fails the test if an err is not nil.
func Ok(tb testing.TB, err error) {
	tb.Helper()

	if err != nil {
		tb.Fatalf("unexpected error: %s", err.Error())
	}
}

// Equals fails the test if exp is not equal to act.
func Equals(tb testing.TB, exp, act interface{}) {
	tb.Helper()

	if !reflect.DeepEqual(exp, act) {
		_, isExpStringer := exp.(fmt.Stringer)
		_, isActStringer := act.(fmt.Stringer)
		if isExpStringer && isActStringer {
			tb.Fatalf("exp: %s\n\n\tgot: %s", exp, act)
		} else {
			tb.Fatalf("exp: %#v\n\n\tgot: %#v", exp, act)
		}
	}
}
