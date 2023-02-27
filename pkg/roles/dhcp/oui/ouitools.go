// Updated version of https://github.com/dutchcoders/go-ouitools
// Package go-oui provides functions to work with MAC and OUI's
package oui

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net"
	"regexp"
	"strconv"
)

// https://code.wireshark.org/review/gitweb?p=wireshark.git;a=blob_plain;f=manuf
// Bigger than we need, not too big to worry about overflow
const big = 0xFFFFFF

var ErrInvalidMACAddress = errors.New("invalid MAC address")

// Hexadecimal to integer starting at &s[i0].
// Returns number, new offset, success.
func xtoi(s string, i0 int) (n int, i int, ok bool) {
	n = 0
	for i = i0; i < len(s); i++ {
		if '0' <= s[i] && s[i] <= '9' {
			n *= 16
			n += int(s[i] - '0')
		} else if 'a' <= s[i] && s[i] <= 'f' {
			n *= 16
			n += int(s[i]-'a') + 10
		} else if 'A' <= s[i] && s[i] <= 'F' {
			n *= 16
			n += int(s[i]-'A') + 10
		} else {
			break
		}
		if n >= big {
			return 0, i, false
		}
	}
	if i == i0 {
		return 0, i, false
	}
	return n, i, true
}

// xtoi2 converts the next two hex digits of s into a byte.
// If s is longer than 2 bytes then the third byte must be e.
// If the first two bytes of s are not hex digits or the third byte
// does not match e, false is returned.
func xtoi2(s string, e byte) (byte, bool) {
	if len(s) > 2 && s[2] != e {
		return 0, false
	}
	n, ei, ok := xtoi(s[:2], 0)
	return byte(n), ok && ei == 2
}

type HardwareAddr net.HardwareAddr

// ParseMAC parses s as an IEEE 802 MAC-48, EUI-48, or EUI-64 using one of the
// following formats:
//
//	01:23:45:67:89:ab
//	01:23:45:67:89:ab:cd:ef
//	01-23-45-67-89-ab
//	01-23-45-67-89-ab-cd-ef
//	0123.4567.89ab
//	0123.4567.89ab.cdef
func ParseOUI(s string, size int) (hw HardwareAddr, err error) {
	if s[2] == ':' || s[2] == '-' {
		if (len(s)+1)%3 != 0 {
			goto error
		}

		n := (len(s) + 1) / 3

		hw = make(HardwareAddr, size)
		for x, i := 0, 0; i < n; i++ {
			var ok bool
			if hw[i], ok = xtoi2(s[x:], s[2]); !ok {
				goto error
			}
			x += 3
		}
	} else {
		goto error
	}
	return hw, nil

error:
	return nil, ErrInvalidMACAddress
}

// Mask returns the result of masking the address with mask.
func (address HardwareAddr) Mask(mask []byte) []byte {
	n := len(address)
	if n != len(mask) {
		return nil
	}
	out := make([]byte, n)
	for i := 0; i < n; i++ {
		out[i] = address[i] & mask[i]
	}
	return out
}

type OuiDb struct {
	Blocks []AddressBlock
}

// New returns a new OUI database loaded from the specified file.
func New(reader io.Reader) *OuiDb {
	db := &OuiDb{}
	if err := db.Load(reader); err != nil {
		return nil
	}
	return db
}

// Lookup finds the OUI the address belongs to
func (m *OuiDb) Lookup(address HardwareAddr) *AddressBlock {
	for _, block := range m.Blocks {
		if block.Contains(address) {
			return &block
		}
	}

	return nil
}

// VendorLookup obtains the vendor organization name from the MAC address s.
func (m *OuiDb) VendorLookup(s string) (string, error) {
	addr, err := net.ParseMAC(s)
	if err != nil {
		return "", err
	}
	block := m.Lookup(HardwareAddr(addr))
	if block == nil {
		return "", ErrInvalidMACAddress
	}
	return block.Organization, nil
}

// VendorLookup obtains the vendor organization name from the MAC address s.
func (m *OuiDb) LookupString(s string) (*AddressBlock, error) {
	addr, err := net.ParseMAC(s)
	if err != nil {
		return nil, err
	}
	block := m.Lookup(HardwareAddr(addr))
	if block == nil {
		return nil, ErrInvalidMACAddress
	}
	return block, nil
}

func byteIndex(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}

func (m *OuiDb) Load(file io.Reader) error {
	// fieldsRe := regexp.MustCompile(`^(\S+)\t+(\S+)(\s+#\s+(\S.*))?`)
	fieldsRe := regexp.MustCompile(`^(\S+)\t+(\S+)(\s+(\S.*))?`)

	re := regexp.MustCompile(`((?:(?:[0-9a-zA-Z]{2})[-:]){2,5}(?:[0-9a-zA-Z]{2}))(?:/(\w{1,2}))?`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" || text[0] == '#' || text[0] == '\t' {
			continue
		}

		block := AddressBlock{}

		// Split input text into address, short organization name
		// and full organization name
		fields := fieldsRe.FindAllStringSubmatch(text, -1)
		addr := fields[0][1]
		if fields[0][4] != "" {
			block.Organization = fields[0][4]
		} else {
			block.Organization = fields[0][2]
		}

		matches := re.FindAllStringSubmatch(addr, -1)
		if len(matches) == 0 {
			continue
		}

		s := matches[0][1]

		i := byteIndex(s, '/')
		if i == -1 {
			block.Oui, _ = ParseOUI(s, 6)
			block.Mask = 24 // len(block.Oui) * 8
		} else {
			block.Oui, _ = ParseOUI(s[:i], 6)
			block.Mask, _ = strconv.Atoi(s[i+1:])
		}

		m.Blocks = append(m.Blocks, block)

		// create smart map
		for i := len(block.Oui) - 1; i >= 0; i-- {
			_ = block.Oui[i]
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func CIDRMask(ones, bits int) []byte {
	l := bits / 8
	m := make([]byte, l)

	n := uint(ones)
	for i := 0; i < l; i++ {
		if n >= 8 {
			m[i] = 0xff
			n -= 8
			continue
		}
		m[i] = ^byte(0xff >> n)
		n = 0
	}

	return (m)
}

// oui, mask, organization
type AddressBlock struct {
	Oui          HardwareAddr
	Mask         int
	Organization string
}

// Contains reports whether the mac address belongs to the OUI
func (b *AddressBlock) Contains(address HardwareAddr) bool {
	return (bytes.Equal(address.Mask(CIDRMask(b.Mask, len(b.Oui)*8)), b.Oui))
}
