package example

import (
	"encoding/json"
	"mig/testutil"
	"net"
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

func TestRunLookupHost(t *testing.T) {
	r := new(Runner)
	r.lookupHostFn = func(name string) (addrs []string, err error) {
		testutil.Assert(t, name == "myhost", "wrong hostname passed to LookupHost")
		var rv []string
		rv = append(rv, "1.2.3.4")
		rv = append(rv, "5.6.7.8")
		return rv, nil
	}
	res, marshalled := testutil.RunModule(t, r, params{false, false, "myhost"})
	testutil.AssertModuleSucceeded(t, res)
	testutil.Assert(t, res.FoundAnything, "FoundAnything not set")

	modres := unmarshal(t, marshalled)
	testutil.Assert(t, modres.Statistics.StuffFound == 2, "statistics.StuffFound != 2")
	testutil.Assertf(t, modres.Elements.LookedUpHost == "myhost=1.2.3.4, 5.6.7.8",
		"LookedUpHost is incorrect: %v", modres.Elements.LookedUpHost)
}

func TestRunGetHostname(t *testing.T) {
	r := new(Runner)
	r.hostnameFn = func() (string, error) {
		return "myhostname", nil
	}
	res, marshalled := testutil.RunModule(t, r, params{true, false, ""})
	testutil.AssertModuleSucceeded(t, res)
	testutil.Assert(t, res.FoundAnything, "FoundAnything not set")

	modres := unmarshal(t, marshalled)
	testutil.Assert(t, modres.Statistics.StuffFound == 1, "statistics.StuffFound != 1")
	testutil.Assert(t, modres.Elements.Hostname == "myhostname", "HostName is empty")
}

func TestRunGetAddresses(t *testing.T) {
	r := new(Runner)
	r.interfaceAddrsFn = func() ([]net.Addr, error) {
		ip := net.IPNet{IP: net.IPv4(128, 135, 11, 1), Mask: net.CIDRMask(24, 32)}
		var ifat []net.Addr
		ifat = append(ifat, &ip)

		return ifat, nil
	}
	res, marshalled := testutil.RunModule(t, r, params{false, true, ""})
	testutil.AssertModuleSucceeded(t, res)
	testutil.Assert(t, res.FoundAnything, "FoundAnything not set")

	modres := unmarshal(t, marshalled)
	testutil.Assert(t, modres.Statistics.StuffFound == 1, "statistics.StuffFound != 1")
	testutil.Assert(t, modres.Elements.Addresses[0] == "128.135.11.1/24", "Addresses is empty")
}
