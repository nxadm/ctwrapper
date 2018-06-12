package main

import (
	"github.com/hashicorp/vault/api"
)

func retrieveVaultSecret(path, key string) (string, error) {
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
