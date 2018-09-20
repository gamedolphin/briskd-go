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

	"time"

	"github.com/piot/brisk-protocol-go/src/connection"
)

type Server struct {
	//userServer communication.Server
	server                    *connection.Server
	lastTimeStatsCalculatedAt int64
	debugEnabled              bool
}

// New : Creates a new server
func New(listenPort int, userServer connection.UserServer, enableDebug bool, dumpPackets bool) (*Server, error) {
	connectionServer, serverErr := connection.NewServer(listenPort, userServer, dumpPackets)
	if serverErr != nil {
		return nil, serverErr
	}

	//defer serverConnection.Close()

	s := &Server{debugEnabled: enableDebug, server: connectionServer}
	return s, nil
}

func (s *Server) tick() error {
	s.server.Tick()
	//s.userServer.Tick()
	/*
		if resultErr != nil {
			fmt.Printf("Error: %s\n", resultErr)
		}
		return resultErr
	*/
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
	ticker := time.NewTicker(time.Millisecond * 33)
	s.start(ticker)

	select {}

	return nil
}
