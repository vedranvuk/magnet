// Copyright 2013 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package magnet is a utility package for parsing magnet links.
package magnet

import (
	"encoding/base32"
	"encoding/hex"
	"errors"
	"net/url"
	"strconv"
	"strings"
)

// Magnet hash type.
type HashType int

// Hash types.
const (
	HashTTH      HashType = iota // Tiger Tree hash.
	HashSHA1                     // Secure Hash Algorithm 1.
	HashBitPrint                 // BitPrint.
	HashED2K                     // eDonkey2000 hash.
	HashAICH                     // Advanced Intelligent Corruption Handler.
	HashKazaa                    // Kazaa hash.
	HashBTIH                     // BitTorrent Info Hash.
	HashMD5                      // Message Digest 5.
)

var HashTypeMap = map[string]HashType{
	"tree:tiger": HashTTH,
	"sha1":       HashSHA1,
	"bitprint":   HashBitPrint,
	"ed2k":       HashED2K,
	"aich":       HashAICH,
	"kzhash":     HashKazaa,
	"btih":       HashBTIH,
	"md5":        HashMD5,
}

// Magnet key type
type KeyType int

// Key types.
const (
	KeyAcceptableSource KeyType = iota // Acceptable sources.
	KeyDisplayName                     // Display name.
	KeyKeywordTopic                    // Search keywords.
	KeyManifestTopic                   // Link to metafile containing list of MAGNETO manifests.
	KeyTrackerAddress                  // Tracker address.
	KeyExactLength                     // Exact length in bytes.
	KeyExactSource                     // p2p link.
	KeyExactTopic                      // URN containing file hash.
	KeySuplement                       // Suplemental keys (extensions).
)

// Maps Magnet string keys to KeyType type.
var KeyTypeMap = map[string]KeyType{
	"as": KeyAcceptableSource,
	"dn": KeyDisplayName,
	"kt": KeyKeywordTopic,
	"mt": KeyManifestTopic,
	"tr": KeyTrackerAddress,
	"xl": KeyExactLength,
	"xs": KeyExactSource,
	"xt": KeyExactTopic,
	"x.": KeySuplement,
}

var (
	ErrInvalidMagnet = errors.New("invalid magnet")
)

// Defines a Magnet hash.
type Hash struct {
	Type HashType // Hash type.
	Data []byte   // Hash data.
}

// Defines a Magnet URN.
type URN struct {
	Hashes []Hash // Hashes, some URNs have multiple.
}

// Create a new *URN structure from an urn string or an error.
func newURN(v string) (*URN, error) {
	r := &URN{}
	a := strings.Split(v, "urn:")
	if len(a) != 2 {
		return nil, ErrInvalidMagnet
	}
	urn := a[1]
	ht := HashType(-1)
	hd := ""
	for d := range HashTypeMap {
		if strings.HasPrefix(strings.ToLower(urn), d) {
			ht = HashTypeMap[d]
			hd = urn[len(d)+1:]
			break
		}
	}

	switch ht {
	case HashTTH, HashSHA1, HashAICH:
		data, err := base32.StdEncoding.DecodeString(hd)
		if err != nil {
			return nil, err
		}
		r.Hashes = append(r.Hashes, Hash{ht, data})
	case HashED2K, HashKazaa, HashBTIH:
		data, err := hex.DecodeString(hd)
		if err != nil {
			return nil, err
		}
		r.Hashes = append(r.Hashes, Hash{ht, data})
	case HashBitPrint:
		b := strings.Split(hd, ".")
		if len(b) != 2 {
			return nil, ErrInvalidMagnet
		}
		data, err := base32.StdEncoding.DecodeString(b[0])
		if err != nil {
			return nil, err
		}
		r.Hashes = append(r.Hashes, Hash{HashSHA1, data})
		data, err = base32.StdEncoding.DecodeString(b[1])
		if err != nil {
			return nil, err
		}
		r.Hashes = append(r.Hashes, Hash{HashTTH, data})
	}
	return r, nil
}

type Suplement struct {
	Key string
	Val string
}

