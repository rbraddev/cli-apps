package scan_test

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"pScan/scan"
)

func TestAdd(t *testing.T) {
	testCases := []struct {
		name      string
		host      string
		expectLen int
		expectErr error
	}{
		{"AddNew", "host2", 2, nil},
		{"AddExisting", "host1", 1, scan.ErrExists},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hl := &scan.HostList{}

			if err := hl.Add("host1"); err != nil {
				t.Fatal(err)
			}

			err := hl.Add(tc.host)

			if tc.expectErr != nil {
				if err == nil {
					t.Fatalf("want error, got nil\n")
				}

				if !errors.Is(err, tc.expectErr) {
					t.Errorf("want error %q, got %q\n", tc.expectErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("want no error, got %q\n", err)
			}

			if len(hl.Hosts) != tc.expectLen {
				t.Errorf("want list length %d, got %d\n", tc.expectLen, len(hl.Hosts))
			}

			if hl.Hosts[1] != tc.host {
				t.Errorf("want hostname %q as index 1, got %q\n", tc.host, hl.Hosts[1])
			}
		})
	}
}

func TestRemove(t *testing.T) {
	testCases := []struct {
		name      string
		host      string
		expectLen int
		expectErr error
	}{
		{"RemoveExisting", "host1", 1, nil},
		{"RemoveNotFound", "host3", 1, scan.ErrNotExists},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hl := &scan.HostList{}

			for _, h := range []string{"host1", "host2"} {
				if err := hl.Add(h); err != nil {
					t.Fatal(err)
				}
			}

			err := hl.Remove(tc.host)

			if tc.expectErr != nil {
				if err == nil {
					t.Fatal("want error, got nil\n")
				}

				if !errors.Is(err, tc.expectErr) {
					t.Errorf("want error %q, got %q\n", tc.expectErr, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("want no error, got %q\n", err)
			}

			if len(hl.Hosts) != tc.expectLen {
				t.Errorf("expect list length %d, got %d\n", tc.expectLen, len(hl.Hosts))
			}

			if hl.Hosts[0] == tc.host {
				t.Errorf("host name %q should not be in the list\n", tc.host)
			}
		})
	}
}

func TestSaveLoad(t *testing.T) {
	hl1 := scan.HostList{}
	hl2 := scan.HostList{}

	hostName := "host1"
	hl1.Add(hostName)

	tf, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("error creating temp file: %s", err)
	}
	defer os.Remove(tf.Name())

	if err := hl1.Save(tf.Name()); err != nil {
		t.Fatalf("error saving list to file: %s", err)
	}

	if err := hl2.Load(tf.Name()); err != nil {
		t.Fatalf("error getting list from file: %s", err)
	}

	if hl1.Hosts[0] != hl2.Hosts[0] {
		t.Errorf("host %q should match %q host", hl1.Hosts[0], hl2.Hosts[0])
	}
}

func TestLoadNoFile(t *testing.T) {
	tf, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("error creating temp file: %s", err)
	}

	if err := os.Remove(tf.Name()); err != nil {
		t.Fatalf("error deleting temp file: %s", err)
	}

	hl := &scan.HostList{}

	if err := hl.Load(tf.Name()); err != nil {
		t.Errorf("want no error, got %q\n", err)
	}
}
