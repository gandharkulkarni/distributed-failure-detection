package msg_handler

import (
	"encoding/binary"
	"net"

	"google.golang.org/protobuf/proto"
)

type MsgHandler struct {
	conn net.Conn
}

func NewMsgHandler(conn net.Conn) *MsgHandler {
	m := &MsgHandler{
		conn: conn,
	}
	return m
}

func (m *MsgHandler) readN(buf []byte) error {
	bytesRead := uint64(0)
	for bytesRead < uint64(len(buf)) {
		n, err := m.conn.Read(buf[bytesRead:])
		if err != nil {
			return err
		}
		bytesRead += uint64(n)
	}
	return nil
}

func (m *MsgHandler) writeN(buf []byte) error {
	bytesWritten := uint64(0)
	for bytesWritten < uint64(len(buf)) {
		n, err := m.conn.Write(buf[bytesWritten:])
		if err != nil {
			return err
		}
		bytesWritten += uint64(n)
	}
	return nil
}

func (m *MsgHandler) Send(wrapper *NodeMsg) error {
	serialized, err := proto.Marshal(wrapper)
	if err != nil {
		return err
	}

	prefix := make([]byte, 8)
	binary.LittleEndian.PutUint64(prefix, uint64(len(serialized)))
	m.writeN(prefix)
	m.writeN(serialized)

	return nil
}

func (m *MsgHandler) Receive() (*NodeMsg, error) {
	prefix := make([]byte, 8)
	m.readN(prefix)

	payloadSize := binary.LittleEndian.Uint64(prefix)
	payload := make([]byte, payloadSize)
	m.readN(payload)

	wrapper := &NodeMsg{}
	err := proto.Unmarshal(payload, wrapper)
	return wrapper, err
}

func (m *MsgHandler) Close() {
	m.conn.Close()
}
