package passtor

import (
	"github.com/sethvargo/go-diceware/diceware"
	"strings"
)

func Passphrase() (string, error) {
	words, err := diceware.Generate(PASSPHRASELENGHT)

	if err != nil {
		return "", err
	}

	return strings.Join(words, PASSPHRASESEP), nil
}
