package scan_test

import (
	"net"
	"strconv"
	"testing"

	"pScan/scan"
)

func TestStateString(t *testing.T) {
	ps := scan.PortState{}

	if ps.Open.String() != "closed" {
		t.Errorf("want %q, got %q\n", "closed", ps.Open.String())
	}

	ps.Open = true

	if ps.Open.String() != "open" {
		t.Errorf("want %q, got %q\n", "open", ps.Open.String())
	}
}

func TestRunHostfound(t *testing.T) {
	testCases := []struct {
		name        string
		expectState string
	}{
		{"OpenPort", "open"},
		{"ClosedPort", "closed"},
	}

	host := "localhost"
	hl := &scan.HostList{}
	hl.Add(host)

	ports := []int{}

	for _, tc := range testCases {
		ln, err := net.Listen("tcp", net.JoinHostPort(host, "0"))
		if err != nil {
			t.Fatal(err)
		}
		defer ln.Close()

		_, portStr, err := net.SplitHostPort(ln.Addr().String())
		if err != nil {
			t.Fatal(err)
		}

		port, err := strconv.Atoi(portStr)
		if err != nil {
			t.Fatal(err)
		}

		ports = append(ports, port)

		if tc.name == "ClosedPort" {
			ln.Close()
		}
	}

	res := scan.Run(hl, ports)

	if len(res) != 1 {
		t.Fatalf("want 1 result, got %d\n", len(res))
	}

	if res[0].Host != host {
		t.Errorf("want host %q, got %q\n", host, res[0].Host)
	}

	if res[0].NotFound {
		t.Errorf("want host %q to be found\n", host)
	}

	if len(res[0].PortStates) != 2 {
		t.Fatalf("want 2 port statesm got %d\n", len(res[0].PortStates))
	}

	for i, tc := range testCases {
		if res[0].PortStates[i].Port != ports[i] {
			t.Errorf("want port %d, got %d\n", ports[0], res[0].PortStates[i].Port)
		}

		if res[0].PortStates[i].Open.String() != tc.expectState {
			t.Errorf("want port %d to be %s\n", ports[i], tc.expectState)
		}
	}
}

func TestRunHostNotFound(t *testing.T) {
	host := "389.389.389.389"
	hl := &scan.HostList{}
	hl.Add(host)

	res := scan.Run(hl, []int{})

	if len(res) != 1 {
		t.Fatalf("want 1 result, got %d\n", len(res))
	}

	if res[0].Host != host {
		t.Errorf("want host %q, got %q\n", host, res[0].Host)
	}

	if !res[0].NotFound {
		t.Errorf("want host %q NOT to be found\n", host)
	}

	if len(res[0].PortStates) != 0 {
		t.Fatalf("want 0 port states, got %d\n", len(res[0].PortStates))
	}
}
