package mess

import (
	"encoding/binary"
	"errors"
)

var ErrMessageInvalid = errors.New("invalid message")

type Message struct {
	Header []byte
	Body   []byte
}

func (m *Message) Encode() []byte {
	buf := make([]byte, 4, len(m.Header)+len(m.Body)+8)
	binary.BigEndian.PutUint32(buf, uint32(len(m.Header)))

	buf = append(buf, m.Header...)
	binary.BigEndian.PutUint32(buf[4+len(m.Header):8+len(m.Header)], uint32(len(m.Body)))

	buf = append(buf, m.Body...)

	return buf
}

func (m *Message) Dncode(data []byte) error {
	if len(data) < 4 {
		return ErrMessageInvalid
	}
	headerLen := int(binary.BigEndian.Uint32(data[0:4]))
	if len(data) < 4+headerLen {
		return ErrMessageInvalid
	}
	header := make([]byte, headerLen, 0)
	copy(header, data[4:4+headerLen])
	m.Header = header

	if len(data) == 4+headerLen {
		m.Body = make([]byte, 0)
		return nil
	}

	if len(data) < 8+headerLen {
		return ErrMessageInvalid
	}

	bodyLen := int(binary.BigEndian.Uint32(data[4+headerLen : 8+headerLen]))
	if len(data) < 8+headerLen+bodyLen {
		return ErrMessageInvalid
	}

	body := make([]byte, bodyLen, 0)
	copy(header, data[8+headerLen:8+headerLen+bodyLen])
	m.Body = body

	return nil
}
