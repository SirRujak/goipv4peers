package goipv4peers

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Peer struct {
	id   *[]byte
	host string
	port uint16
}

type Peers struct {
	IDLength  *int
	EntrySize *int
}

func IDLength(idLength int) (*Peers, error) {
	if idLength <= 0 {
		return nil, errors.New("peers idlength must be greater than zero")
	}
	return &Peers{IDLength: &idLength}, nil
}

func (ipv4peers Peers) EncodingLength(peers *[]Peer) (int, error) {
	var entrySize int
	if ipv4peers.IDLength != nil {
		entrySize = *ipv4peers.IDLength + 6
	} else {
		entrySize = 6
	}
	if *peers != nil {
		var length int
		length = len(*peers) * entrySize
		return length, nil
	}
	return 0, errors.New("unable to find encoding length, peer list is nil")

}

func (ipv4peers Peers) Encode(peersBase *[]Peer, bufBase *[]byte, offsetBase *int) (*[]byte, error) {
	var peers []Peer
	peers = *peersBase
	var buf []byte
	var err error
	var offset int
	if offsetBase != nil {
		offset = *offsetBase
	} else {
		offset = 0
	}

	if bufBase != nil {
		buf = *bufBase
	} else {
		var length int
		length, err = ipv4peers.EncodingLength(peersBase)
		if err != nil {
			return nil, err
		}

		buf = make([]byte, offset+length, offset+length)
	}

	var n int
	for i := 0; i < len(peers); i++ {
		if ipv4peers.IDLength != nil {
			copy(buf[offset:offset+*ipv4peers.IDLength], (*peers[i].id)[:])
			if err != nil {
				return nil, fmt.Errorf("error writing id to buffer in peers. attempted to write %v bytes", n)
			}
			offset += *ipv4peers.IDLength
		}

		// Now deal with the host and port information.
		var host []string
		host = strings.Split(peers[i].host, ".")
		var port uint16
		port = peers[i].port
		var host0, host1, host2, host3 int
		host0, err = strconv.Atoi(host[0])
		if err != nil {
			return nil, fmt.Errorf("error converting first element of host (%v) to byte", host[0])
		}
		buf[offset] = byte(host0)
		offset++

		host1, err = strconv.Atoi(host[1])
		if err != nil {
			return nil, fmt.Errorf("error converting second element of host (%v) to byte", host[1])
		}
		buf[offset] = byte(host1)
		offset++

		host2, err = strconv.Atoi(host[2])
		if err != nil {
			return nil, fmt.Errorf("error converting third element of host (%v) to byte", host[2])
		}
		buf[offset] = byte(host2)
		offset++

		host3, err = strconv.Atoi(host[3])
		if err != nil {
			return nil, fmt.Errorf("error converting fourth element of host (%v) to byte", host[3])
		}
		buf[offset] = byte(host3)
		offset++
		var tempByteArray []byte
		tempByteArray = make([]byte, 2)
		binary.BigEndian.PutUint16(tempByteArray, port)
		copy(buf[offset:offset+2], tempByteArray[:])
		offset += 2
	}

	return &buf, nil
}

func (ipv4peers Peers) Decode(buf *[]byte, offsetBase *int, endBase *int) (*[]Peer, error) {
	var offset int
	if offsetBase == nil {
		offset = 0
	} else {
		offset = *offsetBase
	}

	var end int
	if endBase == nil {
		end = len(*buf)
	} else {
		end = *endBase
	}

	var entrySize int
	if ipv4peers.IDLength != nil {
		entrySize = *ipv4peers.IDLength + 6
	} else {
		entrySize = 6
	}

	var peers []Peer
	peers = make([]Peer, int(math.Floor((float64(end-offset) / float64(entrySize)))))

	for i := 0; i < len(peers); i++ {
		var id *[]byte
		if ipv4peers.IDLength != nil {
			var tempStruct = make([]byte, *ipv4peers.IDLength, *ipv4peers.IDLength)
			id = &tempStruct
			copy(*id, (*buf)[offset:offset+*ipv4peers.IDLength])
			offset += *ipv4peers.IDLength
		} else {
			id = nil
		}
		var host string
		host = strconv.Itoa(int((*buf)[offset])) + "." + strconv.Itoa(int((*buf)[offset+1])) + "." + strconv.Itoa(int((*buf)[offset+2])) + "." + strconv.Itoa(int((*buf)[offset+3]))
		offset += 4
		var port uint16
		//port = buf.Next(1)
		port = binary.BigEndian.Uint16((*buf)[offset : offset+2])
		if port == 0 {
			return nil, errors.New("port should be > 0 and < 65536")
		}
		var tempPeer Peer
		tempPeer = Peer{
			id:   id,
			host: host,
			port: port,
		}
		peers[i] = tempPeer
		offset += 2
	}

	return &peers, nil
}
