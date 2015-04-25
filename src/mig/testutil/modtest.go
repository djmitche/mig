package testutil

import (
    "testing"
    "mig"
    "encoding/json"
)

func unexpectedErrors(t *testing.T, errors []string) {
    for _, msg := range errors {
        t.Errorf("Unexpected error: %s", msg)
    }
}

func RunModule(t *testing.T, m mig.Moduler, args interface{}) (mig.ModuleResult, []byte) {
    var marshalled_args []byte
    marshalled_args, err := json.Marshal(args)
    if (err != nil) {
        t.Error(err)
    }

    t.Logf("Running module with args %s", marshalled_args)
    marshalled_res := []byte(m.Run(marshalled_args))
    t.Logf("Result: %s", marshalled_res)

    var res mig.ModuleResult
    err = json.Unmarshal(marshalled_res, &res)
    if (err != nil) {
        t.Error(err)
    }
    return res, marshalled_res
}

func AssertModuleSucceeded(t *testing.T, res mig.ModuleResult) {
    if (len(res.Errors) > 0) {
        unexpectedErrors(t, res.Errors)
        t.FailNow()
    }
    if !res.Success {
        t.Fatal("Not successful")
    }
}

func AssertModuleError(t *testing.T, res mig.ModuleResult, exp_error string) {
    if res.Success {
        t.Fatal("Unexpected success")
    }
    if len(res.Errors) == 0 {
        t.Fatal("No errors")
    }
    if len(res.Errors) > 1 {
        unexpectedErrors(t, res.Errors)
        t.FailNow()
    }

    if res.Errors[0] != exp_error {
        unexpectedErrors(t, res.Errors)
        t.FailNow()
    }
}
