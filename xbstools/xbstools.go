package xbstools

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/yang3yen/xxtea-go/xxtea"
)

var xxTeaKey = []byte{0xe5, 0x87, 0xbc, 0xe8, 0xa4, 0x86, 0xe6, 0xbb, 0xbf, 0xe9, 0x87, 0x91, 0xe6, 0xba, 0xa1, 0xe5}

func XBS2Json(buffer []byte) ([]byte, error) {

	out, err := xxtea.Decrypt(buffer, xxTeaKey, false, 0)
	if err != nil {
		return nil, err
	}
	var n uint32 = uint32(len(buffer))
	n = n - 4
	m := binary.LittleEndian.Uint32(out[n:])
	if m < n-3 || m > n {
		return nil, fmt.Errorf("decode error")
	}
	n = m
	return out[:n], nil
}

func Json2XBS(buffer []byte) ([]byte, error) {
	var buffer_len uint32 = uint32(len(buffer))
	var n uint32 = 0
	var buffer_enc_len []byte
	if (buffer_len & 3) == 0 {
		n = buffer_len >> 2
	} else {
		n = (buffer_len >> 2) + 1
	}
	for i := buffer_len; i < (n << 2); i++ {
		buffer_enc_len = append(buffer_enc_len, 0x0)
	}

	buffer_enc_len = binary.LittleEndian.AppendUint32(buffer_enc_len, buffer_len)

	buffer = append(buffer, buffer_enc_len...)
	out, err := xxtea.Encrypt(buffer, xxTeaKey, false, 0)
	return out, err
}

func LoadFile(filepath string) ([]byte, error) {
	if _, err := os.Stat(filepath); err == nil {
		return os.ReadFile(filepath)
	} else {
		return nil, err
	}
}
