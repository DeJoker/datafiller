// Package datafiller implements function for randomly filling passed
// datastructures by mock sample data.
package datafiller

import (
	"math/rand"
	"reflect"
	"time"
)

func init() {
	packagesInit()
}

const (
	taggedStructKey = "datafiller"
)

// Function Fill takes a pointer to variable of any type and fills the variable
// by with sample data. It panics if the passed value is not a pointer.
func Fill(i interface{}) {
	f := NewFiller()
	f.Seed = time.Now().Unix() + rand.Int63n(100)
	f.Fill(i)
}

type Filler struct {
	Seed     int64
	randSeed *rand.Rand
	ArrayMin int
	ArrayMax int
}

type FillerArg interface {
	AddArg(*Filler)
}

type Seed int64

func (sed Seed) AddArg(f *Filler) {
	f.Seed = int64(sed)
}

type ArrayMin int

func (amn ArrayMin) AddArg(f *Filler) {
	f.ArrayMin = int(amn)
}

type ArrayMax int

func (amx ArrayMax) AddArg(f *Filler) {
	f.ArrayMax = int(amx)
}

func NewFiller(args ...FillerArg) *Filler {
	f := &Filler{
		Seed:     -1,
		ArrayMin: 1,
		ArrayMax: 3,
	}
	for _, a := range args {
		a.AddArg(f)
	}
	return f
}

func (f *Filler) Fill(i interface{}) {
	valPtr := reflect.ValueOf(i)

	if valPtr.Kind() != reflect.Ptr && valPtr.Kind() != reflect.UnsafePointer {
		panic("Passed argument is not a pointer.")
	}

	val := reflect.Indirect(valPtr)

	f.randSeed = rand.New(rand.NewSource(f.Seed))
	f.recursiveSet(val)

}

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

//会产生大量的丢弃的case,造成重选和浪费。rand.Int63会产生63bit的随机数
//如果我们把它分成6份，那么一次就可以产生10个6bit的随机数。这样就减少了浪费
func RandStringBytesMaskImpr(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}


func (self *Filler) recursiveSet(val reflect.Value) {
	if val.CanSet() {
		var fullPath string
		fullPath = val.Type().PkgPath() + "." + val.Type().Name()
		pkgVal, ok := packages[fullPath]
		if ok {
			val.Set(pkgVal)
			return
		}

		if val.Kind() == reflect.Int ||
			val.Kind() == reflect.Int8 ||
			val.Kind() == reflect.Int16 ||
			val.Kind() == reflect.Int32 ||
			val.Kind() == reflect.Int64 {
			val.SetInt(self.randSeed.Int63n(100))
			return
		} else if val.Kind() == reflect.Uint ||
			val.Kind() == reflect.Uint8 ||
			val.Kind() == reflect.Uint16 ||
			val.Kind() == reflect.Uint32 ||
			val.Kind() == reflect.Uint64 {
			val.SetUint(uint64(self.randSeed.Int63n(100)))
			return
		} else if val.Kind() == reflect.Float32 ||
			val.Kind() == reflect.Float64 {
			val.SetFloat(float64(self.randSeed.Float32()))
			return
		} else if val.Kind() == reflect.Complex64 {
			cpx := complex128(complex(self.randSeed.Float32(), self.randSeed.Float32()))
			val.SetComplex(cpx)
			return
		} else if val.Kind() == reflect.Complex128 {
			cpx := complex128(complex(self.randSeed.Float64(), self.randSeed.Float64()))
			val.SetComplex(cpx)
			return
		} else if val.Kind() == reflect.Bool {
			if self.randSeed.Int63n(2) == 0 {
				val.SetBool(false)
			} else {
				val.SetBool(true)
			}
			return
		} else if val.Kind() == reflect.String {
			val.SetString(RandStringBytesMaskImpr(20))
			return
		} else if val.Kind() == reflect.Struct {
			lngth := val.NumField()
			strType := val.Type()

			for i := 0; i < lngth; i++ {
				if strType.Field(i).Tag.Get(taggedStructKey) == "" {
					self.recursiveSet(val.Field(i))
				} else if strType.Field(i).Tag.Get(taggedStructKey) == "-" {
				} else {
					// advStrTag := strType.Field(i).Tag.Get(taggedStructKey)
				}
			}
			return
		} else if val.Kind() == reflect.Ptr {
			tp := val.Type().Elem()
			nw := reflect.New(tp)
			val.Set(nw)
			self.recursiveSet(reflect.Indirect(val))
			return
		} else if val.Kind() == reflect.Map {
			// nm := reflect.NewMap
			// TODO(tvi): Finish.
			return
		} else if val.Kind() == reflect.Slice {
			// TODO(tvi): Fix slice length.
			typ := val.Type()
			elem := typ.Elem()
			nw := reflect.Zero(elem)
			m := reflect.MakeSlice(typ, 0, 1)
			m = reflect.Append(m, nw)
			m = reflect.Append(m, nw)
			val.Set(m)
			self.recursiveSet(val.Index(0))
			self.recursiveSet(val.Index(1))
			return
		}
	}
}
