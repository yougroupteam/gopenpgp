package helper

import (
	"regexp"
	"testing"

	"github.com/yougroupteam/gopenpgp/v2/internal"

	"github.com/stretchr/testify/assert"
	"github.com/yougroupteam/gopenpgp/v2/crypto"
)

const inputPlainText = "  Signed message\n  \n  "
const signedPlainText = "  Signed message\n\n"

var signedMessageTest = regexp.MustCompile(
	"(?s)^-----BEGIN PGP SIGNED MESSAGE-----.*-----BEGIN PGP SIGNATURE-----.*-----END PGP SIGNATURE-----$")

func TestSignClearText(t *testing.T) {
	// Password defined in base_test
	armored, err := SignCleartextMessageArmored(
		readTestFile("keyring_privateKey", false),
		testMailboxPassword,
		inputPlainText,
	)

	if err != nil {
		t.Fatal("Cannot armor message:", err)
	}

	assert.Regexp(t, signedMessageTest, armored)

	verified, err := VerifyCleartextMessageArmored(
		readTestFile("keyring_publicKey", false),
		armored,
		crypto.GetUnixTime(),
	)
	if err != nil {
		t.Fatal("Cannot verify message:", err)
	}

	assert.Exactly(t, signedPlainText, verified)

	clearTextMessage, err := crypto.NewClearTextMessageFromArmored(armored)
	if err != nil {
		t.Fatal("Cannot parse message:", err)
	}
	assert.Exactly(t, internal.CanonicalizeAndTrim(inputPlainText), string(clearTextMessage.GetBinary()))
}
