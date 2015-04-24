package example_test

import (
	"encoding/json"
	"mig/modules/example"
	"testing"
)

type validatetest struct {
	LookupHost string
	valid      bool
}

func TestValidateParameters(t *testing.T) {
	var hostnames = []validatetest{
		{"valid.hostname.com", true},
		{"in!valid-hostname", false},
		{"-invalid-hostname", false},
		{"invalid-hostname-", false},
		{"invalid_hostname", false},
		{"invalid hostname", false},
		{"http://mozilla.org", false},
	}

	r := new(example.Runner)

	for _, vt := range hostnames {
		r.Parameters.LookupHost = vt.LookupHost
		var err = r.ValidateParameters()
		if !vt.valid && err == nil {
			t.Errorf("invalid LookupHost '%s' considered valid", vt.LookupHost)
		} else if vt.valid && err != nil {
			t.Errorf("valid LookupHost '%s' considered invalid: %s", vt.LookupHost, err)
		}
	}
}

func tryRun(t *testing.T, r *example.Runner) (err error) {
	args, err := json.Marshal(r.Parameters)
	t.Logf("running module with args %s", string(args))
	if err != nil {
		return err
	}
	res := []byte(r.Run(args))
	return json.Unmarshal(res, &r.Results)
}

func TestRunInvalidLookupHost(t *testing.T) {
	r := new(example.Runner)
	r.Parameters.LookupHost = "in valid"
	err := tryRun(t, r)
	if err != nil {
		t.Error(err)
	}
	if r.Results.Success {
		t.Error("unexpected run success")
	}
	if len(r.Results.Errors) != 1 {
		t.Error("no error response")
	}
}