type Magnet struct {
	AcceptableSources []url.URL     // Fall-back sources, direct download from a web server.
	DisplayNames      []string      // Filename/display name.
	KeywordTopics     []string      // Search keywords.
	ManifestTopics    []interface{} // Link to a list of links. URL or URN.
	TrackerAddresses  []string      // Tracker addresses.
	ExactLength       int64         // Filesize. By logic only one key of thsi type should exist.
	ExactSources      []string
	ExactTopics       []URN
	Suplements        map[string][]Suplement
}

// Defines a magnet key.
type magnetKey struct {
	Type KeyType // One of supported key types.
	Indx int     // Optional index if it's a multiple key:value type pair.
	Supl string  // If it's a suplemental key, this is its' value.
}

// Creates a new *magnetKey structure from a magnet key string "k" or an error.
func newMagnetKey(k string) (*magnetKey, error) {
	r := magnetKey{}

	// Fail on unsupported key.
	if len(k) < 2 {
		return nil, ErrInvalidMagnet
	}
	var ok bool
	if r.Type, ok = KeyTypeMap[k[0:2]]; !ok {
		return nil, ErrInvalidMagnet
	}

	// KeySuplement special case.
	if r.Type == KeySuplement {
		b := strings.SplitAfterN(k, ".", 1)
		if len(b) < 2 {
			return nil, ErrInvalidMagnet
		}
		r.Supl = b[1]
	}

	// Get index.
	c := strings.Split(k, ".")
	if len(c) > 2 {
		return nil, ErrInvalidMagnet
	}
	if len(c) == 2 {
		v, err := strconv.Atoi(c[1])
		if err != nil {
			return nil, ErrInvalidMagnet
		}
		r.Indx = v
	}
	return &r, nil
}

// Parses key:value pairs, converts to go types and adds to self.
func (m *Magnet) parseKeyVal(k, v string) error {
	mk, err := newMagnetKey(k)
	if err != nil {
		return err
	}
	switch mk.Type {
	case KeyAcceptableSource:
		u, err := url.Parse(v)
		if err != nil {
			return err
		}
		m.AcceptableSources = append(m.AcceptableSources, *u)
	case KeyDisplayName:
		u, err := url.QueryUnescape(v)
		if err != nil {
			return err
		}
		m.DisplayNames = append(m.DisplayNames, u)
	case KeyKeywordTopic:
		u, err := url.QueryUnescape(v)
		if err != nil {
			return err
		}
		m.KeywordTopics = append(m.KeywordTopics, u)
	case KeyManifestTopic:
		if strings.HasPrefix(strings.ToLower(v), "urn") {
			u, err := newURN(v)
			if err != nil {
				return err
			}
			m.ManifestTopics = append(m.ManifestTopics, u)
		} else {
			u, err := url.QueryUnescape(v)
			if err != nil {
				return err
			}
			m.ManifestTopics = append(m.ManifestTopics, u)
		}
	case KeyTrackerAddress:
		u, err := url.QueryUnescape(v)
		if err != nil {
			return err
		}
		m.TrackerAddresses = append(m.TrackerAddresses, u)
	case KeyExactLength:
		u, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
		m.ExactLength = u
	case KeyExactSource:
		m.ExactSources = append(m.ExactSources, v)
	case KeyExactTopic:
		u, err := newURN(v)
		if err != nil {
			return err
		}
		m.ManifestTopics = append(m.ManifestTopics, u)
	case KeySuplement:
		m.Suplements[k[0:2]] = append(m.Suplements[k[0:2]], Suplement{mk.Supl, v})
	}
	return nil
}

// Does the main split then iterates over key:value pairs.
func (m *Magnet) parseMagnet(s string) error {
	a := strings.Split(s, ":?")
	if len(a) != 2 {
		goto error
	}
	if strings.ToLower(a[0]) != "magnet" {
		goto error
	}
	a = strings.Split(a[1], "&")
	if len(a) == 0 {
		goto error
	}
	for _, v := range a {
		b := strings.Split(v, "=")
		if len(b) < 1 {
			goto error
		}
		if err := m.parseKeyVal(b[0], b[1]); err != nil {
			return err
		}
	}
	return nil
error:
	return ErrInvalidMagnet
}

// Creates a new *Magnet structure from a magnet string "s" or an error.
func NewMagnet(s string) (*Magnet, error) {
	m := &Magnet{}
	if err := m.parseMagnet(s); err != nil {
		return nil, err
	}
	return m, nil
}
