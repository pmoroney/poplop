package poplop

import (
	"bytes"
	"crypto/md5"
	"crypto/sha512"
	"encoding/ascii85"
	"encoding/base64"
	"errors"
	"io"
	"strconv"
	"strings"
	"unicode"
)

type Scheme struct {
	Name      string
	Counter   int
	Username  string
	URL       string
	Notes     string
	Forbidden string
	MaxLength int
	Legacy    bool
	// "!\"#$%&'()*+,-./:;<=>?@[\\]^_'",
}

func (n Scheme) mapping() func(rune) rune {
	return func(r rune) rune {
		if strings.ContainsRune(n.Forbidden, r) {
			return -1
		}
		return r
	}
}

func (n Scheme) Hash(master string) (string, error) {
	if n.Legacy {
		return oplop(master, n)
	}
	return poplop(master, n)
}

func oplop(master string, n Scheme) (string, error) {
	h := md5.New()
	io.WriteString(h, master)
	io.WriteString(h, n.Name)
	if n.Counter > 0 {
		io.WriteString(h, strconv.Itoa(n.Counter))
	}
	str := requireDigit(base64.RawStdEncoding.EncodeToString(h.Sum(nil)))
	if n.Forbidden != "" {
		str = strings.Map(n.mapping(), str)
	}

	length := 8
	if n.MaxLength > 0 && n.MaxLength < length {
		length = n.MaxLength
	}

	if len(str) < length {
		return "", errors.New("hash not long enough")
	}
	return str[:length], nil
}

func requireDigit(pass string) string {
	var digits bytes.Buffer
	for i, r := range pass {
		if unicode.IsDigit(r) {
			if i < 8 {
				return pass
			}
			digits.WriteRune(r)
		} else {
			if digits.Len() > 0 {
				return digits.String() + pass
			}
		}
	}
	if digits.Len() > 0 {
		return digits.String() + pass
	}
	return "1" + pass
}

func poplop(master string, n Scheme) (string, error) {
	h := sha512.New()
	io.WriteString(h, master)
	io.WriteString(h, n.Name)
	if n.Counter > 0 {
		io.WriteString(h, strconv.Itoa(n.Counter))
	}
	hash := h.Sum(nil)
	dst := make([]byte, ascii85.MaxEncodedLen(len(hash)))
	ascii85.Encode(dst, hash)

	str := string(dst)
	if n.Forbidden != "" {
		str = strings.Map(n.mapping(), str)
	}

	length := 20
	if n.MaxLength > 0 {
		length = n.MaxLength
	}

	if len(str) < length {
		return "", errors.New("hash not long enough " + str)
	}
	return str[:length], nil
}
