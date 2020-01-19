package passtor

import (
	"net"

	"go.dedis.ch/protobuf"
)

// SendMessage send the given message to the remote peer over udp
// returns the reply message once it has arrived
func (p *Passtor) SendMessage(msg Message, dst net.UDPAddr,
	retries int) *Message {

	msg.Sender = &p.Addr

	if msg.Reply {
		bytes, err := protobuf.Encode(&msg)
		checkErr(err)

		_, err = p.PConn.WriteTo(bytes, &dst)
		checkErr(err)
	} else {
		// this message is a request generate new message id

		// register message channel to get the reply
		id := p.Messages.GetMessageID()
		c := make(chan Message)

		// take new message id
		p.Messages.Mutex.Lock()
		p.Messages.PendingMsg[id] = &c
		p.Messages.Mutex.Unlock()

		msg.ID = id

		bytes, err := protobuf.Encode(&msg)
		checkErr(err)

		// tries to send a message up to retries times
		for i := 0; i < retries; i++ {
			// send the message
			_, err = p.PConn.WriteTo(bytes, &dst)
			checkErr(err)

			// define a timeout
			timeout := Timeout(TIMEOUT)

			select {
			case <-*timeout:
				// on timeout resend message
				continue
			case rep := <-c:
				// on message reply return it
				return &rep
			}
		}

	}
	// if this is reached, no reply was received
	return nil
}

// HandleMessage handles incoming messages
func (p *Passtor) HandleMessage(protobufed []byte) {

	// unprotobuf the message
	rep := Message{}
	err := protobuf.Decode(protobufed, &rep)
	checkErr(err)

	// add the sender to bucket if necessary
	p.AddPeerToBucket(*rep.Sender)

	if rep.Reply {
		// if the message is a replyd istribute reply to thread that sent req
		*p.Messages.PendingMsg[rep.ID] <- rep
		return
	}

	if rep.Ping != nil {
		//bootstrap message
		rep.Reply = true
		p.SendMessage(rep, rep.Sender.Addr, MINRETRIES)
	}
}

// ListenToPasstors listen on the udp connection used to communicate with other
// passtors, and distribute received messages to HandleMessage()
func (p *Passtor) ListenToPasstors() {
	buf := make([]byte, BUFFERSIZE)

	for {
		// read new message
		m, _, err := p.PConn.ReadFromUDP(buf)
		checkErr(err)

		// copy the receive buffer to avoid that it is modified while being used
		tmp := make([]byte, m)
		copy(tmp, buf[:m])

		go p.HandleMessage(tmp)
	}
}
