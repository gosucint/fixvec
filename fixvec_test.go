package fixvec

import (
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"testing"
)

const (
	testIter = 10
)

func TestEmptyFixVec(t *testing.T) {
	vv := New(0, 0)
	Convey("When an input is empty", t, func() {
		So(vv.Num(), ShouldEqual, 0)
		So(vv.Blen(), ShouldEqual, 0)
	})
}

func TestLargeFixVec(t *testing.T) {
	Convey("When large FixVec is build", t, func() {
		for blen := uint8(0); blen < 64; blen++ {
			for iter := 0; iter < testIter; iter++ {
				num := uint64(rand.Int31n(1000))
				vv := New(num, blen)
				orig := make([]uint64, num)
				for i := uint64(0); i < num; i++ {
					orig[i] = uint64(rand.Int63() % (1 << blen))
					vv.Set(i, orig[i])
				}
				for i := 0; i < 10; i++ {
					ind := uint64(rand.Int31n(int32(num)))
					So(vv.Get(ind), ShouldEqual, orig[ind])
				}
			}
		}
	})
}

func TestLargeFixVecFromArray(t *testing.T) {
	Convey("When large FixVec is build", t, func() {
		for blen := uint8(0); blen < 64; blen++ {
			for iter := 0; iter < testIter; iter++ {
				num := uint64(rand.Int31n(1000))
				orig := make([]uint64, num)
				for i := uint64(0); i < num; i++ {
					orig[i] = uint64(rand.Int63() % (1 << blen))
				}
				vv := NewFromArray(orig)
				for i := 0; i < 10; i++ {
					ind := uint64(rand.Int31n(int32(num)))
					So(vv.Get(ind), ShouldEqual, orig[ind])
				}

				out, err := vv.MarshalBinary()
				So(err, ShouldBeNil)
				newvv := New(0, 0)
				err = newvv.UnmarshalBinary(out)
				So(err, ShouldBeNil)
				So(newvv.Num(), ShouldEqual, vv.Num())
				So(newvv.Blen(), ShouldEqual, vv.Blen())
				for i := 0; i < 10; i++ {
					ind := uint64(rand.Int31n(int32(num)))
					So(newvv.Get(ind), ShouldEqual, orig[ind])
				}
			}
		}
	})
}
