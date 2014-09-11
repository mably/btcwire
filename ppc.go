// Copyright (c) 2014-2014 PPCD developers.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package btcwire

import (
	"encoding/binary"
	"io"
	"math/big"
)

type Meta struct {
	StakeModifier         uint64
	StakeModifierChecksum uint32 // checksum of index; in-memeory only (main.h)
	HashProofOfStake      ShaHash
	Flags                 uint32
	ChainTrust            big.Int
}

func (m *Meta) Serialize(w io.Writer) error {
	e := binary.Write(w, binary.LittleEndian, &m.StakeModifier)
	if e != nil {
		return e
	}
	binary.Write(w, binary.LittleEndian, &m.StakeModifierChecksum)
	if e != nil {
		return e
	}
	binary.Write(w, binary.LittleEndian, &m.Flags)
	if e != nil {
		return e
	}
	binary.Write(w, binary.LittleEndian, &m.HashProofOfStake)
	if e != nil {
		return e
	}
	bytes := m.ChainTrust.Bytes()
	var blen byte
	blen = byte(len(bytes))
	binary.Write(w, binary.LittleEndian, &blen)
	if e != nil {
		return e
	}
	binary.Write(w, binary.LittleEndian, &bytes)
	if e != nil {
		return e
	}
	return nil
}

func (m *Meta) Deserialize(r io.Reader) error {
	e := binary.Read(r, binary.LittleEndian, &m.StakeModifier)
	if e != nil {
		return e
	}
	e = binary.Read(r, binary.LittleEndian, &m.StakeModifierChecksum)
	if e != nil {
		return e
	}
	e = binary.Read(r, binary.LittleEndian, &m.Flags)
	if e != nil {
		return e
	}
	e = binary.Read(r, binary.LittleEndian, &m.HashProofOfStake)
	if e != nil {
		return e
	}

	var blen byte
	e = binary.Read(r, binary.LittleEndian, &blen)
	if e != nil {
		return e
	}
	var arr = make([]byte, blen)
	e = binary.Read(r, binary.LittleEndian, &arr)
	if e != nil {
		return e
	}
	m.ChainTrust.SetBytes(arr)
	return nil
}
