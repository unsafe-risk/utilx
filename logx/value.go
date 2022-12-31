package logx

import (
	"math"
	"runtime"
	"sync"
	"time"
)

var pool1Bytes = sync.Pool{
	New: func() interface{} {
		v := make([]byte, 1)
		return &v
	},
}

var pool2Bytes = sync.Pool{
	New: func() interface{} {
		v := make([]byte, 2)
		return &v
	},
}

var pool4Bytes = sync.Pool{
	New: func() interface{} {
		v := make([]byte, 4)
		return &v
	},
}

var pool8Bytes = sync.Pool{
	New: func() interface{} {
		v := make([]byte, 8)
		return &v
	},
}

type varBytesPool struct {
	poolMap map[int]*sync.Pool
	sync.RWMutex
}

var poolVarBytes = varBytesPool{}

func (p *varBytesPool) Get(size int) *[]byte {
	p.RLock()
	pool, ok := p.poolMap[size]
	p.RUnlock()
	if !ok {
		p.Lock()
		pool, ok = p.poolMap[size]
		if !ok {
			pool = &sync.Pool{
				New: func() interface{} {
					v := make([]byte, size)
					return &v
				},
			}
			p.poolMap[size] = pool
		}
		p.Unlock()
	}
	return pool.Get().(*[]byte)
}

func (p *varBytesPool) Put(v *[]byte) {
	p.RLock()
	pool, ok := p.poolMap[len(*v)]
	p.RUnlock()
	if ok {
		pool.Put(v)
	}
}

type Value struct {
	key  *[]byte
	val  *[]byte
	kind int
}

func NewValue(key string) *Value {
	length := len(key)
	value := new(Value)
	bs := poolVarBytes.Get(length)
	copy(*bs, []byte(key))
	value.key = bs
	runtime.SetFinalizer(value, func(v *Value) {
		v.returnValue()
	})
	return value
}

func (v *Value) returnValue() {
	if v.val == nil {
		return
	}
	switch v.kind {
	case Int8:
		pool1Bytes.Put(v.val)
	case Int16:
		pool2Bytes.Put(v.val)
	case Int32:
		pool4Bytes.Put(v.val)
	case Int64:
		pool8Bytes.Put(v.val)
	case Uint8:
		pool1Bytes.Put(v.val)
	case Uint16:
		pool2Bytes.Put(v.val)
	case Uint32:
		pool4Bytes.Put(v.val)
	case Uint64:
		pool8Bytes.Put(v.val)
	case Float32:
		pool4Bytes.Put(v.val)
	case Float64:
		pool8Bytes.Put(v.val)
	case Bool:
		pool1Bytes.Put(v.val)
	case String:
		poolVarBytes.Put(v.val)
	case Bytes:
		poolVarBytes.Put(v.val)
	case Time:
		pool8Bytes.Put(v.val)
	}
	v.val = nil
}

func (v *Value) SetString(val string) {
	v.returnValue()
	length := len(val)
	bs := poolVarBytes.Get(length)
	copy(*bs, []byte(val))
	v.val = bs
	v.kind = String
}

func (v *Value) SetBytes(val []byte) {
	v.returnValue()
	length := len(val)
	bs := poolVarBytes.Get(length)
	copy(*bs, val)
	v.val = bs
	v.kind = Bytes
}

func (v *Value) SetInt8(val int8) {
	v.returnValue()
	bs := pool1Bytes.Get().(*[]byte)
	(*bs)[0] = byte(val)
	v.val = bs
	v.kind = Int8
}

func (v *Value) SetInt16(val int16) {
	v.returnValue()
	bs := pool2Bytes.Get().(*[]byte)
	(*bs)[0] = byte(val)
	(*bs)[1] = byte(val >> 8)
	v.val = bs
	v.kind = Int16
}

func (v *Value) SetInt32(val int32) {
	v.returnValue()
	bs := pool4Bytes.Get().(*[]byte)
	(*bs)[0] = byte(val)
	(*bs)[1] = byte(val >> 8)
	(*bs)[2] = byte(val >> 16)
	(*bs)[3] = byte(val >> 24)
	v.val = bs
	v.kind = Int32
}

