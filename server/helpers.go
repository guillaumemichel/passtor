package main

import (
	"../../passtor"
)

type accountCountPair struct {
	Account passtor.Account
	Count   int
}

// MostRepresented returns the most represented verified (in the sense of signature equality)
func MostRepresented(accounts []passtor.Account, min int) (*passtor.Account, bool) {

	verified := make([]passtor.Account, 0)
	for _, account := range accounts {
		if account.Verify() {
			verified = append(verified, account)
		}
	}

	if len(verified) == 0 {
		return nil, false
	}

	signatureCounts := make(map[passtor.Signature]accountCountPair)
	for _, account := range verified {
		if count, alreadyExists := signatureCounts[account.Signature]; alreadyExists {
			signatureCounts[account.Signature] = accountCountPair{Account: count.Account, Count: count.Count + 1}
		} else {
			signatureCounts[account.Signature] = accountCountPair{Account: account, Count: 1}
		}
	}

	var mostRepresentedAccount passtor.Account
	mostRepresentedOccurences := 0
	for _, count := range signatureCounts {
		if count.Count > mostRepresentedOccurences {
			mostRepresentedOccurences = count.Count
			mostRepresentedAccount = count.Account
		}
	}

	threshIsMet := mostRepresentedOccurences >= min

	return &mostRepresentedAccount, threshIsMet

}
