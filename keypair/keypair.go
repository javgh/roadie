package keypair

import (
	"github.com/HyperspaceApp/ed25519"
)

type Keypair struct {
	PubKey  ed25519.PublicKey
	PrivKey ed25519.PrivateKey
}

func Generate() (Keypair, error) {
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	return Keypair{PubKey: pubKey, PrivKey: privKey}, err
}

func JointSignAlice(aliceKeypair Keypair, bobPubKey ed25519.PublicKey,
	noncePoints []ed25519.CurvePoint, msg []byte) ([]byte, error) {
	pubKeys := []ed25519.PublicKey{aliceKeypair.PubKey, bobPubKey}
	aliceN := 0
	jointPrivateKey, err := ed25519.GenerateJointPrivateKey(
		pubKeys, aliceKeypair.PrivKey, aliceN)
	if err != nil {
		return nil, err
	}

	return ed25519.JointSign(
		aliceKeypair.PrivKey, jointPrivateKey, noncePoints, msg), nil
}

func JointSignBob(bobKeypair Keypair, alicePubKey ed25519.PublicKey,
	noncePoints []ed25519.CurvePoint, msg []byte) ([]byte, error) {
	pubKeys := []ed25519.PublicKey{alicePubKey, bobKeypair.PubKey}
	bobN := 1
	jointPrivateKey, err := ed25519.GenerateJointPrivateKey(
		pubKeys, bobKeypair.PrivKey, bobN)
	if err != nil {
		return nil, err
	}

	return ed25519.JointSign(
		bobKeypair.PrivKey, jointPrivateKey, noncePoints, msg), nil
}

func JointSignWithAdaptorBob(bobKeypair Keypair, alicePubKey ed25519.PublicKey,
	noncePoints []ed25519.CurvePoint, adaptorPubKey ed25519.CurvePoint, msg []byte) ([]byte, error) {
	pubKeys := []ed25519.PublicKey{alicePubKey, bobKeypair.PubKey}
	bobN := 1
	jointPrivateKey, err := ed25519.GenerateJointPrivateKey(
		pubKeys, bobKeypair.PrivKey, bobN)
	if err != nil {
		return nil, err
	}

	return ed25519.JointSignWithAdaptor(
		bobKeypair.PrivKey, jointPrivateKey, noncePoints[0], noncePoints[1], adaptorPubKey, msg), nil
}

func JointSignWithAdaptorAlice(aliceKeypair Keypair, bobPubKey ed25519.PublicKey,
	noncePoints []ed25519.CurvePoint, adaptorPubKey ed25519.CurvePoint, msg []byte) ([]byte, error) {
	pubKeys := []ed25519.PublicKey{aliceKeypair.PubKey, bobPubKey}
	aliceN := 0
	jointPrivateKey, err := ed25519.GenerateJointPrivateKey(
		pubKeys, aliceKeypair.PrivKey, aliceN)
	if err != nil {
		return nil, err
	}

	return ed25519.JointSignWithAdaptor(
		aliceKeypair.PrivKey, jointPrivateKey, noncePoints[0], noncePoints[1], adaptorPubKey, msg), nil
}
