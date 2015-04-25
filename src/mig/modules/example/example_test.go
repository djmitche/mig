package example

import (
    "encoding/json"
    "testing"
    "mig/testutil"
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
		// an empty string indicates "don't look up a host" and is valid
		{"", true},
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

func unmarshal(t *testing.T, marshalled []byte) (res results) {
    err := json.Unmarshal(marshalled, &res)
    if err != nil {
        t.Fatal(err)
    }
    return
}

func TestRunInvalidLookupHost(t *testing.T) {
    res, _ := testutil.RunModule(t, new(Runner), params{false, false, "in valid"})
    testutil.AssertModuleError(t, res,
        "ValidateParameters: LookupHost parameter is not a valid FQDN.")
}

func TestRunGetHostname(t *testing.T) {
    res, marshalled := testutil.RunModule(t, new(Runner), params{true, false, ""})
    testutil.AssertModuleSucceeded(t, res)
    testutil.Assert(t, res.FoundAnything, "FoundAnything not set")
    
    modres := unmarshal(t, marshalled)
    testutil.Assert(t, modres.Statistics.StuffFound == 1, "statistics.StuffFound != 1")
    testutil.Assert(t, modres.Elements.Hostname != "", "HostName is empty")
}

func TestRunGetAddresses(t *testing.T) {
    res, marshalled := testutil.RunModule(t, new(Runner), params{false, true, ""})
    testutil.AssertModuleSucceeded(t, res)
    testutil.Assert(t, res.FoundAnything, "FoundAnything not set")
    
    modres := unmarshal(t, marshalled)
    testutil.Assert(t, modres.Statistics.StuffFound == 1, "statistics.StuffFound != 1")
    testutil.Assert(t, len(modres.Elements.Addresses) != 0, "Addresses is empty")
}
