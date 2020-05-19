package partialsum

import (
	"github.com/ugorji/go/codec"
	"github.com/vsivsi/rsdic"
)

//////////////////////////////////////////////////////////////////
// Note: interface removed for performance and marshaling reasons
//////////////////////////////////////////////////////////////////

// // PartialSumInt stores non-negative integers V[0...N)
// // and supports Sum, Find in O(1) time
// // using at most (S + N) bits where S is the sum of V[0...N)
// type PartialSumInt interface {
// 	// Increment add V[ind] += val
// 	// ind should hold ind >= Num
// 	IncTail(ind uint64, val uint64)

// 	// Num returns the number of vals
// 	Num() uint64

// 	// AllSum returns the sum of all vals
// 	AllSum() uint64

// 	// Lookup returns V[i] in O(1) time
// 	Lookup(ind uint64) (val uint64)

// 	// Sum returns V[0]+V[1]+...+V[ind-1] in O(1) time
// 	Sum(ind uint64) (sum uint64)

// 	// Lookup returns V[i] and V[0]+V[1]+...+V[i-1] in O(1) time
// 	LookupAndSum(ind uint64) (val uint64, sum uint64)

// 	// Find returns ind satisfying Sum(ind) <= val < Sum(ind+1)
// 	// and val - Sum(ind). If there are multiple inds
// 	// satisfiying this condition, return the minimum one.
// 	Find(val uint64) (ind uint64, offset uint64)

// 	// MarshalBinary encodes VecString into a binary form and returns the result.
// 	MarshalBinary() ([]byte, error)

// 	// UnmarshalBinary decodes the FixVec form a binary from generated MarshalBinary
// 	UnmarshalBinary([]byte) error
// }

// PartialSum stores non-negative integers V[0...N)
// and supports Sum, Find in O(1) time
// using at most (S + N) bits where S is the sum of V[0...N)
type PartialSum struct {
	dic *rsdic.RSDic
}

// New returns a new partial sum data struct
func New() PartialSum {
	return PartialSum{
		dic: rsdic.New(),
	}
}

// IncTail adds: V[ind] += val
// ind should hold ind >= Num
func (ps *PartialSum) IncTail(ind uint64, val uint64) {
	for i := ps.Num(); i <= ind; i++ {
		ps.dic.PushBack(true)
	}
	for i := uint64(0); i < val; i++ {
		ps.dic.PushBack(false)
	}
}

// Num returns the number of vals
func (ps PartialSum) Num() uint64 {
	return ps.dic.OneNum()
}

// AllSum returns the sum of all vals
func (ps PartialSum) AllSum() uint64 {
	return ps.dic.ZeroNum()
}

// Lookup returns V[i] in O(1) time
func (ps PartialSum) Lookup(ind uint64) (val uint64) {
	return ps.dic.Select(ind+1, true) - ps.dic.Select(ind, true) - 1
}

// Sum returns V[0]+V[1]+...+V[ind-1] in O(1) time
func (ps PartialSum) Sum(ind uint64) (sum uint64) {
	return ps.dic.Rank(ps.dic.Select(ind, true), false)
}

// LookupAndSum returns V[i] and V[0]+V[1]+...+V[i-1] in O(1) time
func (ps PartialSum) LookupAndSum(ind uint64) (val uint64, sum uint64) {
	indPos := ps.dic.Select(ind, true)
	sum = ps.dic.Rank(indPos, false)
	val = ps.dic.Select(ind+1, true) - indPos - 1
	return
}

// Find returns ind satisfying Sum(ind) <= val < Sum(ind+1)
// and val - Sum(ind). If there are multiple inds
// satisfiying this condition, return the minimum one.
func (ps PartialSum) Find(val uint64) (ind uint64, offset uint64) {
	pos := ps.dic.Select(val, false)
	ind = ps.dic.Rank(pos, true) - 1
	offset = pos - ps.dic.Select(ind, true) - 1
	return
}

// MarshalBinary encodes VecString into a binary form and returns the result.
func (ps PartialSum) MarshalBinary() (out []byte, err error) {
	var mh codec.MsgpackHandle
	enc := codec.NewEncoderBytes(&out, &mh)
	err = enc.Encode(ps.dic)
	if err != nil {
		return
	}
	return
}

// UnmarshalBinary decodes the FixVec form a binary from generated MarshalBinary
func (ps *PartialSum) UnmarshalBinary(in []byte) (err error) {
	var mh codec.MsgpackHandle
	dec := codec.NewDecoderBytes(in, &mh)
	err = dec.Decode(&ps.dic)
	if err != nil {
		return
	}
	return
}
