package crypto

import (
	"errors"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yougroupteam/gopenpgp/v2/constants"
)

const signedPlainText = "Signed message\n"

var textSignature, binSignature *PGPSignature
var message *PlainMessage
var signatureTest = regexp.MustCompile("(?s)^-----BEGIN PGP SIGNATURE-----.*-----END PGP SIGNATURE-----$")

func TestSignTextDetached(t *testing.T) {
	var err error

	message = NewPlainMessageFromString(signedPlainText)
	textSignature, err = keyRingTestPrivate.SignDetached(message)
	if err != nil {
		t.Fatal("Cannot generate signature:", err)
	}

	armoredSignature, err := textSignature.GetArmored()
	if err != nil {
		t.Fatal("Cannot armor signature:", err)
	}

	assert.Regexp(t, signatureTest, armoredSignature)
}

func TestVerifyTextDetachedSig(t *testing.T) {
	verificationError := keyRingTestPublic.VerifyDetached(message, textSignature, testTime)
	if verificationError != nil {
		t.Fatal("Cannot verify plaintext signature:", verificationError)
	}
}

func TestVerifyTextDetachedSigWrong(t *testing.T) {
	fakeMessage := NewPlainMessageFromString("wrong text")
	verificationError := keyRingTestPublic.VerifyDetached(fakeMessage, textSignature, testTime)

	assert.EqualError(t, verificationError, "Signature Verification Error: Invalid signature")

	err := &SignatureVerificationError{}
	_ = errors.As(verificationError, err)
	assert.Exactly(t, constants.SIGNATURE_FAILED, err.Status)
}

func TestSignBinDetached(t *testing.T) {
	var err error

	message = NewPlainMessage([]byte(signedPlainText))
	binSignature, err = keyRingTestPrivate.SignDetached(message)
	if err != nil {
		t.Fatal("Cannot generate signature:", err)
	}

	armoredSignature, err := binSignature.GetArmored()
	if err != nil {
		t.Fatal("Cannot armor signature:", err)
	}

	assert.Regexp(t, signatureTest, armoredSignature)
}

func TestVerifyBinDetachedSig(t *testing.T) {
	verificationError := keyRingTestPublic.VerifyDetached(message, binSignature, testTime)
	if verificationError != nil {
		t.Fatal("Cannot verify binary signature:", verificationError)
	}
}
