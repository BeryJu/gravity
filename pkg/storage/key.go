package storage

import "strings"

const SEP = '/'

type Key struct {
	parts  []string
	prefix bool
}

func KeyFromString(raw string) *Key {
	// Remove first element as keys always start with a slash
	parts := strings.Split(raw, string(SEP))[1:]
	prefix := false
	if strings.HasSuffix(raw, string(SEP)) {
		prefix = true
		parts = parts[:len(parts)-1]
	}
	return &Key{
		parts:  parts,
		prefix: prefix,
	}
}

func (c *Client) Key(parts ...string) *Key {
	return &Key{
		parts: parts,
	}
}

func (k *Key) Add(parts ...string) *Key {
	k.parts = append(k.parts, parts...)
	return k
}

func (k *Key) Up() *Key {
	k.parts = k.parts[:len(k.parts)-1]
	return k
}

func (k *Key) Prefix(val bool) *Key {
	k.prefix = val
	return k
}

func (k *Key) IsPrefix() bool {
	return k.prefix
}

func (k *Key) Copy() *Key {
	p := make([]string, len(k.parts))
	copy(p, k.parts)
	return &Key{
		parts:  p,
		prefix: k.prefix,
	}
}

func (k *Key) String() string {
	b := strings.Builder{}
	b.WriteRune(SEP)
	for idx, part := range k.parts {
		if strings.HasPrefix(part, string(SEP)) {
			b.WriteString(part[1:])
		} else {
			b.WriteString(part)
		}
		if idx != len(k.parts)-1 {
			b.WriteRune(SEP)
		}
	}
	if k.prefix {
		b.WriteRune(SEP)
	}
	return b.String()
}
