package ping

import (
	//	"encoding/json"
	"fmt"
	"github.com/bmizerany/assert"
	//"mig/testutil"
	"testing"
)

type validatetest struct {
	LookupHost string
	valid      bool
}

func TestValidateParametersBadProtocol(t *testing.T) {
	r := new(Runner)
	r.Parameters.Protocol = "not-a-protocol"
	err := r.ValidateParameters()
	assert.Equal(t, err, fmt.Errorf("not-a-protocol is not a supported ping protocol"))
}

func TestValidateParametersBadPort(t *testing.T) {
	var bad_ports = []int{-10, 65539}
	var protocols = []string{"tcp", "udp"}

	r := new(Runner)
	for _, proto := range protocols {
		for _, port := range bad_ports {
			r.Parameters.Protocol = proto
			r.Parameters.DestinationPort = float64(port)
			err := r.ValidateParameters()
			assert.Equal(t, err,
				fmt.Errorf("%s ping requires a valid destination port between 0 and 65535, got %d",
					proto, port))
		}
	}
}

func TestValidateParametersRawIP(t *testing.T) {
	r := new(Runner)
    r.Parameters.Protocol = "icmp"
    r.Parameters.Destination = "127.0.0.1"
    err := r.ValidateParameters()
    assert.Equal(t, err, nil)
    assert.Equal(t, r.Parameters.ipDest, "127.0.0.1")
}

func TestValidateParametersInvalidIP(t *testing.T) {
	r := new(Runner)
    r.Parameters.Protocol = "icmp"
    r.Parameters.Destination = "127.0.0"
    err := r.ValidateParameters()
    assert.Equal(t, err, fmt.Errorf("destination IP is invalid: 127.0.0"))
}

func TestValidateParametersDefaults(t *testing.T) {
	r := new(Runner)
    r.Parameters.Protocol = "icmp"
    r.Parameters.Destination = "127.0.0.1"
    err := r.ValidateParameters()
    assert.Equal(t, err, nil)
    assert.Equal(t, r.Parameters.Timeout, 5.0)
    assert.Equal(t, r.Parameters.Count, 3.0)
}
