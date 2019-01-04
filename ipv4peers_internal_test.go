package ipv4peers

import (
	"bytes"
	"fmt"
	"testing"
)

func compareStruct(peer1 Peer, peer2 Peer) bool {
	if peer1.host == peer2.host && peer1.port == peer2.port {
		if peer1.id == nil && peer2.id == nil {
			return true
		} else if bytes.Compare(*peer1.id, *peer2.id) == 0 {
			return true
		}
	}
	return false
}

func comparePeers(peers1 *[]Peer, peers2 *[]Peer) bool {
	if len(*peers1) != len(*peers2) {
		return false
	}
	for i := 0; i < len(*peers1); i++ {
		if !compareStruct((*peers1)[i], (*peers2)[i]) {
			return false
		}
	}
	return true
}

func TestEncode(t *testing.T) {
	peers := Peers{}
	tempPeers := []Peer{
		Peer{
			host: "127.0.0.1",
			port: 80,
		},
		Peer{
			host: "127.0.0.1",
			port: 8080,
		},
	}
	n, err := peers.Encode(&tempPeers, nil, nil)
	if err != nil {
		t.Error(err)
	}
	if len(*n) != 12 {
		t.Errorf("Length was not %v, instead got %v", 12, len(*n))
	}

	tempPeers = []Peer{
		Peer{
			host: "127.0.0.1",
			port: 80,
		},
	}
	n, err = peers.Encode(&tempPeers, nil, nil)
	if err != nil {
		t.Error(err)
	}
	if len(*n) != 6 {
		t.Errorf("Length was not %v, instead got %v", 6, len(*n))
	}
}

func TestEncodingLength(t *testing.T) {
	peers := Peers{}
	tempPeers := []Peer{
		Peer{
			host: "127.0.0.1",
			port: 80,
		},
		Peer{
			host: "127.0.0.1",
			port: 8080,
		},
	}
	n, err := peers.EncodingLength(&tempPeers)
	if err != nil {
		t.Error(err)
	}
	if n != 12 {
		t.Errorf("Length was not %v, instead got %v", 12, n)
	}

	tempPeers = []Peer{
		Peer{
			host: "127.0.0.1",
			port: 80,
		},
	}
	n, err = peers.EncodingLength(&tempPeers)
	if err != nil {
		t.Error(err)
	}
	if n != 6 {
		t.Errorf("Length was not %v, instead got %v", 6, n)
	}

	tempPeers = nil
	n, err = peers.EncodingLength(&tempPeers)
	if err == nil {
		t.Error(err)
	}
	if n != 0 {
		t.Errorf("Length was not %v, instead got %v", 6, n)
	}
}

func TestEncodeDecode(t *testing.T) {
	peers := Peers{}
	tempPeers := []Peer{
		Peer{
			host: "127.0.0.1",
			port: 80,
		},
		Peer{
			host: "127.0.0.1",
			port: 8080,
		},
	}
	n, err := peers.Encode(&tempPeers, nil, nil)
	if err != nil {
		t.Error(err)
	} else {
		var newPeers *[]Peer
		newPeers, err = peers.Decode(n, nil, nil)
		if err != nil {
			t.Error(err)
		} else {
			if !comparePeers(&tempPeers, newPeers) {
				t.Errorf("peers do not match when encoded then decoded")
				fmt.Printf("%v != %v\n", tempPeers, newPeers)
			}
		}

	}

	tempPeers = []Peer{
		Peer{
			host: "127.0.0.1",
			port: 80,
		},
	}
	n, err = peers.Encode(&tempPeers, nil, nil)

	if err != nil {
		t.Error(err)
	} else {
		var newPeers *[]Peer
		newPeers, err = peers.Decode(n, nil, nil)
		if err != nil {
			t.Error(err)
		} else {
			if !comparePeers(&tempPeers, newPeers) {
				t.Errorf("peers do not match when encoded then decoded")
			}
		}

	}
}

