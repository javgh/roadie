package ethereum

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	maxAge = big.NewInt(24 * 60 * 60)
)

func TestEthereum(t *testing.T) {
	ethChain, err := NewSimulatedBlockchain()
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second) // wait for smart contract to deploy

	t.Run("StartsOutEmpty", func(t *testing.T) {
		serverDetails, err := ethChain.FetchServers(*maxAge)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 0, len(serverDetails), "expected no server details")
	})

	t.Run("CanRegisterServer", func(t *testing.T) {
		target := "target"
		cert := []byte{}

		err := ethChain.RegisterServer(target, cert)
		if err != nil {
			t.Fatal(err)
		}

		serverDetails, err := ethChain.FetchServers(*maxAge)
		if err != nil {
			t.Fatal(err)
		}

		require.Equal(t, 1, len(serverDetails), "expected server details")
		assert.Equal(t, target, serverDetails[0].Target)
		assert.Equal(t, cert, serverDetails[0].Cert)
	})

	t.Run("CanRegisterMultipleServers", func(t *testing.T) {
		target := "target"
		cert := []byte{}

		err := ethChain.RegisterServer(target, cert)
		if err != nil {
			t.Fatal(err)
		}

		err = ethChain.RegisterServer(target, cert)
		if err != nil {
			t.Fatal(err)
		}

		serverDetails, err := ethChain.FetchServers(*maxAge)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 3, len(serverDetails), "expected server details")
	})

	t.Run("FiltersOutOldEntries", func(t *testing.T) {
		zeroMaxAge := big.NewInt(0)
		serverDetails, err := ethChain.FetchServers(*zeroMaxAge)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 0, len(serverDetails), "expected no server details")
	})
}
