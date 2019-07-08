package ethereum

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEthereum(t *testing.T) {
	ethChain, err := NewSimulatedBlockchain()
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second) // wait for smart contract to deploy

	t.Run("StartsOutEmpty", func(t *testing.T) {
		zero := big.NewInt(0)
		serverDetails, err := ethChain.FetchServers(*zero)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 0, len(serverDetails), "expected no server details")
	})

	t.Run("CanRegisterServer", func(t *testing.T) {
		target := "target"
		cert := "cert"

		err := ethChain.RegisterServer(target, cert)
		if err != nil {
			t.Fatal(err)
		}

		zero := big.NewInt(0)
		serverDetails, err := ethChain.FetchServers(*zero)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 1, len(serverDetails), "expected server details")
		assert.Equal(t, target, serverDetails[0].Target)
		assert.Equal(t, cert, serverDetails[0].Cert)
	})

	t.Run("CanRegisterMultipleServers", func(t *testing.T) {
		target := "target"
		cert := "cert"

		err := ethChain.RegisterServer(target, cert)
		if err != nil {
			t.Fatal(err)
		}

		err = ethChain.RegisterServer(target, cert)
		if err != nil {
			t.Fatal(err)
		}

		zero := big.NewInt(0)
		serverDetails, err := ethChain.FetchServers(*zero)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 3, len(serverDetails), "expected server details")
	})

	t.Run("FiltersOutOldEntries", func(t *testing.T) {
		veryHighLaterThan := new(big.Int).Exp(big.NewInt(2), big.NewInt(64), nil)
		serverDetails, err := ethChain.FetchServers(*veryHighLaterThan)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 0, len(serverDetails), "expected no server details")
	})
}
