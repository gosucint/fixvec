// package fixvec provides a vector representation of value using fixed bits
package fixvec

import (
	"github.com/ugorji/go/codec"
)

// FixVec provides a vector representation of value using fixed bits
// Conceptually, FixVec represents a vector V[0...num), and each
// value V[i] can represent in [0...2^(blen))
// The total working space is num * blen bits (+ some small overhead).
type FixVec interface {
	// Get returns　V[ind]
	Get(ind uint64) uint64

	// Set sets V[ind] = val (val will be masked by &((1 << blen) -1))
	Set(ind uint64, val uint64)

	// Blen returns the number of bits for value representation
	Blen() uint8

	// Num returns the number of elemens in V
	Num() uint64

	// MarshalBinary encodes FixVec into a binary form and returns the result.
	MarshalBinary() ([]byte, error)

	// UnmarshalBinary decodes the FixVec from a binary from generated MarshalBinary
	UnmarshalBinary([]byte) error
}

// New returns FixVec represents V[0...n) where each element is less than (1 << blen).
func New(n uint64, bl uint8) FixVec {
	if bl > 64 {
		return nil
	}
	return &fixVecImpl{
		bits: make([]uint64, floorBlock(n*uint64(bl), 64)),
		num:  n,
		blen: bl,
	}
}

// NewFromArray returns a FixVec represents V[0...n) where each element V[i] equals to vs[i]
func NewFromArray(vs []uint64) FixVec {
	dim := uint64(0)
	for _, v := range vs {
		if v >= dim {
			dim = v + 1
		}
	}
	num := uint64(len(vs))
	if dim == 0 {
		return New(num, 0)
	}
	max := dim - 1
	blen := uint8(0)
	for ; (max >> blen) > 0; blen++ {
	}
	fv := New(num, blen)
	for i := uint64(0); i < num; i++ {
		fv.Set(i, vs[i])
	}
	return fv
}

type fixVecImpl struct {
	bits []uint64
	num  uint64
	blen uint8
}

func (vv fixVecImpl) Get(ind uint64) uint64 {
	blen := uint64(vv.blen)
	if blen == 0 {
		return 0
	}
	pos := ind * blen
	block, offset := pos/64, pos%64
	ret := vv.bits[block] >> offset
	if offset+blen > 64 {
		ret |= vv.bits[block+1] << (64 - offset)
	}
	if vv.blen == 64 {
		return ret
	}
	return ret & ((1 << vv.blen) - 1)
}

func (vv *fixVecImpl) Set(ind uint64, val uint64) {
	blen := uint64(vv.blen)
	if blen == 0 {
		return
	}
	pos := ind * blen
	block, offset := pos/64, pos%64
	vv.bits[block] |= val << offset
	if offset+blen > 64 {
		vv.bits[block+1] |= (val >> (64 - offset))
	}
}

func (vv fixVecImpl) Blen() uint8 {
	return vv.blen
}

func (vv fixVecImpl) Num() uint64 {
	return vv.num
}

func (vv fixVecImpl) PushBack(val uint64) uint64 {

	return vv.num
}

func (vv fixVecImpl) MarshalBinary() (out []byte, err error) {
	var bh codec.MsgpackHandle
	enc := codec.NewEncoderBytes(&out, &bh)
	err = enc.Encode(vv.bits)
	if err != nil {
		return
	}
	err = enc.Encode(vv.num)
	if err != nil {
		return
	}
	err = enc.Encode(vv.blen)
	if err != nil {
		return
	}
	return
}

func (vv *fixVecImpl) UnmarshalBinary(in []byte) (err error) {
	var bh codec.MsgpackHandle
	dec := codec.NewDecoderBytes(in, &bh)
	err = dec.Decode(&vv.bits)
	if err != nil {
		return
	}
	err = dec.Decode(&vv.num)
	if err != nil {
		return
	}
	err = dec.Decode(&vv.blen)
	if err != nil {
		return
	}
	return nil
}

func floorBlock(num uint64, div uint64) uint64 {
	return (num + div - 1) / div
}
