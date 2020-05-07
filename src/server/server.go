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

	"github.com/piot/chrono-go/src/chrono"

	"github.com/piot/brisk-protocol-go/src/connection"
	"github.com/piot/brisk-protocol-go/src/meta"
	"github.com/piot/log-go/src/clog"
)

type Server struct {
	//userServer communication.Server
	server                    *connection.Server
	lastTimeStatsCalculatedAt int64
	debugEnabled              bool
	frequency                 int
	log                       *clog.Log
}

// New : Creates a new server
func New(listenPort int, userServer connection.UserServer, updateFrequency int, log *clog.Log, dumpPackets bool, disconnectTimoutMs uint, header meta.Header, schemaPayload []byte, findFreePort bool) (*Server, int, error) {
	connectionServer, foundPort, serverErr := connection.NewServer(listenPort, userServer, dumpPackets, disconnectTimoutMs, header, schemaPayload, findFreePort, log)
	if serverErr != nil {
		return nil, 0, serverErr
	}
	if updateFrequency == 0 {
		return nil, 0, fmt.Errorf("illegal update frequency")
	}

	if connectionServer == nil {
		return nil, 0, fmt.Errorf("must have valid connection server")
	}

	//defer serverConnection.Close()

	s := &Server{log: log, server: connectionServer, frequency: updateFrequency}
	return s, foundPort, nil
}

func (s *Server) tick() error {
	s.server.Tick()
	return nil
}

func (s *Server) updateFn() bool {
	if s == nil {
		return false
	}
	err := s.tick()
	if err != nil {
		s.log.Error("start error", clog.Error("start error", err))
	}
	return true
}

func (s *Server) Forever() error {
	updater, updaterErr := chrono.NewUpdater(s.frequency, s.updateFn)
	if updaterErr != nil {
		return updaterErr
	}
	updater.Start()

	select {}

	return nil
}
