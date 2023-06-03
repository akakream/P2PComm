package identity

import (
	"crypto/rand"
	"os"
	"path/filepath"

	crypto "github.com/libp2p/go-libp2p/core/crypto"
	peer "github.com/libp2p/go-libp2p/core/peer"
)

type Identity struct {
	KeyPath string
	PrivKey crypto.PrivKey
	PeerID  peer.ID
}

func NewIdentity(datapath string) (*Identity, error) {
	privateKeyPath := filepath.Join(datapath, "privateKey")
	var privKey crypto.PrivKey

	_, err := os.Stat(privateKeyPath)
	if os.IsNotExist(err) {
		privKey, _, err = crypto.GenerateRSAKeyPair(2048, rand.Reader)
		if err != nil {
			return nil, err
		}
		privKeyBytes, err := crypto.MarshalPrivateKey(privKey)
		if err != nil {
			return nil, err
		}
		err = os.WriteFile(privateKeyPath, privKeyBytes, 0400)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		key, err := os.ReadFile(privateKeyPath)
		if err != nil {
			return nil, err
		}
		privKey, err = crypto.UnmarshalPrivateKey(key)
		if err != nil {
			return nil, err
		}
	}

	pid, err := peer.IDFromPublicKey(privKey.GetPublic())
	if err != nil {
		return nil, err
	}

	return &Identity{KeyPath: privateKeyPath, PrivKey: privKey, PeerID: pid}, err
}
