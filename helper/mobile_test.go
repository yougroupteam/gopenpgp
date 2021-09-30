package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yougroupteam/gopenpgp/v2/constants"
	"github.com/yougroupteam/gopenpgp/v2/crypto"
)

func TestMobileSignedMessageDecryption(t *testing.T) {
	privateKey, _ := crypto.NewKeyFromArmored(readTestFile("keyring_privateKey", false))
	// Password defined in base_test
	privateKey, err := privateKey.Unlock(testMailboxPassword)
	if err != nil {
		t.Fatal("Expected no error unlocking privateKey, got:", err)
	}
	testPrivateKeyRing, _ := crypto.NewKeyRing(privateKey)

	publicKey, _ := crypto.NewKeyFromArmored(readTestFile("mime_publicKey", false))
	testPublicKeyRing, _ := crypto.NewKeyRing(publicKey)

	pgpMessage, err := crypto.NewPGPMessageFromArmored(readTestFile("message_signed", false))
	if err != nil {
		t.Fatal("Expected no error when unarmoring, got:", err)
	}

	decrypted, err := DecryptExplicitVerify(pgpMessage, testPrivateKeyRing, testPublicKeyRing, crypto.GetUnixTime())
	if err != nil {
		t.Fatal("Expected no error when decrypting, got:", err)
	}

	assert.Exactly(t, constants.SIGNATURE_NO_VERIFIER, decrypted.SignatureVerificationError.Status)
	assert.Exactly(t, readTestFile("message_plaintext", true), decrypted.Message.GetString())

	publicKey, _ = crypto.NewKeyFromArmored(readTestFile("keyring_publicKey", false))
	testPublicKeyRing, _ = crypto.NewKeyRing(publicKey)

	pgpMessage, err = testPublicKeyRing.Encrypt(decrypted.Message, testPrivateKeyRing)
	if err != nil {
		t.Fatal("Expected no error when encrypting, got:", err)
	}

	decrypted, err = DecryptExplicitVerify(pgpMessage, testPrivateKeyRing, testPublicKeyRing, crypto.GetUnixTime())
	if err != nil {
		t.Fatal("Expected no error when decrypting, got:", err)
	}

	assert.Nil(t, decrypted.SignatureVerificationError)
	assert.Exactly(t, readTestFile("message_plaintext", true), decrypted.Message.GetString())

	decrypted, err = DecryptExplicitVerify(pgpMessage, testPublicKeyRing, testPublicKeyRing, crypto.GetUnixTime())
	assert.NotNil(t, err)
	assert.Nil(t, decrypted)
}

func TestMobileSignedMessageDecryptionWithSessionKey(t *testing.T) {
	var message = crypto.NewPlainMessageFromString(
		"The secret code is... 1, 2, 3, 4, 5. I repeat: the secret code is... 1, 2, 3, 4, 5",
	)

	privateKey, _ := crypto.NewKeyFromArmored(readTestFile("keyring_privateKey", false))
	// Password defined in base_test
	privateKey, err := privateKey.Unlock(testMailboxPassword)
	if err != nil {
		t.Fatal("Expected no error unlocking privateKey, got:", err)
	}
	testPrivateKeyRing, _ := crypto.NewKeyRing(privateKey)

	publicKey, _ := crypto.NewKeyFromArmored(readTestFile("keyring_publicKey", false))
	testPublicKeyRing, _ := crypto.NewKeyRing(publicKey)

	sk, err := crypto.GenerateSessionKey()
	if err != nil {
		t.Fatal("Expected no error generating session key, got:", err)
	}

	pgpMessage, err := sk.Encrypt(message)
	if err != nil {
		t.Fatal("Expected no error when unarmoring, got:", err)
	}

	decrypted, err := DecryptSessionKeyExplicitVerify(pgpMessage, sk, testPublicKeyRing, crypto.GetUnixTime())
	if err != nil {
		t.Fatal("Expected no error when decrypting, got:", err)
	}

	assert.Exactly(t, constants.SIGNATURE_NOT_SIGNED, decrypted.SignatureVerificationError.Status)
	assert.Exactly(t, message.GetString(), decrypted.Message.GetString())

	publicKey, _ = crypto.NewKeyFromArmored(readTestFile("keyring_publicKey", false))
	testPublicKeyRing, _ = crypto.NewKeyRing(publicKey)

	pgpMessage, err = sk.EncryptAndSign(message, testPrivateKeyRing)
	if err != nil {
		t.Fatal("Expected no error when encrypting, got:", err)
	}

	decrypted, err = DecryptSessionKeyExplicitVerify(pgpMessage, sk, testPublicKeyRing, crypto.GetUnixTime())
	if err != nil {
		t.Fatal("Expected no error when decrypting, got:", err)
	}

	assert.Nil(t, decrypted.SignatureVerificationError)
	assert.Exactly(t, message.GetString(), decrypted.Message.GetString())
}

func TestGetJsonSHA256FingerprintsV4(t *testing.T) {
	sha256Fingerprints, err := GetJsonSHA256Fingerprints(readTestFile("keyring_publicKey", false))
	if err != nil {
		t.Fatal("Cannot unarmor key:", err)
	}

	assert.Exactly(t, []byte("[\"d9ac0b857da6d2c8be985b251a9e3db31e7a1d2d832d1f07ebe838a9edce9c24\",\"203dfba1f8442c17e59214d9cd11985bfc5cc8721bb4a71740dd5507e58a1a0d\"]"), sha256Fingerprints)
}
