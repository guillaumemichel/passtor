package main

import (
	"crypto/sha256"
)

// SetIdentity set the identity of the passtor instance to the hash of the given
// name
func (p *Passtor) SetIdentity() {
	p.NodeID = sha256.Sum256([]byte(p.Name))
}
