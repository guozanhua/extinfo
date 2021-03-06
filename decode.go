package extinfo

import (
	"bytes"
	"log"
)

// the current position in a response ([]byte)
// needed, since values are encoded in variable amount of bytes
// global to not have to pass around an int on every dump
var positionInResponse int

// decodes the bytes read from the connection into ints
// returns the decoded byte slice as int and the amount of bytes used up of the slice
func getInt(buf []byte) (int, int) {
	// n is the size of the buffer
	n := len(buf)
	// b is the first byte in buf
	b := buf[0]

	// 0x80 means: value is contained in the next 2 more bytes
	if b == 0x80 {
		if n < 3 {
			log.Fatal("buf too short!")
		}
		// 2 next bytes = cd (as in ABCDEFGH...)
		cd := int(buf[1]) + int(buf[2])<<8
		// return the decoded int and the amount of bytes used
		return cd, 3
	}

	// 0x81 means: value is contained in the next 4 more bytes
	if b == 0x81 {
		if n < 5 {
			log.Fatal("buf too short!")
		}
		// 4 next bytes = cdef (as in ABCDEFGH...)
		cdef := int(buf[1]) + int(buf[2])<<8 + int(buf[3])<<16 + int(buf[4])<<24
		// return the decoded int and the amount of bytes used
		return cdef, 5
	}

	// return the decoded int and the amount of bytes used
	if b > 0x7F {
		return int(b) - int(1<<8), 1
	}
	return int(b), 1
}

// converts the next bytes up to the first \0 byte into a string
func getString(buf []byte) (string, int) {
	end := bytes.Index(buf, []byte{0}) + 1
	str := string(decodeCubecode(buf[:end]))
	return str, end
}

// returns a decoded int and sets the position to the next attribute's first byte
func dumpInt(buf []byte) int {
	decodedInt, bytesRead := getInt(buf[positionInResponse:])
	positionInResponse = positionInResponse + bytesRead
	return decodedInt
}

// returns a string and sets the position to the next attribute's first byte
func dumpString(buf []byte) string {
	decodedString, bytesRead := getString(buf[positionInResponse:])
	positionInResponse = positionInResponse + bytesRead
	return decodedString
}
