package userid

import "encoding/binary"

import "crypto/rand"

const IDLength = 4

func GenerateUserID() (uint32, error) {
	b := make([]byte, IDLength)
	_, err := rand.Read(b)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(b), nil
}
