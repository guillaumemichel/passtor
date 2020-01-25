package passtor

<<<<<<< HEAD
=======
import "gitlab.gnugen.ch/gmichel/passtor/crypto"

>>>>>>> master
// SetIdentity set the identity of the passtor instance to the hash of the given
// name
func (p *Passtor) SetIdentity() {
	p.NodeID = H([]byte(p.Name))
}