func TestEncodeDecodeOffset(t *testing.T) {
	peers := Peers{}
	tempPeers := []Peer{
		Peer{
			host: "127.0.0.1",
			port: 80,
		},
		Peer{
			host: "127.0.0.1",
			port: 8080,
		},
	}
	offset := 2
	length := 12
	var tempByteList []byte
	tempByteList = make([]byte, offset+length, offset+length)
	n, err := peers.Encode(&tempPeers, &tempByteList, &offset)
	if err != nil {
		t.Error(err)
	} else {
		var newPeers *[]Peer
		newPeers, err = peers.Decode(n, &offset, nil)
		if err != nil {
			t.Error(err)
		} else {
			if !comparePeers(&tempPeers, newPeers) {
				t.Errorf("peers do not match when encoded then decoded")
			}
		}

	}

	tempPeers = []Peer{
		Peer{
			host: "127.0.0.1",
			port: 80,
		},
	}
	length = 6
	tempByteList = make([]byte, offset+length, offset+length)
	n, err = peers.Encode(&tempPeers, &tempByteList, &offset)

	if err != nil {
		t.Error(err)
	} else {
		var newPeers *[]Peer
		newPeers, err = peers.Decode(n, &offset, nil)
		if err != nil {
			t.Error(err)
		} else {
			if !comparePeers(&tempPeers, newPeers) {
				t.Errorf("peers do not match when encoded then decoded")
			}
		}

	}
}

func TestEncodeDecodeOffsetEnd(t *testing.T) {
	peers := Peers{}
	tempPeers := []Peer{
		Peer{
			host: "127.0.0.1",
			port: 80,
		},
		Peer{
			host: "127.0.0.1",
			port: 8080,
		},
	}
	offset := 2
	length := 12
	var tempByteList []byte
	tempByteList = make([]byte, offset+length, offset+length)
	n, err := peers.Encode(&tempPeers, &tempByteList, &offset)
	if err != nil {
		t.Error(err)
	} else {
		var end int
		end = offset + length
		var newPeers *[]Peer
		newPeers, err = peers.Decode(n, &offset, &end)
		if err != nil {
			t.Error(err)
		} else {
			if !comparePeers(&tempPeers, newPeers) {
				t.Errorf("peers do not match when encoded then decoded")
			}
		}

	}

	tempPeers = []Peer{
		Peer{
			host: "127.0.0.1",
			port: 80,
		},
	}
	length = 6
	tempByteList = make([]byte, offset+length, offset+length)
	n, err = peers.Encode(&tempPeers, &tempByteList, &offset)

	if err != nil {
		t.Error(err)
	} else {
		var end int
		end = offset + length
		var newPeers *[]Peer
		newPeers, err = peers.Decode(n, &offset, &end)
		if err != nil {
			t.Error(err)
		} else {
			if !comparePeers(&tempPeers, newPeers) {
				t.Errorf("peers do not match when encoded then decoded")
			}
		}

	}
}

func TestPortZeroNotAllowed(t *testing.T) {
	peers := Peers{}
	tempPeers := []Peer{
		Peer{
			host: "127.0.0.1",
			port: 0,
		},
	}
	n, err := peers.Encode(&tempPeers, nil, nil)

	if err != nil {
		t.Error(err)
	} else {
		_, err := peers.Decode(n, nil, nil)
		if err == nil {
			t.Errorf("decoder does not catch disallowed zero port")
		}

	}
}

func TestEncodeDecodeOffsetEndId(t *testing.T) {
	var idLength = 5
	peers, err := IDLength(idLength)
	if err != nil {
		t.Error(err)
	}
	var id1 = []byte("hello")
	var id2 = []byte("hiaya")
	tempPeers := []Peer{
		Peer{
			id:   &id1,
			host: "127.0.0.1",
			port: 80,
		},
		Peer{
			id:   &id2,
			host: "127.0.0.1",
			port: 8080,
		},
	}
	offset := 0
	length := 22
	var tempByteList []byte
	tempByteList = make([]byte, offset+length, offset+length)
	n, err := peers.Encode(&tempPeers, &tempByteList, &offset)
	if err != nil {
		t.Error(err)
	} else {
		var end int
		end = offset + length
		var newPeers *[]Peer
		newPeers, err = peers.Decode(n, &offset, &end)
		if err != nil {
			t.Error(err)
		} else {
			if !comparePeers(&tempPeers, newPeers) {
				t.Errorf("peers do not match when encoded then decoded")
			}
		}

	}

	tempPeers = []Peer{
		Peer{
			id:   &id1,
			host: "127.0.0.1",
			port: 80,
		},
	}
	offset = 0
	length = 11
	tempByteList = make([]byte, offset+length, offset+length)
	n, err = peers.Encode(&tempPeers, &tempByteList, &offset)

	if err != nil {
		t.Error(err)
	} else {
		var end int
		end = offset + length
		var newPeers *[]Peer
		newPeers, err = peers.Decode(n, &offset, &end)
		if err != nil {
			t.Error(err)
		} else {
			if !comparePeers(&tempPeers, newPeers) {
				t.Errorf("peers do not match when encoded then decoded")
			}
		}

	}
}
