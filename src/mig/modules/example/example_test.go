package example

import (
	"mig/testutil"
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

	r := new(Runner)

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

func TestRunInvalidLookupHost(t *testing.T) {
	res := testutil.RunModule(t, new(Runner), params{false, false, "in valid"})
	testutil.AssertModuleError(t, res,
		"ValidateParameters: LookupHost parameter is not a valid FQDN.")
}
