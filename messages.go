package passtor

import (
	"fmt"
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

	rep.Reply = true
	if rep.Ping != nil {
		//ping message
		p.SendMessage(rep, rep.Sender.Addr, MINRETRIES)
	} else if rep.LookupReq != nil {
		// reply to lookup
		p.LookupRep(rep)
	} else if rep.AllocationReq != nil {
		// allocate and reply to allocation request
		p.HandleAllocation(rep)
	} else if rep.FetchReq != nil {
		// searches and reply with corresponding data
		p.HandleFetch(rep)
	}
}

func (p *Passtor) HandleClientMessage(accounts Accounts, message ClientMessage) *ServerResponse {

	// Update or store new account
	if message.Push != nil {
		p.Printer.Print(fmt.Sprint("Push request ", message.Push.ID), V2)
		providers := p.Allocate(message.Push.ID, REPL, *message.Push)
		if len(providers) == REPL {
			return &ServerResponse{
				Status: "ok",
			}
		}
		msg := fmt.Sprintln("couldn't allocate data to ", REPL, "peers")
		return &ServerResponse{
			Status: "warning",
			Debug:  &msg,
		}
	}

	// Retrieve account
	if message.Pull != nil {
		p.Printer.Print(fmt.Sprint("Pull request ", *message.Pull), V2)
		account := p.FetchData(message.Pull, THRESHOLD)

		if account != nil {
			accountNetwork := account.ToAccountNetwork()
			return &ServerResponse{
				Status: "ok",
				Data:   &accountNetwork,
			}
		}

		msg := "This account does not exist"
		return &ServerResponse{
			Status: "error",
			Debug:  &msg,
		}
	}
	return nil
}
