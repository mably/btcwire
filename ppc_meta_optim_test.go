package btcwire_test

import (
	"fmt"
	"github.com/mably/btcwire"
	"github.com/mably/ppcutil"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type metaElement struct {
	meta *btcwire.Meta
	next *metaElement
}

var _time time.Time

func tick(msg string) {
	if msg == "" {
		_time = time.Now()
		return
	}
	now := time.Now()
	fmt.Printf("%s: %v\n", msg, now.Sub(_time))
	_time = now
}

func TestReadCBlockIndex(t *testing.T) {
	tick("")
	root := ppcutil.ReadCBlockIndex("../ppcutil/testdata/blkindex.csv")
	tick("csv loading")
	r := root
	if r.Height != 0 {
		t.Errorf("bad root height, have %d, want %d", r.Height, 0)
	}
	for r.Next != nil {
		r = r.Next
	}
	tick("iteration")
	if r.Height != 131325 {
		t.Errorf("bad head height, have %d, want %d", r.Height, 131325)
	}
	name := filepath.Join(os.TempDir(), "ppc-meta.ser")
	file, err := os.Create(name)
	if err != nil {
		t.Error(err)
		return
	}
	defer file.Close()
	tick("")
	var rootMeta, prev *metaElement
	for r = root; r != nil; r = r.Next {
		m := new(metaElement)
		m.meta = new(btcwire.Meta)
		hash, _ := btcwire.NewShaHash(r.HashProofOfStake)
		m.meta.HashProofOfStake.SetBytes(hash.Bytes())
		m.meta.StakeModifier = new(big.Int).SetBytes(r.StakeModifier).Uint64()
		m.meta.StakeModifierChecksum = uint32(new(big.Int).SetBytes(r.StakeModifier).Uint64())
		m.meta.ChainTrust.SetBytes(r.ChainTrust)
		if r.GeneratedModifier {
			m.meta.Flags = 1 << 0
		}
		if r.EntropyBit {
			m.meta.Flags |= 1 << 1
		}
		if r.ProofOfStake {
			m.meta.Flags |= 1 << 2
		}

		if prev != nil {
			prev.next = m
		} else {
			rootMeta = m
		}
		prev = m
	}
	tick("convert")
	for m := rootMeta; m != nil; m = m.next {
		m.meta.Serialize(file)
	}
	tick("serialize")
	file.Sync()
	fi, err := os.Stat(name)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("file size: %d\n", fi.Size())

	file, err = os.Open(name)
	if err != nil {
		t.Error(err)
		return
	}
	defer file.Close()
	tick("")
	for err == nil {
		err = new(btcwire.Meta).Deserialize(file)
	}
	tick("deserialize")
	if err != io.EOF {
		t.Error(err)
		return
	}
}
