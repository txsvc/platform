package id

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"io"
)

var (
	// see https://golang.org/pkg/hash/crc32/
	crc32q *crc32.Table = crc32.MakeTable(0xD5828281)
)

// Checksum creates a unique CRC32 checksum from a string
func Checksum(a string) string {
	return fmt.Sprintf("%08x", crc32.Checksum([]byte(a), crc32q))
}

// Fingerprint creates a unique ID (MD5) from a string
func Fingerprint(a string) string {
	hash := md5.Sum([]byte(fmt.Sprint(a)))
	return hex.EncodeToString(hash[:])
}

// UUID generates a random UUID according to RFC 4122
func UUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

// RandomToken generates a random token similar to a to a RFC 4122 UID
func RandomToken(prefix string) (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%s-%x-%x", prefix, uuid[0:10], uuid[10:]), nil
}

// ShortUUID returns a short (6 bytes) UID
func ShortUUID() (string, error) {
	uuid := make([]byte, 6)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	uuid[4] = uuid[4]&^0xc0 | 0x80
	uuid[2] = uuid[2]&^0xf0 | 0x40

	return fmt.Sprintf("%x", uuid[0:6]), nil
}

// SimpleUUID generates a random UUID according to RFC 4122, without any dashes
func SimpleUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x", uuid), nil
}
