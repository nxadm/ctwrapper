package main

import (
	"errors"
	"github.com/hashicorp/vault/api"
	"strings"
)

func retrieveVaultSecret(path string) (string, error) {

	// Separate vaultPath and key
	split := strings.SplitAfter(path, "/")
	backendAndPath := strings.Join(split[0:len(split)-1], "")
	key := split[len(split)-1]

	/* Retrieve VAULT_ADDR, VAULT_TOKEN and other VAULT_* env variables */
	vaultConfig := api.DefaultConfig()
	if err := vaultConfig.ReadEnvironment(); err != nil {
		return "", err
	}

	/* Retrieve the secret */
	client, err := api.NewClient(vaultConfig)
	if err != nil {
		return "", err
	}
	secretsRaw, err := client.Logical().Read(backendAndPath)
	if err != nil {
		return "", err
	}
	if secretsRaw == nil {
		return "", errors.New("can not find the requested secret")
	}

	var secret string
	for k, v := range secretsRaw.Data {
		if k == key {
			secret = v.(string)
			break
		}
	}

	return secret, nil
}
