package vaku_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/lingrino/vaku/vaku"
)

var seededOnce = false
var targetSeededOnce = false

// Initialize a new simple vault client to be used for tests
func clientInitForTests(t *testing.T) *vaku.Client {
	return clientInitForTestsCommon(t, vaultToken, vaultAddr, "VAKU_VAULT_ADDR", &seededOnce)
}

func clientInitForTestsCommon(t *testing.T, token string, address string, addressEnvVar string, seeded *bool) *vaku.Client {
	var err error
	// Initialize a new vault client
	vclient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		t.Fatal(fmt.Errorf("failed to create a vault client for testing: %w", err))
	}
	// Initialize a new vaku client and attach the vault client
	client := vaku.NewClient()
	client.Client = vclient
	// Set the token to the test value
	client.SetToken(token)
	// Set the address to the env var VAKU_VAULT_ADDR or the default constant
	err = client.SetAddress(address)
	if err != nil {
		t.Errorf("failed to set client address to %s", address)
	}
	if os.Getenv(addressEnvVar) != "" {
		err = client.SetAddress(os.Getenv(addressEnvVar))
		if err != nil {
			t.Errorf("failed to set client address to %s", os.Getenv(addressEnvVar))
		}
	} else {
		err = client.SetAddress(address)
		if err != nil {
			t.Errorf("failed to set client address to %s", address)
		}
	}

	// Seed the client if it has never been seeded
	if !*seeded {
		err := seed(t, client)
		if err != nil {
			t.Fatal(fmt.Errorf("failed to seed the vault client: %w", err))
		}
		*seeded = true
	}

	return client
}

// Initialize a new copy client to be used for tests
func copyClientInitForTests(t *testing.T) *vaku.CopyClient {
	// Initialize a new copy client and attach the source and target client
	copyClient := vaku.NewCopyClient()

	copyClient.Source = clientInitForTestsCommon(t, vaultToken, vaultAddr, "VAKU_VAULT_ADDR", &seededOnce)
	copyClient.Target = clientInitForTestsCommon(t, targetVaultToken, targetVaultAddr, "VAKU_TARGET_VAULT_ADDR", &targetSeededOnce)
	return copyClient
}

// seed uses a client to write dummy data used for testing to vault.
// Strings generated here: https://www.random.org/strings
func seed(t *testing.T, c *vaku.Client) error {
	t.Helper()
	var err error

	// Turn on logging to stdout
	err = c.Sys().EnableAuditWithOptions("audit_stdout", &api.EnableAuditOptions{
		Type: "file",
		Options: map[string]string{
			"file_path": "stdout",
			"log_raw":   "true",
		},
	})
	if err != nil {
		// We don't care about errors trying to mount to a path that we have already
		// mounted to. A better option here would be to check if the mount exists
		// before attempting the mount, but this is only used in tests so it's not
		// worth the effort. Same with the next two error checks.
		if !strings.Contains(err.Error(), "path already in use") {
			t.Error(fmt.Errorf("failed to turn on vault logging: %w", err))
		}
	}

	// Mount the two secret backends
	err = c.Sys().Mount("secretv1/", &api.MountInput{
		Type: "kv",
		Options: map[string]string{
			"version": "1"},
	})
	if err != nil {
		if !strings.Contains(err.Error(), "path is already in use at secretv1/") {
			t.Error(fmt.Errorf("failed to mount secretv1/: %w", err))
		}
	}
	err = c.Sys().Mount("secretv2/", &api.MountInput{
		Type: "kv",
		Options: map[string]string{
			"version": "2",
		},
	})
	if err != nil {
		if !strings.Contains(err.Error(), "path is already in use at secretv2/") {
			t.Error(fmt.Errorf("failed to mount secretv2/: %w", err))
		}
	}

	seeds := map[string]map[string]interface{}{
		"test/foo": {
			"value": "bar",
		},
		"test/value": {
			"fizz": "buzz",
			"foo":  "bar",
		},
		"test/fizz": {
			"fizz": "buzz",
			"foo":  "bar",
		},
		"test/HToOeKKD": {
			"3zqxVbJY": "TvOjGxvC",
		},
		"test/inner/WKNC3muM": {
			"IY1C148K": "JxBfEt91",
			"iwVzPqbY": "0NH9GlR1",
		},
		"test/inner/A2xlzTfE": {
			"Eg5ljS7t": "BHRMKjj1",
			"quqr32S5": "pcidzSMW",
		},
		"test/inner/again/inner/UCrt6sZT": {
			"Eg5ljS7t": "6F1B5nBg",
			"quqr32S5": "81iY4HAN",
			"r6R0JUzX": "rs1mCRB5",
		},
	}

	v1Seeds := make(map[string]map[string]interface{})
	v2Seeds := make(map[string]map[string]interface{})
	for k, v := range seeds {
		v1Seeds[c.PathJoin("secretv1", k)] = v
		v2Seeds[c.PathJoin("secretv2", k)] = v
	}

	err = c.FolderWrite(v1Seeds)
	if err != nil {
		return fmt.Errorf("failed to seed secretv1 path: %w", err)
	}
	err = c.FolderWrite(v2Seeds)
	if err != nil {
		return fmt.Errorf("failed to seed secretv2 path: %w", err)
	}

	return err
}
