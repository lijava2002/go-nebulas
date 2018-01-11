// Copyright (C) 2017 go-nebulas authors
//
// This file is part of the go-nebulas library.
//
// the go-nebulas library is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// the go-nebulas library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with the go-nebulas library.  If not, see <http://www.gnu.org/licenses/>.
//

package p2p

import (
	"github.com/nebulasio/go-nebulas/net"
)

// NetService service for nebulas p2p network
type NetService struct {
	node       *Node
	quitCh     chan bool
	dispatcher *net.Dispatcher
}

// NewNetService create netService
func NewNetService(n Neblet) (*NetService, error) {
	node, err := NewNode(NewP2PConfig(n))
	if err != nil {
		return nil, err
	}

	ns := &NetService{
		node:       node,
		quitCh:     make(chan bool, 10),
		dispatcher: net.NewDispatcher(),
	}
	node.SetNetService(ns)

	return ns, nil
}

// Node return the peer node
func (ns *NetService) Node() *Node {
	return ns.node
}

// Start start p2p manager.
func (ns *NetService) Start() error {
	// start dispatcher.
	ns.dispatcher.Start()

	// start node.
	if err := ns.node.Start(); err != nil {
		ns.dispatcher.Stop()
		return err
	}

	return nil
}

// Stop stop p2p manager.
func (ns *NetService) Stop() {
	ns.dispatcher.Stop()
	ns.quitCh <- true // TODO: @robin quitCh is used by multi goroutine.
}

// Register register the subscribers.
func (ns *NetService) Register(subscribers ...*net.Subscriber) {
	ns.dispatcher.Register(subscribers...)
}

// Deregister Deregister the subscribers.
func (ns *NetService) Deregister(subscribers ...*net.Subscriber) {
	ns.dispatcher.Deregister(subscribers...)
}

// PutMessage put message to dispatcher.
func (ns *NetService) PutMessage(msg net.Message) {
	ns.dispatcher.PutMessage(msg)
}

// Broadcast message.
func (ns *NetService) Broadcast(name string, msg net.Serializable, priority int) {
	ns.node.BroadcastMessage(name, msg, priority)
}

// Relay message.
func (ns *NetService) Relay(name string, msg net.Serializable, priority int) {
	ns.node.RelayMessage(name, msg, priority)
}

// BroadcastNetworkID broadcast networkID when changed.
func (ns *NetService) BroadcastNetworkID(msg []byte) {
	// TODO: @robin networkID.
}

// BuildRawMessageData return the raw NebMessage content data.
func (ns *NetService) BuildRawMessageData(data []byte, msgName string) []byte {
	message, err := NewNebMessage(ns.node.config.ChainID, DefaultReserved, 0, msgName, data)
	if err != nil {
		return nil
	}

	return message.Content()
}

// SendMsg send message to a peer.
func (ns *NetService) SendMsg(msgName string, msg []byte, target string, priority int) error {
	return ns.node.SendMessageToPeer(target, msgName, msg, priority)
}
