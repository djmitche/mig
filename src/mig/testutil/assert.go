package testutil

import (
    "testing"
)

func Assert(t *testing.T, stmt bool, format string) {
    if !stmt {
        t.Fatal(format)
    }
}

func Assertf(t *testing.T, stmt bool, format string, args ...interface{}) {
    if !stmt {
        t.Fatalf(format, args...)
    }
}
