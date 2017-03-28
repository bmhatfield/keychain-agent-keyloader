package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	// Actually depends on the "support-description" branch
	// See https://github.com/keybase/go-keychain/pull/7
	keychain "github.com/bmhatfield/go-keychain"

	"golang.org/x/crypto/ssh/agent"
)

func getKeychainCredential(name string) ([]byte, error) {
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetService(name)
	query.SetDescription("ssh key passphrase")
	query.SetMatchLimit(keychain.MatchLimitOne)
	query.SetReturnData(true)

	results, err := keychain.QueryItem(query)

	if err != nil {
		return []byte(""), err
	}

	if len(results) != 1 {
		return []byte(""), fmt.Errorf("Did not find passphrase for %s in Keychain", name)
	}

	return results[0].Data, nil
}

func main() {
	keyPath := os.ExpandEnv("${HOME}/.ssh/keys/TODO-KEYPATH")

	pemBytes, err := ioutil.ReadFile(keyPath)

	if err != nil {
		fmt.Println("Error loading key file: ", err)
		os.Exit(1)
	}

	pemBlock, _ := pem.Decode(pemBytes)

	privateKeyBytes := pemBlock.Bytes

	if x509.IsEncryptedPEMBlock(pemBlock) {
		passphrase, err := getKeychainCredential(keyPath)

		if err != nil {
			fmt.Println("Unable to retrieve passphrase from Keychain: ", err)
			os.Exit(1)
		}

		privateKeyBytes, err = x509.DecryptPEMBlock(pemBlock, passphrase)

		if err != nil {
			fmt.Println("Error decrypting private key: ", err)
			os.Exit(1)
		}
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBytes)

	if err != nil {
		fmt.Println("Error parsing private key: ", err)
		os.Exit(1)
	}

	agentKey := agent.AddedKey{
		PrivateKey: privateKey,
	}

	sshAuthSock, exists := os.LookupEnv("SSH_AUTH_SOCK")

	if exists {
		sshAuthSocket, err := net.Dial("unix", sshAuthSock)

		if err != nil {
			fmt.Println("Error connecting to SSH Agent: ", err)
			os.Exit(1)
		}

		sshAgent := agent.NewClient(sshAuthSocket)

		err = sshAgent.Add(agentKey)

		if err != nil {
			fmt.Println("Error adding private key to agent: ", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("SSH_AUTH_SOCK not set")
		os.Exit(1)
	}

	fmt.Printf("Successfully added %s to running ssh-agent\n", keyPath)
}
