package gate

import (
	"kada/core"
	"kada/log"
)

const (
	MESSAGE_LENGTH int32 = 4 // 消息长度
)

//Enpack 封包
func Enpack(head int32, data []byte) []byte {
	body := append(core.Int32ToBytes(head), data...)
	return append(core.Int32ToBytes(int32(len(body))), body...)
}

//Depack 解包
func Depack(sid string, buffer []byte) []byte {
	log.Debug("Gate Package Depack --", buffer)

	length := int32(len(buffer))

	var i int32
	body := make([]byte, 1024)
	for i = 0; i < length; i = i + 1 {
		if length < i+MESSAGE_LENGTH {
			break
		}

		readLen := core.BytesToInt32(buffer[i : i+MESSAGE_LENGTH])
		if length < i+MESSAGE_LENGTH+readLen {
			break
		}

		body = buffer[i+MESSAGE_LENGTH : i+MESSAGE_LENGTH+readLen]

		head := core.BytesToInt32(body[0:4])

		// data := make([]byte, 1024)
		data := body[4:]

		Call(sid, head, data)

		i = i + MESSAGE_LENGTH + readLen - 1
	}

	if i >= length {
		return make([]byte, 0)
	}

	return body
}
