package storage

import "strings"

const SEP = '/'

type Key struct {
	parts  []string
	prefix bool
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

func (k *Key) Prefix(val bool) *Key {
	k.prefix = val
	return k
}

func (k *Key) IsPrefix() bool {
	return k.prefix
}

func (k *Key) Copy() *Key {
	p := make([]string, len(k.parts))
	copy(k.parts, p)
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
