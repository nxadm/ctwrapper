package main

import (
	"github.com/hashicorp/vault/api"
	"strings"
)

func retrieveVaultSecret(vaultInfo string) (string, error) {

	// Separate path and key
	split := strings.SplitAfter(vaultInfo, "/")
	path := strings.Join(split[0:len(split) - 2], "")
	key := split[len(split) - 1]

	// Retrieve the secret
	vaultConfig := api.DefaultConfig()
	if err := vaultConfig.ReadEnvironment(); err != nil {
		return "", err
	}
	client, err := api.NewClient(vaultConfig)
	if err != nil {
		return "", err
	}
	secretsRaw, err := client.Logical().Read(path)
	if err != nil {
		return "", err
	}
	var secret string
	for k, v := range secretsRaw.Data {
		if k == key {
			secret = v.(string)
		}
	}

	return secret, nil
}
