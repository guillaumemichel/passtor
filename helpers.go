package passtor

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"
	"net"
	"sort"
	"strings"
	"time"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// ParsePeers parse peer list in string format to udp addresses
func ParsePeers(peerList string) []net.UDPAddr {

	addresses := make([]net.UDPAddr, 0)

	if peerList == "" {
		return addresses
	}
	// split up the different addresses
	peers := strings.Split(peerList, ",")

	// parse the addresses and add them to the slice
	for _, p := range peers {
		udpAddr, err := net.ResolveUDPAddr("udp4", p)
		checkErrMsg(err, "invalid address \""+p+"\"")
		addresses = append(addresses, *udpAddr)
	}

	return addresses
}

// Timeout creates a clock that writes to the returned channel after the
// time value given as argument
func Timeout(timeout time.Duration) *chan bool {
	c := make(chan bool)
	go func() {
		time.Sleep(timeout)
		c <- true
	}()
	return &c
}

// NewLookupStatus returns new lookup status structure for given nodeaddr
func NewLookupStatus(nodeAddr NodeAddr) *LookupStatus {
	return &LookupStatus{
		NodeAddr: nodeAddr,
		Failed:   false,
		Tested:   false,
	}
}

// RandomBytes generates an array of random bytes of the given size
func RandomBytes(size uint) ([]byte, error) {

	bytes := make([]byte, size)

	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// BytesToNonce converts a byte array to a Nonce type.
func BytesToNonce(array []byte) Nonce {
	if len(array) != NONCESIZE {
		panic("Array is expected to have size " + string(NONCESIZE))
	}

	var nonce = Nonce{}
	for i, b := range array {
		nonce[i] = b
	}

	return nonce
}

func NonceToBytes(nonce Nonce) []byte {
	array := make([]byte, len(nonce))

	for i, b := range nonce {
		array[i] = b
	}

	return array
}

// BytesToSymmetricKey creates a symmetric key from an array of bytes
func BytesToSymmetricKey(array []byte) SymmetricKey {
	if len(array) != SYMMKEYSIZE {
		panic("Array is expected to have size " + string(SYMMKEYSIZE))
	}

	var symmKey = [SYMMKEYSIZE]byte{}
	for i, b := range array {
		symmKey[i] = b
	}

	return symmKey
}

// SymmetricKeyToBytes converts a symmetric key to a raw array of bytes
func SymmetricKeyToBytes(symmK SymmetricKey) []byte {

	array := make([]byte, len(symmK))

	for i, b := range symmK {
		array[i] = b
	}

	return array
}

func BytesToSignature(array []byte) Signature {
	if len(array) != SIGNATURESIZE {
		panic("Array is expected to have size " + string(SIGNATURESIZE))
	}

	var sig = Signature{}
	for i, b := range array {
		sig[i] = b
	}

	return sig
}

func SignatureToBytes(signature Signature) []byte {
	array := make([]byte, len(signature))

	for i, b := range signature {
		array[i] = b
	}

	return array
}

func KDFToSecret(array []byte) Secret {
	if len(array) != SECRETLENGTH {
		panic("Array is expected to have size " + string(SECRETLENGTH))
	}

	var secret = Secret{}
	for i, b := range array {
		secret[i] = b
	}

	return secret
}

func HashToBytes(h Hash) []byte {
	array := make([]byte, len(h))

	for i, b := range h {
		array[i] = b
	}

	return array
}

func BytesToHash(array []byte) Hash {
	if len(array) != HASHSIZE {
		panic("Array is expected to have size " + string(HASHSIZE))
	}

	var h = Hash{}
	for i, b := range array {
		h[i] = b
	}

	return h
}

func GetKeysSorted(data map[Hash]Login) []Hash {
	keysString := make([]string, len(data))

	i := 0
	for k := range data {
		keysString[i] = base64.StdEncoding.EncodeToString(HashToBytes(k))
		i++
	}
	sort.Strings(keysString)

	keysHash := make([]Hash, len(data))

	i = 0
	for _, s := range keysString {
		h, err := base64.StdEncoding.DecodeString(s)
		if err != nil {
			panic("base64 decoding failed")
		}
		keysHash[i] = BytesToHash(h)
		i++
	}

	return keysHash
}

func DuplicateMap(data map[Hash]Login) map[Hash]Login {
	newMap := make(map[Hash]Login, len(data))

	for k, v := range data {
		newMap[k] = v
	}

	return newMap
}

// MostRepresented returns the most represented verified (in the sense of signature equality)
func MostRepresented(accounts []Account, min int) (*Account, bool) {

	verified := make([]Account, 0)
	for _, account := range accounts {
		if account.Verify() {
			verified = append(verified, account)
		}
	}

	if len(verified) == 0 {
		return nil, false
	}

	signatureCounts := make(map[Signature]accountCountPair)
	for _, account := range verified {
		if count, alreadyExists := signatureCounts[account.Signature]; alreadyExists {
			signatureCounts[account.Signature] = accountCountPair{Account: count.Account, Count: count.Count + 1}
		} else {
			signatureCounts[account.Signature] = accountCountPair{Account: account, Count: 1}
		}
	}

	var mostRepresentedAccount Account
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

// RandInt generate a random int64 between 0 and given n
func RandInt(n int64) int64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(n))
	if err != nil {
		panic(err)
	}
	return nBig.Int64()
}
