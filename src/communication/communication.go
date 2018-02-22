package communication

import (
	"github.com/piot/briskd-go/src/connection"
	"github.com/piot/brook-go/src/instream"
)

type Connection interface {
	HandleStream(stream *instream.InStream) error
}

type Server interface {
	CreateConnection(id connection.ID) Connection
}
