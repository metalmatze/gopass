package secret

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/justwatchcom/gopass/store"
)

const (
	yamlKey      = "bar"
	yamlValue    = "baz"
	yamlPassword = "zzz"
)

/*
- empty doc (get/set)
- only pw (get/set)
- no pw (get/set)
- pw and yaml (get/set)
- no sep / sep (get/set)
*/
func TestYAML(t *testing.T) {
	for _, tc := range []struct {
		name string
		tf   func(t *testing.T)
	}{
		{
			name: "Get Key from empty Secret",
			tf: func(t *testing.T) {
				s := &Secret{}
				_, err := s.Value(yamlKey)
				if err == nil {
					t.Errorf("Should complain about missing YAML marker")
				}
			},
		},
		{
			name: "Set Key to empty Secret",
			tf: func(t *testing.T) {
				s := &Secret{}
				// write key
				err := s.SetValue(yamlKey, yamlValue)
				if err != nil {
					t.Fatalf("%s", err)
				}
				// read back key
				content, err := s.Value(yamlKey)
				if err != nil {
					t.Fatalf("%s", err)
				}
				if string(content) != yamlValue {
					t.Errorf("Wrong value: %s", content)
				}
				// read back whole entry
				buf, err := s.Bytes()
				if err != nil {
					t.Fatalf("%s", err)
				}
				want := "\n---\n" + yamlKey + ": " + yamlValue + "\n"
				if string(buf) != want {
					t.Errorf("Wrong value: '%s' != '%s'", content, want)
				}
			},
		},
		{
			name: "Get key from password-only secret",
			tf: func(t *testing.T) {
				s, err := Parse([]byte(yamlPassword))
				if err != nil {
					t.Fatalf("faile to parse secret")
				}
				// read (non-existing) key
				_, err = s.Value(yamlKey)
				if err == nil {
					t.Errorf("Should complain about missing YAML marker")
				}
				// read back whole entry
				content, err := s.Bytes()
				if err != nil {
					t.Fatalf("%s", err)
				}
				want := string(yamlPassword)
				if string(content) != want+"\n" {
					t.Errorf("Wrong value: '%s' != '%s'", content, want)
				}
			},
		},
		{
			name: "Set key to password-only secret",
			tf: func(t *testing.T) {
				s, err := Parse([]byte(yamlPassword))
				if err != nil {
					t.Fatalf("%s", err)
				}
				// set new key
				err = s.SetValue(yamlKey, yamlValue)
				if err != nil {
					t.Fatalf("Failed to write new key: %s", err)
				}
				// read back the password
				if s.Password() != yamlPassword {
					t.Errorf("Wrong password: %s", s.Password())
				}
				// read back the key
				content, err := s.Value(yamlKey)
				if err != nil {
					t.Fatalf("Failed to read key %s: %s", yamlKey, err)
				}
				if string(content) != yamlValue {
					t.Errorf("Wrong value: %s", content)
				}
				// read back whole entry
				bv, err := s.Bytes()
				if err != nil {
					t.Fatalf("%s", err)
				}
				want := yamlPassword + "\n---\nbar: baz\n"
				if string(bv) != want {
					t.Errorf("Wrong value: '%s' != '%s'", content, want)
				}
			},
		},
		{
			name: "Bare YAML - no document marker - read key",
			tf: func(t *testing.T) {
				in := "bar: baz\nzab: 123\n"
				s, err := Parse([]byte(in))
				if err != nil {
					t.Fatalf("%s", err)
				}
				// read back a key
				_, err = s.Value(yamlKey)
				if err != store.ErrYAMLNoMark {
					t.Fatalf("Should fail to read YAML without document marker")
				}
				// read back whole entry
				content, err := s.Bytes()
				if err != nil {
					t.Fatalf("%s", err)
				}
				if string(content)+"\n" != in {
					t.Errorf("Wrong value: '%s' != '%s'", content, in)
				}
			},
		},
		{
			name: "Set multiple keys to a secret",
			tf: func(t *testing.T) {
				s, err := Parse([]byte(yamlPassword))
				if err != nil {
					t.Fatalf("%s", err)
				}
				want := yamlPassword + "\n---\n"
				numKey := 100
				for i := 0; i < numKey; i++ {
					// set key
					key := fmt.Sprintf("%s-%d", yamlKey, i)
					if err := s.SetValue(key, yamlValue); err != nil {
						t.Fatalf("Failed to write new key: %s", err)
						continue
					}
					want += key + ": " + yamlValue + "\n"
				}
				// read back the password
				if s.Password() != yamlPassword {
					t.Errorf("Wrong password: %s", s.Password())
				}
				// read back the keys
				for i := 0; i < numKey; i++ {
					key := yamlKey + "-" + strconv.Itoa(i)
					content, err := s.Value(key)
					if err != nil {
						t.Fatalf("Failed to read key %s: %s", key, err)
					}
					if content != yamlValue {
						t.Errorf("Wrong value: %s", content)
					}
				}
				// read back whole entry
				content, err := s.Bytes()
				if err != nil {
					t.Fatalf("%s", err)
				}
				if string(content) != want {
					t.Errorf("Wrong value: '%s' != '%s'", content, want)
				}
			},
		},
		{
			name: "Get Multi-Line Value containing three dashes",
			tf: func(t *testing.T) {
				mlValue := `-----BEGIN PGP PRIVATE KEY BLOCK-----
aaa
bbb
ccc
-----END PGP PRIVATE KEY BLOCK-----`
				s := &Secret{}
				// write key
				err := s.SetValue(yamlKey, mlValue)
				if err != nil {
					t.Fatalf("%s", err)
				}
				// read back key
				content, err := s.Value(yamlKey)
				if err != nil {
					t.Fatalf("%s", err)
				}
				if string(content) != mlValue {
					t.Errorf("Wrong value: '%s' - Expected: '%s'", content, mlValue)
				}
			},
		},
	} {
		// run test case
		t.Run(tc.name, tc.tf)
	}
}
