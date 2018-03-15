/*

MIT License

Copyright (c) 2017 Peter Bjorklund

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/

package server

import (
	"fmt"

	"github.com/piot/briskd-go/src/commandcreator"
	"github.com/piot/briskd-go/src/commands"
	"github.com/piot/briskd-go/src/communication"
	"github.com/piot/briskd-go/src/connection"
	"github.com/piot/briskd-go/src/endpoint"
	"github.com/piot/briskd-go/src/message"
	"github.com/piot/briskd-go/src/packet"
	"github.com/piot/briskd-go/src/sequence"
	brisktime "github.com/piot/briskd-go/src/time"
	"github.com/piot/brook-go/src/instream"
	"github.com/piot/brook-go/src/outstream"

	"net"
	"time"
)

type Server struct {
	connection                     *net.UDPConn
	connections                    map[connection.ID]*Connection
	userServer                     communication.Server
	lastAllocatedConnectionIDValue uint16
	lastTimeStatsCalculatedAt      int64
}

func (s *Server) SendPacketToConnection(conn *Connection, stream *outstream.OutStream) {
	addr := conn.Addr()
	octetCount := uint(len(stream.Octets()))
	conn.SentPacket(octetCount)
	s.SendPacketToEndpoint(addr, stream)
}

func (s *Server) findConnection(addr *endpoint.Endpoint, connectionID connection.ID) (*Connection, error) {
	connection, foundIt := s.connections[connectionID]
	if !foundIt {
		return nil, nil
	}
	return connection, nil
}

func (s *Server) fetchConnection(addr *endpoint.Endpoint, connectionID connection.ID) (*Connection, error) {
	connection, err := s.findConnection(addr, connectionID)
	if err != nil {
		return connection, err
	}
	if connection == nil {
		return nil, fmt.Errorf("couldn't find connection %d", connectionID)
	}
	return connection, nil
}

func (s *Server) onTimeSync(addr *endpoint.Endpoint, timesyncRequest *commands.TimeSyncRequest, connection *Connection) error {
	fmt.Printf("on_timesync: %v\n", timesyncRequest)
	localTime := uint64(brisktime.MonotonicMilliseconds())
	response := commands.NewTimeSyncResponse(timesyncRequest.RemoteTime, localTime)
	s.SendMessageToConnection(connection, response, packet.OobMode)

	return nil
}

func (s *Server) challenge(addr *endpoint.Endpoint, challengeMessage *commands.ChallengeMessage) error {
	fmt.Printf("on_challenge:%s\n", challengeMessage)
	existingConnection := s.findExistingConnectionFromEndpointAndChallenge(addr, challengeMessage.ClientDeterminedNonce)
	if existingConnection == nil {
		newConnection, err := s.createConnection(addr, challengeMessage.ClientDeterminedNonce)
		if err != nil {
			return err
		}
		userConnection := s.userServer.CreateConnection(newConnection.ID())
		newConnection.SetUserConnection(userConnection)
		response := commands.NewChallengeResponseMessage(challengeMessage.ClientDeterminedNonce, newConnection.ID())
		s.SendMessageToEndpoint(addr, response)
	} else {
		return fmt.Errorf("We already have a connection for nonce: %d", challengeMessage.ClientDeterminedNonce)
	}
	return nil
}

func (s *Server) handleOOBMessage(addr *endpoint.Endpoint, msg message.Message, connection *Connection) error {
	fmt.Printf("OOB %s\n ", msg)

	switch msg.Command() {
	case packet.OobPacketTypeChallenge:
		err := s.challenge(addr, msg.(*commands.ChallengeMessage))
		if err != nil {
			return err
		}
	case packet.OobPacketTypeTimeSyncRequest:
		err := s.onTimeSync(addr, msg.(*commands.TimeSyncRequest), connection)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("not handled %d", msg.Command())
	}
	return nil

}

func (s *Server) handleOOBStream(addr *endpoint.Endpoint, inStream *instream.InStream, connection *Connection) error {
	msg, msgErr := commandcreator.CreateMessage(inStream)
	if msgErr != nil {
		return msgErr
	}
	err := s.handleOOBMessage(addr, msg, connection)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) handlePacket(buf []byte, addr *endpoint.Endpoint) error {
	inStream := instream.New(buf)
	header, _ := packet.ReadHeader(&inStream)

	if header.Mode != packet.NormalMode && header.Mode != packet.OobMode {
		return fmt.Errorf("unknown mode")
	}

	for !inStream.IsEOF() {
		if header.ConnectionID == 0 {
			s.handleOOBStream(addr, &inStream, nil)
		} else {
			connection, findConnectionErr := s.fetchConnection(addr, header.ConnectionID)
			if findConnectionErr != nil {
				return findConnectionErr
			}
			connection.LastReceivedPacketAt = brisktime.MonotonicMilliseconds()
			if header.Mode == packet.OobMode {
				handleErr := s.handleOOBStream(addr, &inStream, connection)
				if handleErr != nil {
					return handleErr
				}
			} else {
				connection.ReceivedPacket(uint(len(buf)))
				connection.DebugIncomingPacket(buf[inStream.Tell():], brisktime.MonotonicMilliseconds())
				handleErr := connection.handleStream(&inStream)
				if handleErr != nil {
					return handleErr
				}
			}
			break
		}
	}

	return nil
}

// New : Creates a new server
func New(userServer communication.Server) Server {
	return Server{connections: make(map[connection.ID]*Connection), userServer: userServer}
}

func (s *Server) handleIncomingUDP() {
	for {
		buf := make([]byte, 1800)
		n, addr, err := s.connection.ReadFromUDP(buf)
		packet := buf[0:n]
		addrEndpoint := endpoint.New(addr)
		//hexPayload := hex.Dump(packet)
		//fmt.Println("Received ", hexPayload, " from ", addr)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		packetErr := s.handlePacket(packet, &addrEndpoint)
		if packetErr != nil {
			fmt.Printf("Problem with packet:%s\n", packetErr)
		}
	}
}

// SendPacketToEndpoint : Sends one packet to endpoint without rate limit
func (s *Server) SendPacketToEndpoint(addr *endpoint.Endpoint, stream *outstream.OutStream) {
	octets := stream.Octets()
	// hexPayload := hex.Dump(octets)
	//fmt.Println("Sending ", hexPayload, " to ", addr)
	s.connection.WriteToUDP(octets, addr.UDPAddr())
}

func headerAndMessageToStream(header *packet.PacketHeader, message2 message.Message) *outstream.OutStream {
	stream := outstream.New()

	packet.WriteHeader(stream, header)
	stream.WriteUint8(uint8(message2.Command()))
	message2.Serialize(stream)

	return stream
}

func (s *Server) SendMessageToEndpoint(addr *endpoint.Endpoint, message2 message.Message) {
	emptySequenceID, _ := sequence.NewID(sequence.IDType(0))
	header := packet.PacketHeader{Mode: packet.NormalMode, Sequence: emptySequenceID, ConnectionID: connection.ID(0)}
	stream := headerAndMessageToStream(&header, message2)

	// fmt.Printf(">> %s %s\n", addr, message2)
	s.SendPacketToEndpoint(addr, stream)
}

func (s *Server) SendMessageToConnection(connection *Connection, message2 message.Message, mode packet.Mode) error {
	stream := writeConnectionHeader(connection, mode)
	stream.WriteUint8(uint8(message2.Command()))
	message2.Serialize(stream)
	// fmt.Printf(">>> %v %v\n", connection, message2)
	s.SendPacketToConnection(connection, stream)
	return nil
}

func (s *Server) tick() error {
	s.userServer.Tick()
	var resultErr error
	for _, connection := range s.connections {
		for i := 0; i < 1; i++ {
			err := s.sendStream(connection)
			if err != nil {
				if resultErr != nil {
					resultErr = err
				}
			}
		}
	}

	nowMilliseconds := brisktime.MonotonicMilliseconds()
	var deletedConnections []*Connection
	for _, connection := range s.connections {
		msSinceReceived := nowMilliseconds - connection.LastReceivedPacketAt
		if msSinceReceived >= 3000 {
			connection.Lost()
			deletedConnections = append(deletedConnections, connection)
		}
	}

	timeSinceStats := nowMilliseconds - s.lastTimeStatsCalculatedAt
	if timeSinceStats > 5000 {
		for _, connection := range s.connections {
			connection.CalculateStats(uint(timeSinceStats))
		}
		s.lastTimeStatsCalculatedAt = nowMilliseconds
	}

	for _, deletedConnection := range deletedConnections {
		delete(s.connections, deletedConnection.ID())
	}

	if resultErr != nil {
		fmt.Printf("Error: %s\n", resultErr)
	}
	return resultErr
}

func writeConnectionHeader(connection *Connection, mode packet.Mode) *outstream.OutStream {
	stream := outstream.New()
	connection.NextOutSequenceID = connection.NextOutSequenceID.Next()
	header := &packet.PacketHeader{Mode: mode, Sequence: connection.NextOutSequenceID, ConnectionID: connection.ID()}
	packet.WriteHeader(stream, header)
	return stream
}

func (s *Server) sendStream(connection *Connection) error {
	stream := writeConnectionHeader(connection, packet.NormalMode)
	startPosition := stream.Tell()
	userErr := connection.userConnection.SendStream(stream)
	if userErr != nil {
		return userErr
	}
	connection.DebugOutgoingPacket(stream.Octets()[startPosition:], brisktime.MonotonicMilliseconds())

	s.SendPacketToConnection(connection, stream)
	return nil
}

func (s *Server) start(ticker *time.Ticker) {
	go func() {
		for range ticker.C {
			err := s.tick()
			if err != nil {
				fmt.Printf("Start err %s \n", err)
			}
		}
	}()
}

func (s *Server) Forever() error {
	const portString = ":32001"
	serverAddress, err := net.ResolveUDPAddr("udp", portString)
	if err != nil {
		return fmt.Errorf("Error:%v ", err)
	}
	serverConnection, err := net.ListenUDP("udp", serverAddress)
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}

	fmt.Printf("Listening to %s\n", portString)

	go s.handleIncomingUDP()
	//defer serverConnection.Close()
	ticker := time.NewTicker(time.Millisecond * 33)

	s.connection = serverConnection
	s.start(ticker)

	select {}

	return nil
}

func (s *Server) findExistingConnectionFromEndpointAndChallenge(addr *endpoint.Endpoint, nonce uint32) *Connection {
	for _, connection := range s.connections {
		if connection.Addr().Equal(addr) && connection.nonce == nonce {
			return connection
		}
	}

	return nil
}

func (s *Server) findFreeConnectionID() (connection.ID, error) {
	idToCheckValue := s.lastAllocatedConnectionIDValue
	for i := 0; i < 65536; i++ {
		idToCheckValue += 61
		if idToCheckValue == 0 {
			continue
		}
		idToCheck := connection.ID(idToCheckValue)
		if _, exists := s.connections[idToCheck]; !exists {
			s.lastAllocatedConnectionIDValue = idToCheckValue
			return idToCheck, nil
		}
	}

	return 0, fmt.Errorf("no free connections")
}

func (s *Server) createConnection(endpoint *endpoint.Endpoint, nonce uint32) (*Connection, error) {
	connectionID, err := s.findFreeConnectionID()
	if err != nil {
		return nil, err
	}
	newConnection := NewConnection(s, connectionID, endpoint, nonce)
	s.connections[connectionID] = newConnection

	return newConnection, nil
}
