package passtor

import "gitlab.gnugen.ch/gmichel/passtor/crypto"

// SetIdentity set the identity of the passtor instance to the hash of the given
// name
func (p *Passtor) SetIdentity() {
	p.NodeID = crypto.H(p.Name)
}