func (v *Value) SetInt64(val int64) {
	v.returnValue()
	bs := pool8Bytes.Get().(*[]byte)
	(*bs)[0] = byte(val)
	(*bs)[1] = byte(val >> 8)
	(*bs)[2] = byte(val >> 16)
	(*bs)[3] = byte(val >> 24)
	(*bs)[4] = byte(val >> 32)
	(*bs)[5] = byte(val >> 40)
	(*bs)[6] = byte(val >> 48)
	(*bs)[7] = byte(val >> 56)
	v.val = bs
	v.kind = Int64
}

func (v *Value) SetUint8(val uint8) {
	v.returnValue()
	bs := pool1Bytes.Get().(*[]byte)
	(*bs)[0] = byte(val)
	v.val = bs
	v.kind = Uint8
}

func (v *Value) SetUint16(val uint16) {
	v.returnValue()
	bs := pool2Bytes.Get().(*[]byte)
	(*bs)[0] = byte(val)
	(*bs)[1] = byte(val >> 8)
	v.val = bs
	v.kind = Uint16
}

func (v *Value) SetUint32(val uint32) {
	v.returnValue()
	bs := pool4Bytes.Get().(*[]byte)
	(*bs)[0] = byte(val)
	(*bs)[1] = byte(val >> 8)
	(*bs)[2] = byte(val >> 16)
	(*bs)[3] = byte(val >> 24)
	v.val = bs
	v.kind = Uint32
}

func (v *Value) SetUint64(val uint64) {
	v.returnValue()
	bs := pool8Bytes.Get().(*[]byte)
	(*bs)[0] = byte(val)
	(*bs)[1] = byte(val >> 8)
	(*bs)[2] = byte(val >> 16)
	(*bs)[3] = byte(val >> 24)
	(*bs)[4] = byte(val >> 32)
	(*bs)[5] = byte(val >> 40)
	(*bs)[6] = byte(val >> 48)
	(*bs)[7] = byte(val >> 56)
	v.val = bs
	v.kind = Uint64
}

func (v *Value) SetFloat32(val float32) {
	v.returnValue()
	bs := pool4Bytes.Get().(*[]byte)
	bits := math.Float32bits(val)
	(*bs)[0] = byte(bits)
	(*bs)[1] = byte(bits >> 8)
	(*bs)[2] = byte(bits >> 16)
	(*bs)[3] = byte(bits >> 24)
	v.val = bs
	v.kind = Float32
}

func (v *Value) SetFloat64(val float64) {
	v.returnValue()
	bs := pool8Bytes.Get().(*[]byte)
	bits := math.Float64bits(val)
	(*bs)[0] = byte(bits)
	(*bs)[1] = byte(bits >> 8)
	(*bs)[2] = byte(bits >> 16)
	(*bs)[3] = byte(bits >> 24)
	(*bs)[4] = byte(bits >> 32)
	(*bs)[5] = byte(bits >> 40)
	(*bs)[6] = byte(bits >> 48)
	(*bs)[7] = byte(bits >> 56)
	v.val = bs
	v.kind = Float64
}

func (v *Value) SetBool(val bool) {
	v.returnValue()
	bs := pool1Bytes.Get().(*[]byte)
	if val {
		(*bs)[0] = 1
	} else {
		(*bs)[0] = 0
	}
	v.val = bs
	v.kind = Bool
}

func (v *Value) SetTime(val time.Time) {
	v.returnValue()
	bs := pool8Bytes.Get().(*[]byte)
	bits := uint64(val.UnixNano())
	(*bs)[0] = byte(bits)
	(*bs)[1] = byte(bits >> 8)
	(*bs)[2] = byte(bits >> 16)
	(*bs)[3] = byte(bits >> 24)
	(*bs)[4] = byte(bits >> 32)
	(*bs)[5] = byte(bits >> 40)
	(*bs)[6] = byte(bits >> 48)
	(*bs)[7] = byte(bits >> 56)
	v.val = bs
	v.kind = Time
}
