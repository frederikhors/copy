package funcs

import (
	"reflect"
	"sync"
	"time"
	"unsafe"
)

type funcKey struct {
	Src  reflect.Type
	Dest reflect.Type
}

func typeOf(v interface{}) reflect.Type {
	return reflect.TypeOf(v)
}

func typeOfPointer(v interface{}) reflect.Type {
	return reflect.PtrTo(reflect.TypeOf(v))
}

// CopyFuncs is the storage of functions intended for copying data.
type CopyFuncs struct {
	mu    sync.RWMutex
	funcs map[funcKey]func(dst, src unsafe.Pointer)
	sizes []func(dst, src unsafe.Pointer)
}

// Get the copy function for the pair of types, if it is not found then nil is returned.
func (t *CopyFuncs) Get(dst, src reflect.Type) func(dst, src unsafe.Pointer) {
	t.mu.RLock()
	f := t.funcs[funcKey{Src: src, Dest: dst}]
	t.mu.RUnlock()
	if f != nil {
		return f
	}

	if dst.Kind() != src.Kind() {
		return nil
	}

	if dst.Kind() == reflect.String {
		// TODO
		return nil
	}

	same := dst == src

	switch dst.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
		same = same || dst.Elem() == src.Elem()
	}

	if same && dst.Size() == src.Size() && src.Size() > 0 && src.Size() <= uintptr(len(t.sizes)) {
		return t.sizes[src.Size()-1]
	}

	return nil
}

// Set the copy function for the pair of types.
func (t *CopyFuncs) Set(dst, src reflect.Type, f func(dst, src unsafe.Pointer)) {
	t.mu.Lock()
	t.funcs[funcKey{Src: src, Dest: dst}] = f
	t.mu.Unlock()
}

// Get the copy function for the pair of types, if it is not found then nil is returned.
func Get(dst, src reflect.Type) func(dst, src unsafe.Pointer) {
	return funcs.Get(dst, src)
}

// Set the copy function for the pair of types.
func Set(dst, src reflect.Type, f func(dst, src unsafe.Pointer)) {
	funcs.Set(dst, src, f)
}

var funcs = &CopyFuncs{
	funcs: map[funcKey]func(dst, src unsafe.Pointer){
		// int to int
		{Src: typeOf(int(0)), Dest: typeOf(int(0))}:               copyIntToInt,
		{Src: typeOfPointer(int(0)), Dest: typeOf(int(0))}:        copyPIntToInt,
		{Src: typeOf(int(0)), Dest: typeOfPointer(int(0))}:        copyIntToPInt,
		{Src: typeOfPointer(int(0)), Dest: typeOfPointer(int(0))}: copyPIntToPInt,
		// int8 to int
		{Src: typeOf(int8(0)), Dest: typeOf(int(0))}:               copyInt8ToInt,
		{Src: typeOfPointer(int8(0)), Dest: typeOf(int(0))}:        copyPInt8ToInt,
		{Src: typeOf(int8(0)), Dest: typeOfPointer(int(0))}:        copyInt8ToPInt,
		{Src: typeOfPointer(int8(0)), Dest: typeOfPointer(int(0))}: copyPInt8ToPInt,
		// int16 to int
		{Src: typeOf(int16(0)), Dest: typeOf(int(0))}:               copyInt16ToInt,
		{Src: typeOfPointer(int16(0)), Dest: typeOf(int(0))}:        copyPInt16ToInt,
		{Src: typeOf(int16(0)), Dest: typeOfPointer(int(0))}:        copyInt16ToPInt,
		{Src: typeOfPointer(int16(0)), Dest: typeOfPointer(int(0))}: copyPInt16ToPInt,
		// int32 to int
		{Src: typeOf(int32(0)), Dest: typeOf(int(0))}:               copyInt32ToInt,
		{Src: typeOfPointer(int32(0)), Dest: typeOf(int(0))}:        copyPInt32ToInt,
		{Src: typeOf(int32(0)), Dest: typeOfPointer(int(0))}:        copyInt32ToPInt,
		{Src: typeOfPointer(int32(0)), Dest: typeOfPointer(int(0))}: copyPInt32ToPInt,
		// int64 to int
		{Src: typeOf(int64(0)), Dest: typeOf(int(0))}:               copyInt64ToInt,
		{Src: typeOfPointer(int64(0)), Dest: typeOf(int(0))}:        copyPInt64ToInt,
		{Src: typeOf(int64(0)), Dest: typeOfPointer(int(0))}:        copyInt64ToPInt,
		{Src: typeOfPointer(int64(0)), Dest: typeOfPointer(int(0))}: copyPInt64ToPInt,
		// uint to int
		{Src: typeOf(uint(0)), Dest: typeOf(int(0))}:               copyUintToInt,
		{Src: typeOfPointer(uint(0)), Dest: typeOf(int(0))}:        copyPUintToInt,
		{Src: typeOf(uint(0)), Dest: typeOfPointer(int(0))}:        copyUintToPInt,
		{Src: typeOfPointer(uint(0)), Dest: typeOfPointer(int(0))}: copyPUintToPInt,
		// uint8 to int
		{Src: typeOf(uint8(0)), Dest: typeOf(int(0))}:               copyUint8ToInt,
		{Src: typeOfPointer(uint8(0)), Dest: typeOf(int(0))}:        copyPUint8ToInt,
		{Src: typeOf(uint8(0)), Dest: typeOfPointer(int(0))}:        copyUint8ToPInt,
		{Src: typeOfPointer(uint8(0)), Dest: typeOfPointer(int(0))}: copyPUint8ToPInt,
		// uint16 to int
		{Src: typeOf(uint16(0)), Dest: typeOf(int(0))}:               copyUint16ToInt,
		{Src: typeOfPointer(uint16(0)), Dest: typeOf(int(0))}:        copyPUint16ToInt,
		{Src: typeOf(uint16(0)), Dest: typeOfPointer(int(0))}:        copyUint16ToPInt,
		{Src: typeOfPointer(uint16(0)), Dest: typeOfPointer(int(0))}: copyPUint16ToPInt,
		// uint32 to int
		{Src: typeOf(uint32(0)), Dest: typeOf(int(0))}:               copyUint32ToInt,
		{Src: typeOfPointer(uint32(0)), Dest: typeOf(int(0))}:        copyPUint32ToInt,
		{Src: typeOf(uint32(0)), Dest: typeOfPointer(int(0))}:        copyUint32ToPInt,
		{Src: typeOfPointer(uint32(0)), Dest: typeOfPointer(int(0))}: copyPUint32ToPInt,
		// uint64 to int
		{Src: typeOf(uint64(0)), Dest: typeOf(int(0))}:               copyUint64ToInt,
		{Src: typeOfPointer(uint64(0)), Dest: typeOf(int(0))}:        copyPUint64ToInt,
		{Src: typeOf(uint64(0)), Dest: typeOfPointer(int(0))}:        copyUint64ToPInt,
		{Src: typeOfPointer(uint64(0)), Dest: typeOfPointer(int(0))}: copyPUint64ToPInt,
		// int to int8
		{Src: typeOf(int(0)), Dest: typeOf(int8(0))}:               copyIntToInt8,
		{Src: typeOfPointer(int(0)), Dest: typeOf(int8(0))}:        copyPIntToInt8,
		{Src: typeOf(int(0)), Dest: typeOfPointer(int8(0))}:        copyIntToPInt8,
		{Src: typeOfPointer(int(0)), Dest: typeOfPointer(int8(0))}: copyPIntToPInt8,
		// int8 to int8
		{Src: typeOf(int8(0)), Dest: typeOf(int8(0))}:               copyInt8ToInt8,
		{Src: typeOfPointer(int8(0)), Dest: typeOf(int8(0))}:        copyPInt8ToInt8,
		{Src: typeOf(int8(0)), Dest: typeOfPointer(int8(0))}:        copyInt8ToPInt8,
		{Src: typeOfPointer(int8(0)), Dest: typeOfPointer(int8(0))}: copyPInt8ToPInt8,
		// int16 to int8
		{Src: typeOf(int16(0)), Dest: typeOf(int8(0))}:               copyInt16ToInt8,
		{Src: typeOfPointer(int16(0)), Dest: typeOf(int8(0))}:        copyPInt16ToInt8,
		{Src: typeOf(int16(0)), Dest: typeOfPointer(int8(0))}:        copyInt16ToPInt8,
		{Src: typeOfPointer(int16(0)), Dest: typeOfPointer(int8(0))}: copyPInt16ToPInt8,
		// int32 to int8
		{Src: typeOf(int32(0)), Dest: typeOf(int8(0))}:               copyInt32ToInt8,
		{Src: typeOfPointer(int32(0)), Dest: typeOf(int8(0))}:        copyPInt32ToInt8,
		{Src: typeOf(int32(0)), Dest: typeOfPointer(int8(0))}:        copyInt32ToPInt8,
		{Src: typeOfPointer(int32(0)), Dest: typeOfPointer(int8(0))}: copyPInt32ToPInt8,
		// int64 to int8
		{Src: typeOf(int64(0)), Dest: typeOf(int8(0))}:               copyInt64ToInt8,
		{Src: typeOfPointer(int64(0)), Dest: typeOf(int8(0))}:        copyPInt64ToInt8,
		{Src: typeOf(int64(0)), Dest: typeOfPointer(int8(0))}:        copyInt64ToPInt8,
		{Src: typeOfPointer(int64(0)), Dest: typeOfPointer(int8(0))}: copyPInt64ToPInt8,
		// uint to int8
		{Src: typeOf(uint(0)), Dest: typeOf(int8(0))}:               copyUintToInt8,
		{Src: typeOfPointer(uint(0)), Dest: typeOf(int8(0))}:        copyPUintToInt8,
		{Src: typeOf(uint(0)), Dest: typeOfPointer(int8(0))}:        copyUintToPInt8,
		{Src: typeOfPointer(uint(0)), Dest: typeOfPointer(int8(0))}: copyPUintToPInt8,
		// uint8 to int8
		{Src: typeOf(uint8(0)), Dest: typeOf(int8(0))}:               copyUint8ToInt8,
		{Src: typeOfPointer(uint8(0)), Dest: typeOf(int8(0))}:        copyPUint8ToInt8,
		{Src: typeOf(uint8(0)), Dest: typeOfPointer(int8(0))}:        copyUint8ToPInt8,
		{Src: typeOfPointer(uint8(0)), Dest: typeOfPointer(int8(0))}: copyPUint8ToPInt8,
		// uint16 to int8
		{Src: typeOf(uint16(0)), Dest: typeOf(int8(0))}:               copyUint16ToInt8,
		{Src: typeOfPointer(uint16(0)), Dest: typeOf(int8(0))}:        copyPUint16ToInt8,
		{Src: typeOf(uint16(0)), Dest: typeOfPointer(int8(0))}:        copyUint16ToPInt8,
		{Src: typeOfPointer(uint16(0)), Dest: typeOfPointer(int8(0))}: copyPUint16ToPInt8,
		// uint32 to int8
		{Src: typeOf(uint32(0)), Dest: typeOf(int8(0))}:               copyUint32ToInt8,
		{Src: typeOfPointer(uint32(0)), Dest: typeOf(int8(0))}:        copyPUint32ToInt8,
		{Src: typeOf(uint32(0)), Dest: typeOfPointer(int8(0))}:        copyUint32ToPInt8,
		{Src: typeOfPointer(uint32(0)), Dest: typeOfPointer(int8(0))}: copyPUint32ToPInt8,
		// uint64 to int8
		{Src: typeOf(uint64(0)), Dest: typeOf(int8(0))}:               copyUint64ToInt8,
		{Src: typeOfPointer(uint64(0)), Dest: typeOf(int8(0))}:        copyPUint64ToInt8,
		{Src: typeOf(uint64(0)), Dest: typeOfPointer(int8(0))}:        copyUint64ToPInt8,
		{Src: typeOfPointer(uint64(0)), Dest: typeOfPointer(int8(0))}: copyPUint64ToPInt8,
		// int to int16
		{Src: typeOf(int(0)), Dest: typeOf(int16(0))}:               copyIntToInt16,
		{Src: typeOfPointer(int(0)), Dest: typeOf(int16(0))}:        copyPIntToInt16,
		{Src: typeOf(int(0)), Dest: typeOfPointer(int16(0))}:        copyIntToPInt16,
		{Src: typeOfPointer(int(0)), Dest: typeOfPointer(int16(0))}: copyPIntToPInt16,
		// int8 to int16
		{Src: typeOf(int8(0)), Dest: typeOf(int16(0))}:               copyInt8ToInt16,
		{Src: typeOfPointer(int8(0)), Dest: typeOf(int16(0))}:        copyPInt8ToInt16,
		{Src: typeOf(int8(0)), Dest: typeOfPointer(int16(0))}:        copyInt8ToPInt16,
		{Src: typeOfPointer(int8(0)), Dest: typeOfPointer(int16(0))}: copyPInt8ToPInt16,
		// int16 to int16
		{Src: typeOf(int16(0)), Dest: typeOf(int16(0))}:               copyInt16ToInt16,
		{Src: typeOfPointer(int16(0)), Dest: typeOf(int16(0))}:        copyPInt16ToInt16,
		{Src: typeOf(int16(0)), Dest: typeOfPointer(int16(0))}:        copyInt16ToPInt16,
		{Src: typeOfPointer(int16(0)), Dest: typeOfPointer(int16(0))}: copyPInt16ToPInt16,
		// int32 to int16
		{Src: typeOf(int32(0)), Dest: typeOf(int16(0))}:               copyInt32ToInt16,
		{Src: typeOfPointer(int32(0)), Dest: typeOf(int16(0))}:        copyPInt32ToInt16,
		{Src: typeOf(int32(0)), Dest: typeOfPointer(int16(0))}:        copyInt32ToPInt16,
		{Src: typeOfPointer(int32(0)), Dest: typeOfPointer(int16(0))}: copyPInt32ToPInt16,
		// int64 to int16
		{Src: typeOf(int64(0)), Dest: typeOf(int16(0))}:               copyInt64ToInt16,
		{Src: typeOfPointer(int64(0)), Dest: typeOf(int16(0))}:        copyPInt64ToInt16,
		{Src: typeOf(int64(0)), Dest: typeOfPointer(int16(0))}:        copyInt64ToPInt16,
		{Src: typeOfPointer(int64(0)), Dest: typeOfPointer(int16(0))}: copyPInt64ToPInt16,
		// uint to int16
		{Src: typeOf(uint(0)), Dest: typeOf(int16(0))}:               copyUintToInt16,
		{Src: typeOfPointer(uint(0)), Dest: typeOf(int16(0))}:        copyPUintToInt16,
		{Src: typeOf(uint(0)), Dest: typeOfPointer(int16(0))}:        copyUintToPInt16,
		{Src: typeOfPointer(uint(0)), Dest: typeOfPointer(int16(0))}: copyPUintToPInt16,
		// uint8 to int16
		{Src: typeOf(uint8(0)), Dest: typeOf(int16(0))}:               copyUint8ToInt16,
		{Src: typeOfPointer(uint8(0)), Dest: typeOf(int16(0))}:        copyPUint8ToInt16,
		{Src: typeOf(uint8(0)), Dest: typeOfPointer(int16(0))}:        copyUint8ToPInt16,
		{Src: typeOfPointer(uint8(0)), Dest: typeOfPointer(int16(0))}: copyPUint8ToPInt16,
		// uint16 to int16
		{Src: typeOf(uint16(0)), Dest: typeOf(int16(0))}:               copyUint16ToInt16,
		{Src: typeOfPointer(uint16(0)), Dest: typeOf(int16(0))}:        copyPUint16ToInt16,
		{Src: typeOf(uint16(0)), Dest: typeOfPointer(int16(0))}:        copyUint16ToPInt16,
		{Src: typeOfPointer(uint16(0)), Dest: typeOfPointer(int16(0))}: copyPUint16ToPInt16,
		// uint32 to int16
		{Src: typeOf(uint32(0)), Dest: typeOf(int16(0))}:               copyUint32ToInt16,
		{Src: typeOfPointer(uint32(0)), Dest: typeOf(int16(0))}:        copyPUint32ToInt16,
		{Src: typeOf(uint32(0)), Dest: typeOfPointer(int16(0))}:        copyUint32ToPInt16,
		{Src: typeOfPointer(uint32(0)), Dest: typeOfPointer(int16(0))}: copyPUint32ToPInt16,
		// uint64 to int16
		{Src: typeOf(uint64(0)), Dest: typeOf(int16(0))}:               copyUint64ToInt16,
		{Src: typeOfPointer(uint64(0)), Dest: typeOf(int16(0))}:        copyPUint64ToInt16,
		{Src: typeOf(uint64(0)), Dest: typeOfPointer(int16(0))}:        copyUint64ToPInt16,
		{Src: typeOfPointer(uint64(0)), Dest: typeOfPointer(int16(0))}: copyPUint64ToPInt16,
		// int to int32
		{Src: typeOf(int(0)), Dest: typeOf(int32(0))}:               copyIntToInt32,
		{Src: typeOfPointer(int(0)), Dest: typeOf(int32(0))}:        copyPIntToInt32,
		{Src: typeOf(int(0)), Dest: typeOfPointer(int32(0))}:        copyIntToPInt32,
		{Src: typeOfPointer(int(0)), Dest: typeOfPointer(int32(0))}: copyPIntToPInt32,
		// int8 to int32
		{Src: typeOf(int8(0)), Dest: typeOf(int32(0))}:               copyInt8ToInt32,
		{Src: typeOfPointer(int8(0)), Dest: typeOf(int32(0))}:        copyPInt8ToInt32,
		{Src: typeOf(int8(0)), Dest: typeOfPointer(int32(0))}:        copyInt8ToPInt32,
		{Src: typeOfPointer(int8(0)), Dest: typeOfPointer(int32(0))}: copyPInt8ToPInt32,
		// int16 to int32
		{Src: typeOf(int16(0)), Dest: typeOf(int32(0))}:               copyInt16ToInt32,
		{Src: typeOfPointer(int16(0)), Dest: typeOf(int32(0))}:        copyPInt16ToInt32,
		{Src: typeOf(int16(0)), Dest: typeOfPointer(int32(0))}:        copyInt16ToPInt32,
		{Src: typeOfPointer(int16(0)), Dest: typeOfPointer(int32(0))}: copyPInt16ToPInt32,
		// int32 to int32
		{Src: typeOf(int32(0)), Dest: typeOf(int32(0))}:               copyInt32ToInt32,
		{Src: typeOfPointer(int32(0)), Dest: typeOf(int32(0))}:        copyPInt32ToInt32,
		{Src: typeOf(int32(0)), Dest: typeOfPointer(int32(0))}:        copyInt32ToPInt32,
		{Src: typeOfPointer(int32(0)), Dest: typeOfPointer(int32(0))}: copyPInt32ToPInt32,
		// int64 to int32
		{Src: typeOf(int64(0)), Dest: typeOf(int32(0))}:               copyInt64ToInt32,
		{Src: typeOfPointer(int64(0)), Dest: typeOf(int32(0))}:        copyPInt64ToInt32,
		{Src: typeOf(int64(0)), Dest: typeOfPointer(int32(0))}:        copyInt64ToPInt32,
		{Src: typeOfPointer(int64(0)), Dest: typeOfPointer(int32(0))}: copyPInt64ToPInt32,
		// uint to int32
		{Src: typeOf(uint(0)), Dest: typeOf(int32(0))}:               copyUintToInt32,
		{Src: typeOfPointer(uint(0)), Dest: typeOf(int32(0))}:        copyPUintToInt32,
		{Src: typeOf(uint(0)), Dest: typeOfPointer(int32(0))}:        copyUintToPInt32,
		{Src: typeOfPointer(uint(0)), Dest: typeOfPointer(int32(0))}: copyPUintToPInt32,
		// uint8 to int32
		{Src: typeOf(uint8(0)), Dest: typeOf(int32(0))}:               copyUint8ToInt32,
		{Src: typeOfPointer(uint8(0)), Dest: typeOf(int32(0))}:        copyPUint8ToInt32,
		{Src: typeOf(uint8(0)), Dest: typeOfPointer(int32(0))}:        copyUint8ToPInt32,
		{Src: typeOfPointer(uint8(0)), Dest: typeOfPointer(int32(0))}: copyPUint8ToPInt32,
		// uint16 to int32
		{Src: typeOf(uint16(0)), Dest: typeOf(int32(0))}:               copyUint16ToInt32,
		{Src: typeOfPointer(uint16(0)), Dest: typeOf(int32(0))}:        copyPUint16ToInt32,
		{Src: typeOf(uint16(0)), Dest: typeOfPointer(int32(0))}:        copyUint16ToPInt32,
		{Src: typeOfPointer(uint16(0)), Dest: typeOfPointer(int32(0))}: copyPUint16ToPInt32,
		// uint32 to int32
		{Src: typeOf(uint32(0)), Dest: typeOf(int32(0))}:               copyUint32ToInt32,
		{Src: typeOfPointer(uint32(0)), Dest: typeOf(int32(0))}:        copyPUint32ToInt32,
		{Src: typeOf(uint32(0)), Dest: typeOfPointer(int32(0))}:        copyUint32ToPInt32,
		{Src: typeOfPointer(uint32(0)), Dest: typeOfPointer(int32(0))}: copyPUint32ToPInt32,
		// uint64 to int32
		{Src: typeOf(uint64(0)), Dest: typeOf(int32(0))}:               copyUint64ToInt32,
		{Src: typeOfPointer(uint64(0)), Dest: typeOf(int32(0))}:        copyPUint64ToInt32,
		{Src: typeOf(uint64(0)), Dest: typeOfPointer(int32(0))}:        copyUint64ToPInt32,
		{Src: typeOfPointer(uint64(0)), Dest: typeOfPointer(int32(0))}: copyPUint64ToPInt32,
		// int to int64
		{Src: typeOf(int(0)), Dest: typeOf(int64(0))}:               copyIntToInt64,
		{Src: typeOfPointer(int(0)), Dest: typeOf(int64(0))}:        copyPIntToInt64,
		{Src: typeOf(int(0)), Dest: typeOfPointer(int64(0))}:        copyIntToPInt64,
		{Src: typeOfPointer(int(0)), Dest: typeOfPointer(int64(0))}: copyPIntToPInt64,
		// int8 to int64
		{Src: typeOf(int8(0)), Dest: typeOf(int64(0))}:               copyInt8ToInt64,
		{Src: typeOfPointer(int8(0)), Dest: typeOf(int64(0))}:        copyPInt8ToInt64,
		{Src: typeOf(int8(0)), Dest: typeOfPointer(int64(0))}:        copyInt8ToPInt64,
		{Src: typeOfPointer(int8(0)), Dest: typeOfPointer(int64(0))}: copyPInt8ToPInt64,
		// int16 to int64
		{Src: typeOf(int16(0)), Dest: typeOf(int64(0))}:               copyInt16ToInt64,
		{Src: typeOfPointer(int16(0)), Dest: typeOf(int64(0))}:        copyPInt16ToInt64,
		{Src: typeOf(int16(0)), Dest: typeOfPointer(int64(0))}:        copyInt16ToPInt64,
		{Src: typeOfPointer(int16(0)), Dest: typeOfPointer(int64(0))}: copyPInt16ToPInt64,
		// int32 to int64
		{Src: typeOf(int32(0)), Dest: typeOf(int64(0))}:               copyInt32ToInt64,
		{Src: typeOfPointer(int32(0)), Dest: typeOf(int64(0))}:        copyPInt32ToInt64,
		{Src: typeOf(int32(0)), Dest: typeOfPointer(int64(0))}:        copyInt32ToPInt64,
		{Src: typeOfPointer(int32(0)), Dest: typeOfPointer(int64(0))}: copyPInt32ToPInt64,
		// int64 to int64
		{Src: typeOf(int64(0)), Dest: typeOf(int64(0))}:               copyInt64ToInt64,
		{Src: typeOfPointer(int64(0)), Dest: typeOf(int64(0))}:        copyPInt64ToInt64,
		{Src: typeOf(int64(0)), Dest: typeOfPointer(int64(0))}:        copyInt64ToPInt64,
		{Src: typeOfPointer(int64(0)), Dest: typeOfPointer(int64(0))}: copyPInt64ToPInt64,
		// uint to int64
		{Src: typeOf(uint(0)), Dest: typeOf(int64(0))}:               copyUintToInt64,
		{Src: typeOfPointer(uint(0)), Dest: typeOf(int64(0))}:        copyPUintToInt64,
		{Src: typeOf(uint(0)), Dest: typeOfPointer(int64(0))}:        copyUintToPInt64,
		{Src: typeOfPointer(uint(0)), Dest: typeOfPointer(int64(0))}: copyPUintToPInt64,
		// uint8 to int64
		{Src: typeOf(uint8(0)), Dest: typeOf(int64(0))}:               copyUint8ToInt64,
		{Src: typeOfPointer(uint8(0)), Dest: typeOf(int64(0))}:        copyPUint8ToInt64,
		{Src: typeOf(uint8(0)), Dest: typeOfPointer(int64(0))}:        copyUint8ToPInt64,
		{Src: typeOfPointer(uint8(0)), Dest: typeOfPointer(int64(0))}: copyPUint8ToPInt64,
		// uint16 to int64
		{Src: typeOf(uint16(0)), Dest: typeOf(int64(0))}:               copyUint16ToInt64,
		{Src: typeOfPointer(uint16(0)), Dest: typeOf(int64(0))}:        copyPUint16ToInt64,
		{Src: typeOf(uint16(0)), Dest: typeOfPointer(int64(0))}:        copyUint16ToPInt64,
		{Src: typeOfPointer(uint16(0)), Dest: typeOfPointer(int64(0))}: copyPUint16ToPInt64,
		// uint32 to int64
		{Src: typeOf(uint32(0)), Dest: typeOf(int64(0))}:               copyUint32ToInt64,
		{Src: typeOfPointer(uint32(0)), Dest: typeOf(int64(0))}:        copyPUint32ToInt64,
		{Src: typeOf(uint32(0)), Dest: typeOfPointer(int64(0))}:        copyUint32ToPInt64,
		{Src: typeOfPointer(uint32(0)), Dest: typeOfPointer(int64(0))}: copyPUint32ToPInt64,
		// uint64 to int64
		{Src: typeOf(uint64(0)), Dest: typeOf(int64(0))}:               copyUint64ToInt64,
		{Src: typeOfPointer(uint64(0)), Dest: typeOf(int64(0))}:        copyPUint64ToInt64,
		{Src: typeOf(uint64(0)), Dest: typeOfPointer(int64(0))}:        copyUint64ToPInt64,
		{Src: typeOfPointer(uint64(0)), Dest: typeOfPointer(int64(0))}: copyPUint64ToPInt64,
		// int to uint
		{Src: typeOf(int(0)), Dest: typeOf(uint(0))}:               copyIntToUint,
		{Src: typeOfPointer(int(0)), Dest: typeOf(uint(0))}:        copyPIntToUint,
		{Src: typeOf(int(0)), Dest: typeOfPointer(uint(0))}:        copyIntToPUint,
		{Src: typeOfPointer(int(0)), Dest: typeOfPointer(uint(0))}: copyPIntToPUint,
		// int8 to uint
		{Src: typeOf(int8(0)), Dest: typeOf(uint(0))}:               copyInt8ToUint,
		{Src: typeOfPointer(int8(0)), Dest: typeOf(uint(0))}:        copyPInt8ToUint,
		{Src: typeOf(int8(0)), Dest: typeOfPointer(uint(0))}:        copyInt8ToPUint,
		{Src: typeOfPointer(int8(0)), Dest: typeOfPointer(uint(0))}: copyPInt8ToPUint,
		// int16 to uint
		{Src: typeOf(int16(0)), Dest: typeOf(uint(0))}:               copyInt16ToUint,
		{Src: typeOfPointer(int16(0)), Dest: typeOf(uint(0))}:        copyPInt16ToUint,
		{Src: typeOf(int16(0)), Dest: typeOfPointer(uint(0))}:        copyInt16ToPUint,
		{Src: typeOfPointer(int16(0)), Dest: typeOfPointer(uint(0))}: copyPInt16ToPUint,
		// int32 to uint
		{Src: typeOf(int32(0)), Dest: typeOf(uint(0))}:               copyInt32ToUint,
		{Src: typeOfPointer(int32(0)), Dest: typeOf(uint(0))}:        copyPInt32ToUint,
		{Src: typeOf(int32(0)), Dest: typeOfPointer(uint(0))}:        copyInt32ToPUint,
		{Src: typeOfPointer(int32(0)), Dest: typeOfPointer(uint(0))}: copyPInt32ToPUint,
		// int64 to uint
		{Src: typeOf(int64(0)), Dest: typeOf(uint(0))}:               copyInt64ToUint,
		{Src: typeOfPointer(int64(0)), Dest: typeOf(uint(0))}:        copyPInt64ToUint,
		{Src: typeOf(int64(0)), Dest: typeOfPointer(uint(0))}:        copyInt64ToPUint,
		{Src: typeOfPointer(int64(0)), Dest: typeOfPointer(uint(0))}: copyPInt64ToPUint,
		// uint to uint
		{Src: typeOf(uint(0)), Dest: typeOf(uint(0))}:               copyUintToUint,
		{Src: typeOfPointer(uint(0)), Dest: typeOf(uint(0))}:        copyPUintToUint,
		{Src: typeOf(uint(0)), Dest: typeOfPointer(uint(0))}:        copyUintToPUint,
		{Src: typeOfPointer(uint(0)), Dest: typeOfPointer(uint(0))}: copyPUintToPUint,
		// uint8 to uint
		{Src: typeOf(uint8(0)), Dest: typeOf(uint(0))}:               copyUint8ToUint,
		{Src: typeOfPointer(uint8(0)), Dest: typeOf(uint(0))}:        copyPUint8ToUint,
		{Src: typeOf(uint8(0)), Dest: typeOfPointer(uint(0))}:        copyUint8ToPUint,
		{Src: typeOfPointer(uint8(0)), Dest: typeOfPointer(uint(0))}: copyPUint8ToPUint,
		// uint16 to uint
		{Src: typeOf(uint16(0)), Dest: typeOf(uint(0))}:               copyUint16ToUint,
		{Src: typeOfPointer(uint16(0)), Dest: typeOf(uint(0))}:        copyPUint16ToUint,
		{Src: typeOf(uint16(0)), Dest: typeOfPointer(uint(0))}:        copyUint16ToPUint,
		{Src: typeOfPointer(uint16(0)), Dest: typeOfPointer(uint(0))}: copyPUint16ToPUint,
		// uint32 to uint
		{Src: typeOf(uint32(0)), Dest: typeOf(uint(0))}:               copyUint32ToUint,
		{Src: typeOfPointer(uint32(0)), Dest: typeOf(uint(0))}:        copyPUint32ToUint,
		{Src: typeOf(uint32(0)), Dest: typeOfPointer(uint(0))}:        copyUint32ToPUint,
		{Src: typeOfPointer(uint32(0)), Dest: typeOfPointer(uint(0))}: copyPUint32ToPUint,
		// uint64 to uint
		{Src: typeOf(uint64(0)), Dest: typeOf(uint(0))}:               copyUint64ToUint,
		{Src: typeOfPointer(uint64(0)), Dest: typeOf(uint(0))}:        copyPUint64ToUint,
		{Src: typeOf(uint64(0)), Dest: typeOfPointer(uint(0))}:        copyUint64ToPUint,
		{Src: typeOfPointer(uint64(0)), Dest: typeOfPointer(uint(0))}: copyPUint64ToPUint,
		// int to uint8
		{Src: typeOf(int(0)), Dest: typeOf(uint8(0))}:               copyIntToUint8,
		{Src: typeOfPointer(int(0)), Dest: typeOf(uint8(0))}:        copyPIntToUint8,
		{Src: typeOf(int(0)), Dest: typeOfPointer(uint8(0))}:        copyIntToPUint8,
		{Src: typeOfPointer(int(0)), Dest: typeOfPointer(uint8(0))}: copyPIntToPUint8,
		// int8 to uint8
		{Src: typeOf(int8(0)), Dest: typeOf(uint8(0))}:               copyInt8ToUint8,
		{Src: typeOfPointer(int8(0)), Dest: typeOf(uint8(0))}:        copyPInt8ToUint8,
		{Src: typeOf(int8(0)), Dest: typeOfPointer(uint8(0))}:        copyInt8ToPUint8,
		{Src: typeOfPointer(int8(0)), Dest: typeOfPointer(uint8(0))}: copyPInt8ToPUint8,
		// int16 to uint8
		{Src: typeOf(int16(0)), Dest: typeOf(uint8(0))}:               copyInt16ToUint8,
		{Src: typeOfPointer(int16(0)), Dest: typeOf(uint8(0))}:        copyPInt16ToUint8,
		{Src: typeOf(int16(0)), Dest: typeOfPointer(uint8(0))}:        copyInt16ToPUint8,
		{Src: typeOfPointer(int16(0)), Dest: typeOfPointer(uint8(0))}: copyPInt16ToPUint8,
		// int32 to uint8
		{Src: typeOf(int32(0)), Dest: typeOf(uint8(0))}:               copyInt32ToUint8,
		{Src: typeOfPointer(int32(0)), Dest: typeOf(uint8(0))}:        copyPInt32ToUint8,
		{Src: typeOf(int32(0)), Dest: typeOfPointer(uint8(0))}:        copyInt32ToPUint8,
		{Src: typeOfPointer(int32(0)), Dest: typeOfPointer(uint8(0))}: copyPInt32ToPUint8,
		// int64 to uint8
		{Src: typeOf(int64(0)), Dest: typeOf(uint8(0))}:               copyInt64ToUint8,
		{Src: typeOfPointer(int64(0)), Dest: typeOf(uint8(0))}:        copyPInt64ToUint8,
		{Src: typeOf(int64(0)), Dest: typeOfPointer(uint8(0))}:        copyInt64ToPUint8,
		{Src: typeOfPointer(int64(0)), Dest: typeOfPointer(uint8(0))}: copyPInt64ToPUint8,
		// uint to uint8
		{Src: typeOf(uint(0)), Dest: typeOf(uint8(0))}:               copyUintToUint8,
		{Src: typeOfPointer(uint(0)), Dest: typeOf(uint8(0))}:        copyPUintToUint8,
		{Src: typeOf(uint(0)), Dest: typeOfPointer(uint8(0))}:        copyUintToPUint8,
		{Src: typeOfPointer(uint(0)), Dest: typeOfPointer(uint8(0))}: copyPUintToPUint8,
		// uint8 to uint8
		{Src: typeOf(uint8(0)), Dest: typeOf(uint8(0))}:               copyUint8ToUint8,
		{Src: typeOfPointer(uint8(0)), Dest: typeOf(uint8(0))}:        copyPUint8ToUint8,
		{Src: typeOf(uint8(0)), Dest: typeOfPointer(uint8(0))}:        copyUint8ToPUint8,
		{Src: typeOfPointer(uint8(0)), Dest: typeOfPointer(uint8(0))}: copyPUint8ToPUint8,
		// uint16 to uint8
		{Src: typeOf(uint16(0)), Dest: typeOf(uint8(0))}:               copyUint16ToUint8,
		{Src: typeOfPointer(uint16(0)), Dest: typeOf(uint8(0))}:        copyPUint16ToUint8,
		{Src: typeOf(uint16(0)), Dest: typeOfPointer(uint8(0))}:        copyUint16ToPUint8,
		{Src: typeOfPointer(uint16(0)), Dest: typeOfPointer(uint8(0))}: copyPUint16ToPUint8,
		// uint32 to uint8
		{Src: typeOf(uint32(0)), Dest: typeOf(uint8(0))}:               copyUint32ToUint8,
		{Src: typeOfPointer(uint32(0)), Dest: typeOf(uint8(0))}:        copyPUint32ToUint8,
		{Src: typeOf(uint32(0)), Dest: typeOfPointer(uint8(0))}:        copyUint32ToPUint8,
		{Src: typeOfPointer(uint32(0)), Dest: typeOfPointer(uint8(0))}: copyPUint32ToPUint8,
		// uint64 to uint8
		{Src: typeOf(uint64(0)), Dest: typeOf(uint8(0))}:               copyUint64ToUint8,
		{Src: typeOfPointer(uint64(0)), Dest: typeOf(uint8(0))}:        copyPUint64ToUint8,
		{Src: typeOf(uint64(0)), Dest: typeOfPointer(uint8(0))}:        copyUint64ToPUint8,
		{Src: typeOfPointer(uint64(0)), Dest: typeOfPointer(uint8(0))}: copyPUint64ToPUint8,
		// int to uint16
		{Src: typeOf(int(0)), Dest: typeOf(uint16(0))}:               copyIntToUint16,
		{Src: typeOfPointer(int(0)), Dest: typeOf(uint16(0))}:        copyPIntToUint16,
		{Src: typeOf(int(0)), Dest: typeOfPointer(uint16(0))}:        copyIntToPUint16,
		{Src: typeOfPointer(int(0)), Dest: typeOfPointer(uint16(0))}: copyPIntToPUint16,
		// int8 to uint16
		{Src: typeOf(int8(0)), Dest: typeOf(uint16(0))}:               copyInt8ToUint16,
		{Src: typeOfPointer(int8(0)), Dest: typeOf(uint16(0))}:        copyPInt8ToUint16,
		{Src: typeOf(int8(0)), Dest: typeOfPointer(uint16(0))}:        copyInt8ToPUint16,
		{Src: typeOfPointer(int8(0)), Dest: typeOfPointer(uint16(0))}: copyPInt8ToPUint16,
		// int16 to uint16
		{Src: typeOf(int16(0)), Dest: typeOf(uint16(0))}:               copyInt16ToUint16,
		{Src: typeOfPointer(int16(0)), Dest: typeOf(uint16(0))}:        copyPInt16ToUint16,
		{Src: typeOf(int16(0)), Dest: typeOfPointer(uint16(0))}:        copyInt16ToPUint16,
		{Src: typeOfPointer(int16(0)), Dest: typeOfPointer(uint16(0))}: copyPInt16ToPUint16,
		// int32 to uint16
		{Src: typeOf(int32(0)), Dest: typeOf(uint16(0))}:               copyInt32ToUint16,
		{Src: typeOfPointer(int32(0)), Dest: typeOf(uint16(0))}:        copyPInt32ToUint16,
		{Src: typeOf(int32(0)), Dest: typeOfPointer(uint16(0))}:        copyInt32ToPUint16,
		{Src: typeOfPointer(int32(0)), Dest: typeOfPointer(uint16(0))}: copyPInt32ToPUint16,
		// int64 to uint16
		{Src: typeOf(int64(0)), Dest: typeOf(uint16(0))}:               copyInt64ToUint16,
		{Src: typeOfPointer(int64(0)), Dest: typeOf(uint16(0))}:        copyPInt64ToUint16,
		{Src: typeOf(int64(0)), Dest: typeOfPointer(uint16(0))}:        copyInt64ToPUint16,
		{Src: typeOfPointer(int64(0)), Dest: typeOfPointer(uint16(0))}: copyPInt64ToPUint16,
		// uint to uint16
		{Src: typeOf(uint(0)), Dest: typeOf(uint16(0))}:               copyUintToUint16,
		{Src: typeOfPointer(uint(0)), Dest: typeOf(uint16(0))}:        copyPUintToUint16,
		{Src: typeOf(uint(0)), Dest: typeOfPointer(uint16(0))}:        copyUintToPUint16,
		{Src: typeOfPointer(uint(0)), Dest: typeOfPointer(uint16(0))}: copyPUintToPUint16,
		// uint8 to uint16
		{Src: typeOf(uint8(0)), Dest: typeOf(uint16(0))}:               copyUint8ToUint16,
		{Src: typeOfPointer(uint8(0)), Dest: typeOf(uint16(0))}:        copyPUint8ToUint16,
		{Src: typeOf(uint8(0)), Dest: typeOfPointer(uint16(0))}:        copyUint8ToPUint16,
		{Src: typeOfPointer(uint8(0)), Dest: typeOfPointer(uint16(0))}: copyPUint8ToPUint16,
		// uint16 to uint16
		{Src: typeOf(uint16(0)), Dest: typeOf(uint16(0))}:               copyUint16ToUint16,
		{Src: typeOfPointer(uint16(0)), Dest: typeOf(uint16(0))}:        copyPUint16ToUint16,
		{Src: typeOf(uint16(0)), Dest: typeOfPointer(uint16(0))}:        copyUint16ToPUint16,
		{Src: typeOfPointer(uint16(0)), Dest: typeOfPointer(uint16(0))}: copyPUint16ToPUint16,
		// uint32 to uint16
		{Src: typeOf(uint32(0)), Dest: typeOf(uint16(0))}:               copyUint32ToUint16,
		{Src: typeOfPointer(uint32(0)), Dest: typeOf(uint16(0))}:        copyPUint32ToUint16,
		{Src: typeOf(uint32(0)), Dest: typeOfPointer(uint16(0))}:        copyUint32ToPUint16,
		{Src: typeOfPointer(uint32(0)), Dest: typeOfPointer(uint16(0))}: copyPUint32ToPUint16,
		// uint64 to uint16
		{Src: typeOf(uint64(0)), Dest: typeOf(uint16(0))}:               copyUint64ToUint16,
		{Src: typeOfPointer(uint64(0)), Dest: typeOf(uint16(0))}:        copyPUint64ToUint16,
		{Src: typeOf(uint64(0)), Dest: typeOfPointer(uint16(0))}:        copyUint64ToPUint16,
		{Src: typeOfPointer(uint64(0)), Dest: typeOfPointer(uint16(0))}: copyPUint64ToPUint16,
		// int to uint32
		{Src: typeOf(int(0)), Dest: typeOf(uint32(0))}:               copyIntToUint32,
		{Src: typeOfPointer(int(0)), Dest: typeOf(uint32(0))}:        copyPIntToUint32,
		{Src: typeOf(int(0)), Dest: typeOfPointer(uint32(0))}:        copyIntToPUint32,
		{Src: typeOfPointer(int(0)), Dest: typeOfPointer(uint32(0))}: copyPIntToPUint32,
		// int8 to uint32
		{Src: typeOf(int8(0)), Dest: typeOf(uint32(0))}:               copyInt8ToUint32,
		{Src: typeOfPointer(int8(0)), Dest: typeOf(uint32(0))}:        copyPInt8ToUint32,
		{Src: typeOf(int8(0)), Dest: typeOfPointer(uint32(0))}:        copyInt8ToPUint32,
		{Src: typeOfPointer(int8(0)), Dest: typeOfPointer(uint32(0))}: copyPInt8ToPUint32,
		// int16 to uint32
		{Src: typeOf(int16(0)), Dest: typeOf(uint32(0))}:               copyInt16ToUint32,
		{Src: typeOfPointer(int16(0)), Dest: typeOf(uint32(0))}:        copyPInt16ToUint32,
		{Src: typeOf(int16(0)), Dest: typeOfPointer(uint32(0))}:        copyInt16ToPUint32,
		{Src: typeOfPointer(int16(0)), Dest: typeOfPointer(uint32(0))}: copyPInt16ToPUint32,
		// int32 to uint32
		{Src: typeOf(int32(0)), Dest: typeOf(uint32(0))}:               copyInt32ToUint32,
		{Src: typeOfPointer(int32(0)), Dest: typeOf(uint32(0))}:        copyPInt32ToUint32,
		{Src: typeOf(int32(0)), Dest: typeOfPointer(uint32(0))}:        copyInt32ToPUint32,
		{Src: typeOfPointer(int32(0)), Dest: typeOfPointer(uint32(0))}: copyPInt32ToPUint32,
		// int64 to uint32
		{Src: typeOf(int64(0)), Dest: typeOf(uint32(0))}:               copyInt64ToUint32,
		{Src: typeOfPointer(int64(0)), Dest: typeOf(uint32(0))}:        copyPInt64ToUint32,
		{Src: typeOf(int64(0)), Dest: typeOfPointer(uint32(0))}:        copyInt64ToPUint32,
		{Src: typeOfPointer(int64(0)), Dest: typeOfPointer(uint32(0))}: copyPInt64ToPUint32,
		// uint to uint32
		{Src: typeOf(uint(0)), Dest: typeOf(uint32(0))}:               copyUintToUint32,
		{Src: typeOfPointer(uint(0)), Dest: typeOf(uint32(0))}:        copyPUintToUint32,
		{Src: typeOf(uint(0)), Dest: typeOfPointer(uint32(0))}:        copyUintToPUint32,
		{Src: typeOfPointer(uint(0)), Dest: typeOfPointer(uint32(0))}: copyPUintToPUint32,
		// uint8 to uint32
		{Src: typeOf(uint8(0)), Dest: typeOf(uint32(0))}:               copyUint8ToUint32,
		{Src: typeOfPointer(uint8(0)), Dest: typeOf(uint32(0))}:        copyPUint8ToUint32,
		{Src: typeOf(uint8(0)), Dest: typeOfPointer(uint32(0))}:        copyUint8ToPUint32,
		{Src: typeOfPointer(uint8(0)), Dest: typeOfPointer(uint32(0))}: copyPUint8ToPUint32,
		// uint16 to uint32
		{Src: typeOf(uint16(0)), Dest: typeOf(uint32(0))}:               copyUint16ToUint32,
		{Src: typeOfPointer(uint16(0)), Dest: typeOf(uint32(0))}:        copyPUint16ToUint32,
		{Src: typeOf(uint16(0)), Dest: typeOfPointer(uint32(0))}:        copyUint16ToPUint32,
		{Src: typeOfPointer(uint16(0)), Dest: typeOfPointer(uint32(0))}: copyPUint16ToPUint32,
		// uint32 to uint32
		{Src: typeOf(uint32(0)), Dest: typeOf(uint32(0))}:               copyUint32ToUint32,
		{Src: typeOfPointer(uint32(0)), Dest: typeOf(uint32(0))}:        copyPUint32ToUint32,
		{Src: typeOf(uint32(0)), Dest: typeOfPointer(uint32(0))}:        copyUint32ToPUint32,
		{Src: typeOfPointer(uint32(0)), Dest: typeOfPointer(uint32(0))}: copyPUint32ToPUint32,
		// uint64 to uint32
		{Src: typeOf(uint64(0)), Dest: typeOf(uint32(0))}:               copyUint64ToUint32,
		{Src: typeOfPointer(uint64(0)), Dest: typeOf(uint32(0))}:        copyPUint64ToUint32,
		{Src: typeOf(uint64(0)), Dest: typeOfPointer(uint32(0))}:        copyUint64ToPUint32,
		{Src: typeOfPointer(uint64(0)), Dest: typeOfPointer(uint32(0))}: copyPUint64ToPUint32,
		// int to uint64
		{Src: typeOf(int(0)), Dest: typeOf(uint64(0))}:               copyIntToUint64,
		{Src: typeOfPointer(int(0)), Dest: typeOf(uint64(0))}:        copyPIntToUint64,
		{Src: typeOf(int(0)), Dest: typeOfPointer(uint64(0))}:        copyIntToPUint64,
		{Src: typeOfPointer(int(0)), Dest: typeOfPointer(uint64(0))}: copyPIntToPUint64,
		// int8 to uint64
		{Src: typeOf(int8(0)), Dest: typeOf(uint64(0))}:               copyInt8ToUint64,
		{Src: typeOfPointer(int8(0)), Dest: typeOf(uint64(0))}:        copyPInt8ToUint64,
		{Src: typeOf(int8(0)), Dest: typeOfPointer(uint64(0))}:        copyInt8ToPUint64,
		{Src: typeOfPointer(int8(0)), Dest: typeOfPointer(uint64(0))}: copyPInt8ToPUint64,
		// int16 to uint64
		{Src: typeOf(int16(0)), Dest: typeOf(uint64(0))}:               copyInt16ToUint64,
		{Src: typeOfPointer(int16(0)), Dest: typeOf(uint64(0))}:        copyPInt16ToUint64,
		{Src: typeOf(int16(0)), Dest: typeOfPointer(uint64(0))}:        copyInt16ToPUint64,
		{Src: typeOfPointer(int16(0)), Dest: typeOfPointer(uint64(0))}: copyPInt16ToPUint64,
		// int32 to uint64
		{Src: typeOf(int32(0)), Dest: typeOf(uint64(0))}:               copyInt32ToUint64,
		{Src: typeOfPointer(int32(0)), Dest: typeOf(uint64(0))}:        copyPInt32ToUint64,
		{Src: typeOf(int32(0)), Dest: typeOfPointer(uint64(0))}:        copyInt32ToPUint64,
		{Src: typeOfPointer(int32(0)), Dest: typeOfPointer(uint64(0))}: copyPInt32ToPUint64,
		// int64 to uint64
		{Src: typeOf(int64(0)), Dest: typeOf(uint64(0))}:               copyInt64ToUint64,
		{Src: typeOfPointer(int64(0)), Dest: typeOf(uint64(0))}:        copyPInt64ToUint64,
		{Src: typeOf(int64(0)), Dest: typeOfPointer(uint64(0))}:        copyInt64ToPUint64,
		{Src: typeOfPointer(int64(0)), Dest: typeOfPointer(uint64(0))}: copyPInt64ToPUint64,
		// uint to uint64
		{Src: typeOf(uint(0)), Dest: typeOf(uint64(0))}:               copyUintToUint64,
		{Src: typeOfPointer(uint(0)), Dest: typeOf(uint64(0))}:        copyPUintToUint64,
		{Src: typeOf(uint(0)), Dest: typeOfPointer(uint64(0))}:        copyUintToPUint64,
		{Src: typeOfPointer(uint(0)), Dest: typeOfPointer(uint64(0))}: copyPUintToPUint64,
		// uint8 to uint64
		{Src: typeOf(uint8(0)), Dest: typeOf(uint64(0))}:               copyUint8ToUint64,
		{Src: typeOfPointer(uint8(0)), Dest: typeOf(uint64(0))}:        copyPUint8ToUint64,
		{Src: typeOf(uint8(0)), Dest: typeOfPointer(uint64(0))}:        copyUint8ToPUint64,
		{Src: typeOfPointer(uint8(0)), Dest: typeOfPointer(uint64(0))}: copyPUint8ToPUint64,
		// uint16 to uint64
		{Src: typeOf(uint16(0)), Dest: typeOf(uint64(0))}:               copyUint16ToUint64,
		{Src: typeOfPointer(uint16(0)), Dest: typeOf(uint64(0))}:        copyPUint16ToUint64,
		{Src: typeOf(uint16(0)), Dest: typeOfPointer(uint64(0))}:        copyUint16ToPUint64,
		{Src: typeOfPointer(uint16(0)), Dest: typeOfPointer(uint64(0))}: copyPUint16ToPUint64,
		// uint32 to uint64
		{Src: typeOf(uint32(0)), Dest: typeOf(uint64(0))}:               copyUint32ToUint64,
		{Src: typeOfPointer(uint32(0)), Dest: typeOf(uint64(0))}:        copyPUint32ToUint64,
		{Src: typeOf(uint32(0)), Dest: typeOfPointer(uint64(0))}:        copyUint32ToPUint64,
		{Src: typeOfPointer(uint32(0)), Dest: typeOfPointer(uint64(0))}: copyPUint32ToPUint64,
		// uint64 to uint64
		{Src: typeOf(uint64(0)), Dest: typeOf(uint64(0))}:               copyUint64ToUint64,
		{Src: typeOfPointer(uint64(0)), Dest: typeOf(uint64(0))}:        copyPUint64ToUint64,
		{Src: typeOf(uint64(0)), Dest: typeOfPointer(uint64(0))}:        copyUint64ToPUint64,
		{Src: typeOfPointer(uint64(0)), Dest: typeOfPointer(uint64(0))}: copyPUint64ToPUint64,
		// float32 to float32
		{Src: typeOf(float32(0)), Dest: typeOf(float32(0))}:               copyFloat32ToFloat32,
		{Src: typeOfPointer(float32(0)), Dest: typeOf(float32(0))}:        copyPFloat32ToFloat32,
		{Src: typeOf(float32(0)), Dest: typeOfPointer(float32(0))}:        copyFloat32ToPFloat32,
		{Src: typeOfPointer(float32(0)), Dest: typeOfPointer(float32(0))}: copyPFloat32ToPFloat32,
		// float64 to float32
		{Src: typeOf(float64(0)), Dest: typeOf(float32(0))}:               copyFloat64ToFloat32,
		{Src: typeOfPointer(float64(0)), Dest: typeOf(float32(0))}:        copyPFloat64ToFloat32,
		{Src: typeOf(float64(0)), Dest: typeOfPointer(float32(0))}:        copyFloat64ToPFloat32,
		{Src: typeOfPointer(float64(0)), Dest: typeOfPointer(float32(0))}: copyPFloat64ToPFloat32,
		// float32 to float64
		{Src: typeOf(float32(0)), Dest: typeOf(float64(0))}:               copyFloat32ToFloat64,
		{Src: typeOfPointer(float32(0)), Dest: typeOf(float64(0))}:        copyPFloat32ToFloat64,
		{Src: typeOf(float32(0)), Dest: typeOfPointer(float64(0))}:        copyFloat32ToPFloat64,
		{Src: typeOfPointer(float32(0)), Dest: typeOfPointer(float64(0))}: copyPFloat32ToPFloat64,
		// float64 to float64
		{Src: typeOf(float64(0)), Dest: typeOf(float64(0))}:               copyFloat64ToFloat64,
		{Src: typeOfPointer(float64(0)), Dest: typeOf(float64(0))}:        copyPFloat64ToFloat64,
		{Src: typeOf(float64(0)), Dest: typeOfPointer(float64(0))}:        copyFloat64ToPFloat64,
		{Src: typeOfPointer(float64(0)), Dest: typeOfPointer(float64(0))}: copyPFloat64ToPFloat64,
		// bool to bool
		{Src: typeOf(bool(false)), Dest: typeOf(bool(false))}:               copyBoolToBool,
		{Src: typeOfPointer(bool(false)), Dest: typeOf(bool(false))}:        copyPBoolToBool,
		{Src: typeOf(bool(false)), Dest: typeOfPointer(bool(false))}:        copyBoolToPBool,
		{Src: typeOfPointer(bool(false)), Dest: typeOfPointer(bool(false))}: copyPBoolToPBool,
		// complex64 to complex64
		{Src: typeOf(complex64(0)), Dest: typeOf(complex64(0))}:               copyComplex64ToComplex64,
		{Src: typeOfPointer(complex64(0)), Dest: typeOf(complex64(0))}:        copyPComplex64ToComplex64,
		{Src: typeOf(complex64(0)), Dest: typeOfPointer(complex64(0))}:        copyComplex64ToPComplex64,
		{Src: typeOfPointer(complex64(0)), Dest: typeOfPointer(complex64(0))}: copyPComplex64ToPComplex64,
		// complex128 to complex64
		{Src: typeOf(complex128(0)), Dest: typeOf(complex64(0))}:               copyComplex128ToComplex64,
		{Src: typeOfPointer(complex128(0)), Dest: typeOf(complex64(0))}:        copyPComplex128ToComplex64,
		{Src: typeOf(complex128(0)), Dest: typeOfPointer(complex64(0))}:        copyComplex128ToPComplex64,
		{Src: typeOfPointer(complex128(0)), Dest: typeOfPointer(complex64(0))}: copyPComplex128ToPComplex64,
		// complex64 to complex128
		{Src: typeOf(complex64(0)), Dest: typeOf(complex128(0))}:               copyComplex64ToComplex128,
		{Src: typeOfPointer(complex64(0)), Dest: typeOf(complex128(0))}:        copyPComplex64ToComplex128,
		{Src: typeOf(complex64(0)), Dest: typeOfPointer(complex128(0))}:        copyComplex64ToPComplex128,
		{Src: typeOfPointer(complex64(0)), Dest: typeOfPointer(complex128(0))}: copyPComplex64ToPComplex128,
		// complex128 to complex128
		{Src: typeOf(complex128(0)), Dest: typeOf(complex128(0))}:               copyComplex128ToComplex128,
		{Src: typeOfPointer(complex128(0)), Dest: typeOf(complex128(0))}:        copyPComplex128ToComplex128,
		{Src: typeOf(complex128(0)), Dest: typeOfPointer(complex128(0))}:        copyComplex128ToPComplex128,
		{Src: typeOfPointer(complex128(0)), Dest: typeOfPointer(complex128(0))}: copyPComplex128ToPComplex128,
		// string to string
		{Src: typeOf(string("")), Dest: typeOf(string(""))}:               copyStringToString,
		{Src: typeOfPointer(string("")), Dest: typeOf(string(""))}:        copyPStringToString,
		{Src: typeOf(string("")), Dest: typeOfPointer(string(""))}:        copyStringToPString,
		{Src: typeOfPointer(string("")), Dest: typeOfPointer(string(""))}: copyPStringToPString,
		// []byte to string
		{Src: typeOf([]byte(nil)), Dest: typeOf(string(""))}:               copyBytesToString,
		{Src: typeOfPointer([]byte(nil)), Dest: typeOf(string(""))}:        copyPBytesToString,
		{Src: typeOf([]byte(nil)), Dest: typeOfPointer(string(""))}:        copyBytesToPString,
		{Src: typeOfPointer([]byte(nil)), Dest: typeOfPointer(string(""))}: copyPBytesToPString,
		// string to []byte
		{Src: typeOf(string("")), Dest: typeOf([]byte(nil))}:               copyStringToBytes,
		{Src: typeOfPointer(string("")), Dest: typeOf([]byte(nil))}:        copyPStringToBytes,
		{Src: typeOf(string("")), Dest: typeOfPointer([]byte(nil))}:        copyStringToPBytes,
		{Src: typeOfPointer(string("")), Dest: typeOfPointer([]byte(nil))}: copyPStringToPBytes,
		// []byte to []byte
		{Src: typeOf([]byte(nil)), Dest: typeOf([]byte(nil))}:               copyBytesToBytes,
		{Src: typeOfPointer([]byte(nil)), Dest: typeOf([]byte(nil))}:        copyPBytesToBytes,
		{Src: typeOf([]byte(nil)), Dest: typeOfPointer([]byte(nil))}:        copyBytesToPBytes,
		{Src: typeOfPointer([]byte(nil)), Dest: typeOfPointer([]byte(nil))}: copyPBytesToPBytes,
		// time.Time to time.Time
		{Src: typeOf(time.Time(time.Time{})), Dest: typeOf(time.Time(time.Time{}))}:               copyTimeToTime,
		{Src: typeOfPointer(time.Time(time.Time{})), Dest: typeOf(time.Time(time.Time{}))}:        copyPTimeToTime,
		{Src: typeOf(time.Time(time.Time{})), Dest: typeOfPointer(time.Time(time.Time{}))}:        copyTimeToPTime,
		{Src: typeOfPointer(time.Time(time.Time{})), Dest: typeOfPointer(time.Time(time.Time{}))}: copyPTimeToPTime,
		// time.Duration to time.Duration
		{Src: typeOf(time.Duration(0)), Dest: typeOf(time.Duration(0))}:               copyDurationToDuration,
		{Src: typeOfPointer(time.Duration(0)), Dest: typeOf(time.Duration(0))}:        copyPDurationToDuration,
		{Src: typeOf(time.Duration(0)), Dest: typeOfPointer(time.Duration(0))}:        copyDurationToPDuration,
		{Src: typeOfPointer(time.Duration(0)), Dest: typeOfPointer(time.Duration(0))}: copyPDurationToPDuration,
	},
	sizes: []func(dst, src unsafe.Pointer){
		copy1, copy2, copy3, copy4, copy5, copy6, copy7, copy8, copy9, copy10, copy11, copy12, copy13, copy14, copy15, copy16, copy17, copy18, copy19, copy20, copy21, copy22, copy23, copy24, copy25, copy26, copy27, copy28, copy29, copy30, copy31, copy32, copy33, copy34, copy35, copy36, copy37, copy38, copy39, copy40, copy41, copy42, copy43, copy44, copy45, copy46, copy47, copy48, copy49, copy50, copy51, copy52, copy53, copy54, copy55, copy56, copy57, copy58, copy59, copy60, copy61, copy62, copy63, copy64, copy65, copy66, copy67, copy68, copy69, copy70, copy71, copy72, copy73, copy74, copy75, copy76, copy77, copy78, copy79, copy80, copy81, copy82, copy83, copy84, copy85, copy86, copy87, copy88, copy89, copy90, copy91, copy92, copy93, copy94, copy95, copy96, copy97, copy98, copy99, copy100, copy101, copy102, copy103, copy104, copy105, copy106, copy107, copy108, copy109, copy110, copy111, copy112, copy113, copy114, copy115, copy116, copy117, copy118, copy119, copy120, copy121, copy122, copy123, copy124, copy125, copy126, copy127, copy128, copy129, copy130, copy131, copy132, copy133, copy134, copy135, copy136, copy137, copy138, copy139, copy140, copy141, copy142, copy143, copy144, copy145, copy146, copy147, copy148, copy149, copy150, copy151, copy152, copy153, copy154, copy155, copy156, copy157, copy158, copy159, copy160, copy161, copy162, copy163, copy164, copy165, copy166, copy167, copy168, copy169, copy170, copy171, copy172, copy173, copy174, copy175, copy176, copy177, copy178, copy179, copy180, copy181, copy182, copy183, copy184, copy185, copy186, copy187, copy188, copy189, copy190, copy191, copy192, copy193, copy194, copy195, copy196, copy197, copy198, copy199, copy200, copy201, copy202, copy203, copy204, copy205, copy206, copy207, copy208, copy209, copy210, copy211, copy212, copy213, copy214, copy215, copy216, copy217, copy218, copy219, copy220, copy221, copy222, copy223, copy224, copy225, copy226, copy227, copy228, copy229, copy230, copy231, copy232, copy233, copy234, copy235, copy236, copy237, copy238, copy239, copy240, copy241, copy242, copy243, copy244, copy245, copy246, copy247, copy248, copy249, copy250, copy251, copy252, copy253, copy254, copy255,
	},
}

// int to int

func copyIntToInt(dst, src unsafe.Pointer) {
	*(*int)(unsafe.Pointer(dst)) = int(*(*int)(unsafe.Pointer(src)))
}

func copyPIntToInt(dst, src unsafe.Pointer) {
	var v int
	if p := *(**int)(unsafe.Pointer(src)); p != nil {
		v = int(*p)
	}
	*(*int)(unsafe.Pointer(dst)) = v
}

func copyIntToPInt(dst, src unsafe.Pointer) {
	v := int(*(*int)(unsafe.Pointer(src)))
	p := (**int)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPIntToPInt(dst, src unsafe.Pointer) {
	var v int
	if p := *(**int)(unsafe.Pointer(src)); p != nil {
		v = int(*p)
	}

	p := (**int)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int8 to int

func copyInt8ToInt(dst, src unsafe.Pointer) {
	*(*int)(unsafe.Pointer(dst)) = int(*(*int8)(unsafe.Pointer(src)))
}

func copyPInt8ToInt(dst, src unsafe.Pointer) {
	var v int
	if p := *(**int8)(unsafe.Pointer(src)); p != nil {
		v = int(*p)
	}
	*(*int)(unsafe.Pointer(dst)) = v
}

func copyInt8ToPInt(dst, src unsafe.Pointer) {
	v := int(*(*int8)(unsafe.Pointer(src)))
	p := (**int)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt8ToPInt(dst, src unsafe.Pointer) {
	var v int
	if p := *(**int8)(unsafe.Pointer(src)); p != nil {
		v = int(*p)
	}

	p := (**int)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int16 to int

func copyInt16ToInt(dst, src unsafe.Pointer) {
	*(*int)(unsafe.Pointer(dst)) = int(*(*int16)(unsafe.Pointer(src)))
}

func copyPInt16ToInt(dst, src unsafe.Pointer) {
	var v int
	if p := *(**int16)(unsafe.Pointer(src)); p != nil {
		v = int(*p)
	}
	*(*int)(unsafe.Pointer(dst)) = v
}

func copyInt16ToPInt(dst, src unsafe.Pointer) {
	v := int(*(*int16)(unsafe.Pointer(src)))
	p := (**int)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt16ToPInt(dst, src unsafe.Pointer) {
	var v int
	if p := *(**int16)(unsafe.Pointer(src)); p != nil {
		v = int(*p)
	}

	p := (**int)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int32 to int

func copyInt32ToInt(dst, src unsafe.Pointer) {
	*(*int)(unsafe.Pointer(dst)) = int(*(*int32)(unsafe.Pointer(src)))
}

func copyPInt32ToInt(dst, src unsafe.Pointer) {
	var v int
	if p := *(**int32)(unsafe.Pointer(src)); p != nil {
		v = int(*p)
	}
	*(*int)(unsafe.Pointer(dst)) = v
}

func copyInt32ToPInt(dst, src unsafe.Pointer) {
	v := int(*(*int32)(unsafe.Pointer(src)))
	p := (**int)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt32ToPInt(dst, src unsafe.Pointer) {
	var v int
	if p := *(**int32)(unsafe.Pointer(src)); p != nil {
		v = int(*p)
	}

	p := (**int)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int64 to int

func copyInt64ToInt(dst, src unsafe.Pointer) {
	*(*int)(unsafe.Pointer(dst)) = int(*(*int64)(unsafe.Pointer(src)))
}

func copyPInt64ToInt(dst, src unsafe.Pointer) {
	var v int
	if p := *(**int64)(unsafe.Pointer(src)); p != nil {
		v = int(*p)
	}
	*(*int)(unsafe.Pointer(dst)) = v
}

func copyInt64ToPInt(dst, src unsafe.Pointer) {
	v := int(*(*int64)(unsafe.Pointer(src)))
	p := (**int)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt64ToPInt(dst, src unsafe.Pointer) {
	var v int
	if p := *(**int64)(unsafe.Pointer(src)); p != nil {
		v = int(*p)
	}

	p := (**int)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint to int

func copyUintToInt(dst, src unsafe.Pointer) {
	*(*int)(unsafe.Pointer(dst)) = int(*(*uint)(unsafe.Pointer(src)))
}

func copyPUintToInt(dst, src unsafe.Pointer) {
	var v int
	if p := *(**uint)(unsafe.Pointer(src)); p != nil {
		v = int(*p)
	}
	*(*int)(unsafe.Pointer(dst)) = v
}

func copyUintToPInt(dst, src unsafe.Pointer) {
	v := int(*(*uint)(unsafe.Pointer(src)))
	p := (**int)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUintToPInt(dst, src unsafe.Pointer) {
	var v int
	if p := *(**uint)(unsafe.Pointer(src)); p != nil {
		v = int(*p)
	}

	p := (**int)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint8 to int

func copyUint8ToInt(dst, src unsafe.Pointer) {
	*(*int)(unsafe.Pointer(dst)) = int(*(*uint8)(unsafe.Pointer(src)))
}

func copyPUint8ToInt(dst, src unsafe.Pointer) {
	var v int
	if p := *(**uint8)(unsafe.Pointer(src)); p != nil {
		v = int(*p)
	}
	*(*int)(unsafe.Pointer(dst)) = v
}

func copyUint8ToPInt(dst, src unsafe.Pointer) {
	v := int(*(*uint8)(unsafe.Pointer(src)))
	p := (**int)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint8ToPInt(dst, src unsafe.Pointer) {
	var v int
	if p := *(**uint8)(unsafe.Pointer(src)); p != nil {
		v = int(*p)
	}

	p := (**int)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint16 to int

func copyUint16ToInt(dst, src unsafe.Pointer) {
	*(*int)(unsafe.Pointer(dst)) = int(*(*uint16)(unsafe.Pointer(src)))
}

func copyPUint16ToInt(dst, src unsafe.Pointer) {
	var v int
	if p := *(**uint16)(unsafe.Pointer(src)); p != nil {
		v = int(*p)
	}
	*(*int)(unsafe.Pointer(dst)) = v
}

func copyUint16ToPInt(dst, src unsafe.Pointer) {
	v := int(*(*uint16)(unsafe.Pointer(src)))
	p := (**int)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint16ToPInt(dst, src unsafe.Pointer) {
	var v int
	if p := *(**uint16)(unsafe.Pointer(src)); p != nil {
		v = int(*p)
	}

	p := (**int)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint32 to int

func copyUint32ToInt(dst, src unsafe.Pointer) {
	*(*int)(unsafe.Pointer(dst)) = int(*(*uint32)(unsafe.Pointer(src)))
}

func copyPUint32ToInt(dst, src unsafe.Pointer) {
	var v int
	if p := *(**uint32)(unsafe.Pointer(src)); p != nil {
		v = int(*p)
	}
	*(*int)(unsafe.Pointer(dst)) = v
}

func copyUint32ToPInt(dst, src unsafe.Pointer) {
	v := int(*(*uint32)(unsafe.Pointer(src)))
	p := (**int)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint32ToPInt(dst, src unsafe.Pointer) {
	var v int
	if p := *(**uint32)(unsafe.Pointer(src)); p != nil {
		v = int(*p)
	}

	p := (**int)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint64 to int

func copyUint64ToInt(dst, src unsafe.Pointer) {
	*(*int)(unsafe.Pointer(dst)) = int(*(*uint64)(unsafe.Pointer(src)))
}

func copyPUint64ToInt(dst, src unsafe.Pointer) {
	var v int
	if p := *(**uint64)(unsafe.Pointer(src)); p != nil {
		v = int(*p)
	}
	*(*int)(unsafe.Pointer(dst)) = v
}

func copyUint64ToPInt(dst, src unsafe.Pointer) {
	v := int(*(*uint64)(unsafe.Pointer(src)))
	p := (**int)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint64ToPInt(dst, src unsafe.Pointer) {
	var v int
	if p := *(**uint64)(unsafe.Pointer(src)); p != nil {
		v = int(*p)
	}

	p := (**int)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int to int8

func copyIntToInt8(dst, src unsafe.Pointer) {
	*(*int8)(unsafe.Pointer(dst)) = int8(*(*int)(unsafe.Pointer(src)))
}

func copyPIntToInt8(dst, src unsafe.Pointer) {
	var v int8
	if p := *(**int)(unsafe.Pointer(src)); p != nil {
		v = int8(*p)
	}
	*(*int8)(unsafe.Pointer(dst)) = v
}

func copyIntToPInt8(dst, src unsafe.Pointer) {
	v := int8(*(*int)(unsafe.Pointer(src)))
	p := (**int8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPIntToPInt8(dst, src unsafe.Pointer) {
	var v int8
	if p := *(**int)(unsafe.Pointer(src)); p != nil {
		v = int8(*p)
	}

	p := (**int8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int8 to int8

func copyInt8ToInt8(dst, src unsafe.Pointer) {
	*(*int8)(unsafe.Pointer(dst)) = int8(*(*int8)(unsafe.Pointer(src)))
}

func copyPInt8ToInt8(dst, src unsafe.Pointer) {
	var v int8
	if p := *(**int8)(unsafe.Pointer(src)); p != nil {
		v = int8(*p)
	}
	*(*int8)(unsafe.Pointer(dst)) = v
}

func copyInt8ToPInt8(dst, src unsafe.Pointer) {
	v := int8(*(*int8)(unsafe.Pointer(src)))
	p := (**int8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt8ToPInt8(dst, src unsafe.Pointer) {
	var v int8
	if p := *(**int8)(unsafe.Pointer(src)); p != nil {
		v = int8(*p)
	}

	p := (**int8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int16 to int8

func copyInt16ToInt8(dst, src unsafe.Pointer) {
	*(*int8)(unsafe.Pointer(dst)) = int8(*(*int16)(unsafe.Pointer(src)))
}

func copyPInt16ToInt8(dst, src unsafe.Pointer) {
	var v int8
	if p := *(**int16)(unsafe.Pointer(src)); p != nil {
		v = int8(*p)
	}
	*(*int8)(unsafe.Pointer(dst)) = v
}

func copyInt16ToPInt8(dst, src unsafe.Pointer) {
	v := int8(*(*int16)(unsafe.Pointer(src)))
	p := (**int8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt16ToPInt8(dst, src unsafe.Pointer) {
	var v int8
	if p := *(**int16)(unsafe.Pointer(src)); p != nil {
		v = int8(*p)
	}

	p := (**int8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int32 to int8

func copyInt32ToInt8(dst, src unsafe.Pointer) {
	*(*int8)(unsafe.Pointer(dst)) = int8(*(*int32)(unsafe.Pointer(src)))
}

func copyPInt32ToInt8(dst, src unsafe.Pointer) {
	var v int8
	if p := *(**int32)(unsafe.Pointer(src)); p != nil {
		v = int8(*p)
	}
	*(*int8)(unsafe.Pointer(dst)) = v
}

func copyInt32ToPInt8(dst, src unsafe.Pointer) {
	v := int8(*(*int32)(unsafe.Pointer(src)))
	p := (**int8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt32ToPInt8(dst, src unsafe.Pointer) {
	var v int8
	if p := *(**int32)(unsafe.Pointer(src)); p != nil {
		v = int8(*p)
	}

	p := (**int8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int64 to int8

func copyInt64ToInt8(dst, src unsafe.Pointer) {
	*(*int8)(unsafe.Pointer(dst)) = int8(*(*int64)(unsafe.Pointer(src)))
}

func copyPInt64ToInt8(dst, src unsafe.Pointer) {
	var v int8
	if p := *(**int64)(unsafe.Pointer(src)); p != nil {
		v = int8(*p)
	}
	*(*int8)(unsafe.Pointer(dst)) = v
}

func copyInt64ToPInt8(dst, src unsafe.Pointer) {
	v := int8(*(*int64)(unsafe.Pointer(src)))
	p := (**int8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt64ToPInt8(dst, src unsafe.Pointer) {
	var v int8
	if p := *(**int64)(unsafe.Pointer(src)); p != nil {
		v = int8(*p)
	}

	p := (**int8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint to int8

func copyUintToInt8(dst, src unsafe.Pointer) {
	*(*int8)(unsafe.Pointer(dst)) = int8(*(*uint)(unsafe.Pointer(src)))
}

func copyPUintToInt8(dst, src unsafe.Pointer) {
	var v int8
	if p := *(**uint)(unsafe.Pointer(src)); p != nil {
		v = int8(*p)
	}
	*(*int8)(unsafe.Pointer(dst)) = v
}

func copyUintToPInt8(dst, src unsafe.Pointer) {
	v := int8(*(*uint)(unsafe.Pointer(src)))
	p := (**int8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUintToPInt8(dst, src unsafe.Pointer) {
	var v int8
	if p := *(**uint)(unsafe.Pointer(src)); p != nil {
		v = int8(*p)
	}

	p := (**int8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint8 to int8

func copyUint8ToInt8(dst, src unsafe.Pointer) {
	*(*int8)(unsafe.Pointer(dst)) = int8(*(*uint8)(unsafe.Pointer(src)))
}

func copyPUint8ToInt8(dst, src unsafe.Pointer) {
	var v int8
	if p := *(**uint8)(unsafe.Pointer(src)); p != nil {
		v = int8(*p)
	}
	*(*int8)(unsafe.Pointer(dst)) = v
}

func copyUint8ToPInt8(dst, src unsafe.Pointer) {
	v := int8(*(*uint8)(unsafe.Pointer(src)))
	p := (**int8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint8ToPInt8(dst, src unsafe.Pointer) {
	var v int8
	if p := *(**uint8)(unsafe.Pointer(src)); p != nil {
		v = int8(*p)
	}

	p := (**int8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint16 to int8

func copyUint16ToInt8(dst, src unsafe.Pointer) {
	*(*int8)(unsafe.Pointer(dst)) = int8(*(*uint16)(unsafe.Pointer(src)))
}

func copyPUint16ToInt8(dst, src unsafe.Pointer) {
	var v int8
	if p := *(**uint16)(unsafe.Pointer(src)); p != nil {
		v = int8(*p)
	}
	*(*int8)(unsafe.Pointer(dst)) = v
}

func copyUint16ToPInt8(dst, src unsafe.Pointer) {
	v := int8(*(*uint16)(unsafe.Pointer(src)))
	p := (**int8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint16ToPInt8(dst, src unsafe.Pointer) {
	var v int8
	if p := *(**uint16)(unsafe.Pointer(src)); p != nil {
		v = int8(*p)
	}

	p := (**int8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint32 to int8

func copyUint32ToInt8(dst, src unsafe.Pointer) {
	*(*int8)(unsafe.Pointer(dst)) = int8(*(*uint32)(unsafe.Pointer(src)))
}

func copyPUint32ToInt8(dst, src unsafe.Pointer) {
	var v int8
	if p := *(**uint32)(unsafe.Pointer(src)); p != nil {
		v = int8(*p)
	}
	*(*int8)(unsafe.Pointer(dst)) = v
}

func copyUint32ToPInt8(dst, src unsafe.Pointer) {
	v := int8(*(*uint32)(unsafe.Pointer(src)))
	p := (**int8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint32ToPInt8(dst, src unsafe.Pointer) {
	var v int8
	if p := *(**uint32)(unsafe.Pointer(src)); p != nil {
		v = int8(*p)
	}

	p := (**int8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint64 to int8

func copyUint64ToInt8(dst, src unsafe.Pointer) {
	*(*int8)(unsafe.Pointer(dst)) = int8(*(*uint64)(unsafe.Pointer(src)))
}

func copyPUint64ToInt8(dst, src unsafe.Pointer) {
	var v int8
	if p := *(**uint64)(unsafe.Pointer(src)); p != nil {
		v = int8(*p)
	}
	*(*int8)(unsafe.Pointer(dst)) = v
}

func copyUint64ToPInt8(dst, src unsafe.Pointer) {
	v := int8(*(*uint64)(unsafe.Pointer(src)))
	p := (**int8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint64ToPInt8(dst, src unsafe.Pointer) {
	var v int8
	if p := *(**uint64)(unsafe.Pointer(src)); p != nil {
		v = int8(*p)
	}

	p := (**int8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int to int16

func copyIntToInt16(dst, src unsafe.Pointer) {
	*(*int16)(unsafe.Pointer(dst)) = int16(*(*int)(unsafe.Pointer(src)))
}

func copyPIntToInt16(dst, src unsafe.Pointer) {
	var v int16
	if p := *(**int)(unsafe.Pointer(src)); p != nil {
		v = int16(*p)
	}
	*(*int16)(unsafe.Pointer(dst)) = v
}

func copyIntToPInt16(dst, src unsafe.Pointer) {
	v := int16(*(*int)(unsafe.Pointer(src)))
	p := (**int16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPIntToPInt16(dst, src unsafe.Pointer) {
	var v int16
	if p := *(**int)(unsafe.Pointer(src)); p != nil {
		v = int16(*p)
	}

	p := (**int16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int8 to int16

func copyInt8ToInt16(dst, src unsafe.Pointer) {
	*(*int16)(unsafe.Pointer(dst)) = int16(*(*int8)(unsafe.Pointer(src)))
}

func copyPInt8ToInt16(dst, src unsafe.Pointer) {
	var v int16
	if p := *(**int8)(unsafe.Pointer(src)); p != nil {
		v = int16(*p)
	}
	*(*int16)(unsafe.Pointer(dst)) = v
}

func copyInt8ToPInt16(dst, src unsafe.Pointer) {
	v := int16(*(*int8)(unsafe.Pointer(src)))
	p := (**int16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt8ToPInt16(dst, src unsafe.Pointer) {
	var v int16
	if p := *(**int8)(unsafe.Pointer(src)); p != nil {
		v = int16(*p)
	}

	p := (**int16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int16 to int16

func copyInt16ToInt16(dst, src unsafe.Pointer) {
	*(*int16)(unsafe.Pointer(dst)) = int16(*(*int16)(unsafe.Pointer(src)))
}

func copyPInt16ToInt16(dst, src unsafe.Pointer) {
	var v int16
	if p := *(**int16)(unsafe.Pointer(src)); p != nil {
		v = int16(*p)
	}
	*(*int16)(unsafe.Pointer(dst)) = v
}

func copyInt16ToPInt16(dst, src unsafe.Pointer) {
	v := int16(*(*int16)(unsafe.Pointer(src)))
	p := (**int16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt16ToPInt16(dst, src unsafe.Pointer) {
	var v int16
	if p := *(**int16)(unsafe.Pointer(src)); p != nil {
		v = int16(*p)
	}

	p := (**int16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int32 to int16

func copyInt32ToInt16(dst, src unsafe.Pointer) {
	*(*int16)(unsafe.Pointer(dst)) = int16(*(*int32)(unsafe.Pointer(src)))
}

func copyPInt32ToInt16(dst, src unsafe.Pointer) {
	var v int16
	if p := *(**int32)(unsafe.Pointer(src)); p != nil {
		v = int16(*p)
	}
	*(*int16)(unsafe.Pointer(dst)) = v
}

func copyInt32ToPInt16(dst, src unsafe.Pointer) {
	v := int16(*(*int32)(unsafe.Pointer(src)))
	p := (**int16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt32ToPInt16(dst, src unsafe.Pointer) {
	var v int16
	if p := *(**int32)(unsafe.Pointer(src)); p != nil {
		v = int16(*p)
	}

	p := (**int16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int64 to int16

func copyInt64ToInt16(dst, src unsafe.Pointer) {
	*(*int16)(unsafe.Pointer(dst)) = int16(*(*int64)(unsafe.Pointer(src)))
}

func copyPInt64ToInt16(dst, src unsafe.Pointer) {
	var v int16
	if p := *(**int64)(unsafe.Pointer(src)); p != nil {
		v = int16(*p)
	}
	*(*int16)(unsafe.Pointer(dst)) = v
}

func copyInt64ToPInt16(dst, src unsafe.Pointer) {
	v := int16(*(*int64)(unsafe.Pointer(src)))
	p := (**int16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt64ToPInt16(dst, src unsafe.Pointer) {
	var v int16
	if p := *(**int64)(unsafe.Pointer(src)); p != nil {
		v = int16(*p)
	}

	p := (**int16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint to int16

func copyUintToInt16(dst, src unsafe.Pointer) {
	*(*int16)(unsafe.Pointer(dst)) = int16(*(*uint)(unsafe.Pointer(src)))
}

func copyPUintToInt16(dst, src unsafe.Pointer) {
	var v int16
	if p := *(**uint)(unsafe.Pointer(src)); p != nil {
		v = int16(*p)
	}
	*(*int16)(unsafe.Pointer(dst)) = v
}

func copyUintToPInt16(dst, src unsafe.Pointer) {
	v := int16(*(*uint)(unsafe.Pointer(src)))
	p := (**int16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUintToPInt16(dst, src unsafe.Pointer) {
	var v int16
	if p := *(**uint)(unsafe.Pointer(src)); p != nil {
		v = int16(*p)
	}

	p := (**int16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint8 to int16

func copyUint8ToInt16(dst, src unsafe.Pointer) {
	*(*int16)(unsafe.Pointer(dst)) = int16(*(*uint8)(unsafe.Pointer(src)))
}

func copyPUint8ToInt16(dst, src unsafe.Pointer) {
	var v int16
	if p := *(**uint8)(unsafe.Pointer(src)); p != nil {
		v = int16(*p)
	}
	*(*int16)(unsafe.Pointer(dst)) = v
}

func copyUint8ToPInt16(dst, src unsafe.Pointer) {
	v := int16(*(*uint8)(unsafe.Pointer(src)))
	p := (**int16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint8ToPInt16(dst, src unsafe.Pointer) {
	var v int16
	if p := *(**uint8)(unsafe.Pointer(src)); p != nil {
		v = int16(*p)
	}

	p := (**int16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint16 to int16

func copyUint16ToInt16(dst, src unsafe.Pointer) {
	*(*int16)(unsafe.Pointer(dst)) = int16(*(*uint16)(unsafe.Pointer(src)))
}

func copyPUint16ToInt16(dst, src unsafe.Pointer) {
	var v int16
	if p := *(**uint16)(unsafe.Pointer(src)); p != nil {
		v = int16(*p)
	}
	*(*int16)(unsafe.Pointer(dst)) = v
}

func copyUint16ToPInt16(dst, src unsafe.Pointer) {
	v := int16(*(*uint16)(unsafe.Pointer(src)))
	p := (**int16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint16ToPInt16(dst, src unsafe.Pointer) {
	var v int16
	if p := *(**uint16)(unsafe.Pointer(src)); p != nil {
		v = int16(*p)
	}

	p := (**int16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint32 to int16

func copyUint32ToInt16(dst, src unsafe.Pointer) {
	*(*int16)(unsafe.Pointer(dst)) = int16(*(*uint32)(unsafe.Pointer(src)))
}

func copyPUint32ToInt16(dst, src unsafe.Pointer) {
	var v int16
	if p := *(**uint32)(unsafe.Pointer(src)); p != nil {
		v = int16(*p)
	}
	*(*int16)(unsafe.Pointer(dst)) = v
}

func copyUint32ToPInt16(dst, src unsafe.Pointer) {
	v := int16(*(*uint32)(unsafe.Pointer(src)))
	p := (**int16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint32ToPInt16(dst, src unsafe.Pointer) {
	var v int16
	if p := *(**uint32)(unsafe.Pointer(src)); p != nil {
		v = int16(*p)
	}

	p := (**int16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint64 to int16

func copyUint64ToInt16(dst, src unsafe.Pointer) {
	*(*int16)(unsafe.Pointer(dst)) = int16(*(*uint64)(unsafe.Pointer(src)))
}

func copyPUint64ToInt16(dst, src unsafe.Pointer) {
	var v int16
	if p := *(**uint64)(unsafe.Pointer(src)); p != nil {
		v = int16(*p)
	}
	*(*int16)(unsafe.Pointer(dst)) = v
}

func copyUint64ToPInt16(dst, src unsafe.Pointer) {
	v := int16(*(*uint64)(unsafe.Pointer(src)))
	p := (**int16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint64ToPInt16(dst, src unsafe.Pointer) {
	var v int16
	if p := *(**uint64)(unsafe.Pointer(src)); p != nil {
		v = int16(*p)
	}

	p := (**int16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int to int32

func copyIntToInt32(dst, src unsafe.Pointer) {
	*(*int32)(unsafe.Pointer(dst)) = int32(*(*int)(unsafe.Pointer(src)))
}

func copyPIntToInt32(dst, src unsafe.Pointer) {
	var v int32
	if p := *(**int)(unsafe.Pointer(src)); p != nil {
		v = int32(*p)
	}
	*(*int32)(unsafe.Pointer(dst)) = v
}

func copyIntToPInt32(dst, src unsafe.Pointer) {
	v := int32(*(*int)(unsafe.Pointer(src)))
	p := (**int32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPIntToPInt32(dst, src unsafe.Pointer) {
	var v int32
	if p := *(**int)(unsafe.Pointer(src)); p != nil {
		v = int32(*p)
	}

	p := (**int32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int8 to int32

func copyInt8ToInt32(dst, src unsafe.Pointer) {
	*(*int32)(unsafe.Pointer(dst)) = int32(*(*int8)(unsafe.Pointer(src)))
}

func copyPInt8ToInt32(dst, src unsafe.Pointer) {
	var v int32
	if p := *(**int8)(unsafe.Pointer(src)); p != nil {
		v = int32(*p)
	}
	*(*int32)(unsafe.Pointer(dst)) = v
}

func copyInt8ToPInt32(dst, src unsafe.Pointer) {
	v := int32(*(*int8)(unsafe.Pointer(src)))
	p := (**int32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt8ToPInt32(dst, src unsafe.Pointer) {
	var v int32
	if p := *(**int8)(unsafe.Pointer(src)); p != nil {
		v = int32(*p)
	}

	p := (**int32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int16 to int32

func copyInt16ToInt32(dst, src unsafe.Pointer) {
	*(*int32)(unsafe.Pointer(dst)) = int32(*(*int16)(unsafe.Pointer(src)))
}

func copyPInt16ToInt32(dst, src unsafe.Pointer) {
	var v int32
	if p := *(**int16)(unsafe.Pointer(src)); p != nil {
		v = int32(*p)
	}
	*(*int32)(unsafe.Pointer(dst)) = v
}

func copyInt16ToPInt32(dst, src unsafe.Pointer) {
	v := int32(*(*int16)(unsafe.Pointer(src)))
	p := (**int32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt16ToPInt32(dst, src unsafe.Pointer) {
	var v int32
	if p := *(**int16)(unsafe.Pointer(src)); p != nil {
		v = int32(*p)
	}

	p := (**int32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int32 to int32

func copyInt32ToInt32(dst, src unsafe.Pointer) {
	*(*int32)(unsafe.Pointer(dst)) = int32(*(*int32)(unsafe.Pointer(src)))
}

func copyPInt32ToInt32(dst, src unsafe.Pointer) {
	var v int32
	if p := *(**int32)(unsafe.Pointer(src)); p != nil {
		v = int32(*p)
	}
	*(*int32)(unsafe.Pointer(dst)) = v
}

func copyInt32ToPInt32(dst, src unsafe.Pointer) {
	v := int32(*(*int32)(unsafe.Pointer(src)))
	p := (**int32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt32ToPInt32(dst, src unsafe.Pointer) {
	var v int32
	if p := *(**int32)(unsafe.Pointer(src)); p != nil {
		v = int32(*p)
	}

	p := (**int32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int64 to int32

func copyInt64ToInt32(dst, src unsafe.Pointer) {
	*(*int32)(unsafe.Pointer(dst)) = int32(*(*int64)(unsafe.Pointer(src)))
}

func copyPInt64ToInt32(dst, src unsafe.Pointer) {
	var v int32
	if p := *(**int64)(unsafe.Pointer(src)); p != nil {
		v = int32(*p)
	}
	*(*int32)(unsafe.Pointer(dst)) = v
}

func copyInt64ToPInt32(dst, src unsafe.Pointer) {
	v := int32(*(*int64)(unsafe.Pointer(src)))
	p := (**int32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt64ToPInt32(dst, src unsafe.Pointer) {
	var v int32
	if p := *(**int64)(unsafe.Pointer(src)); p != nil {
		v = int32(*p)
	}

	p := (**int32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint to int32

func copyUintToInt32(dst, src unsafe.Pointer) {
	*(*int32)(unsafe.Pointer(dst)) = int32(*(*uint)(unsafe.Pointer(src)))
}

func copyPUintToInt32(dst, src unsafe.Pointer) {
	var v int32
	if p := *(**uint)(unsafe.Pointer(src)); p != nil {
		v = int32(*p)
	}
	*(*int32)(unsafe.Pointer(dst)) = v
}

func copyUintToPInt32(dst, src unsafe.Pointer) {
	v := int32(*(*uint)(unsafe.Pointer(src)))
	p := (**int32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUintToPInt32(dst, src unsafe.Pointer) {
	var v int32
	if p := *(**uint)(unsafe.Pointer(src)); p != nil {
		v = int32(*p)
	}

	p := (**int32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint8 to int32

func copyUint8ToInt32(dst, src unsafe.Pointer) {
	*(*int32)(unsafe.Pointer(dst)) = int32(*(*uint8)(unsafe.Pointer(src)))
}

func copyPUint8ToInt32(dst, src unsafe.Pointer) {
	var v int32
	if p := *(**uint8)(unsafe.Pointer(src)); p != nil {
		v = int32(*p)
	}
	*(*int32)(unsafe.Pointer(dst)) = v
}

func copyUint8ToPInt32(dst, src unsafe.Pointer) {
	v := int32(*(*uint8)(unsafe.Pointer(src)))
	p := (**int32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint8ToPInt32(dst, src unsafe.Pointer) {
	var v int32
	if p := *(**uint8)(unsafe.Pointer(src)); p != nil {
		v = int32(*p)
	}

	p := (**int32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint16 to int32

func copyUint16ToInt32(dst, src unsafe.Pointer) {
	*(*int32)(unsafe.Pointer(dst)) = int32(*(*uint16)(unsafe.Pointer(src)))
}

func copyPUint16ToInt32(dst, src unsafe.Pointer) {
	var v int32
	if p := *(**uint16)(unsafe.Pointer(src)); p != nil {
		v = int32(*p)
	}
	*(*int32)(unsafe.Pointer(dst)) = v
}

func copyUint16ToPInt32(dst, src unsafe.Pointer) {
	v := int32(*(*uint16)(unsafe.Pointer(src)))
	p := (**int32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint16ToPInt32(dst, src unsafe.Pointer) {
	var v int32
	if p := *(**uint16)(unsafe.Pointer(src)); p != nil {
		v = int32(*p)
	}

	p := (**int32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint32 to int32

func copyUint32ToInt32(dst, src unsafe.Pointer) {
	*(*int32)(unsafe.Pointer(dst)) = int32(*(*uint32)(unsafe.Pointer(src)))
}

func copyPUint32ToInt32(dst, src unsafe.Pointer) {
	var v int32
	if p := *(**uint32)(unsafe.Pointer(src)); p != nil {
		v = int32(*p)
	}
	*(*int32)(unsafe.Pointer(dst)) = v
}

func copyUint32ToPInt32(dst, src unsafe.Pointer) {
	v := int32(*(*uint32)(unsafe.Pointer(src)))
	p := (**int32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint32ToPInt32(dst, src unsafe.Pointer) {
	var v int32
	if p := *(**uint32)(unsafe.Pointer(src)); p != nil {
		v = int32(*p)
	}

	p := (**int32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint64 to int32

func copyUint64ToInt32(dst, src unsafe.Pointer) {
	*(*int32)(unsafe.Pointer(dst)) = int32(*(*uint64)(unsafe.Pointer(src)))
}

func copyPUint64ToInt32(dst, src unsafe.Pointer) {
	var v int32
	if p := *(**uint64)(unsafe.Pointer(src)); p != nil {
		v = int32(*p)
	}
	*(*int32)(unsafe.Pointer(dst)) = v
}

func copyUint64ToPInt32(dst, src unsafe.Pointer) {
	v := int32(*(*uint64)(unsafe.Pointer(src)))
	p := (**int32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint64ToPInt32(dst, src unsafe.Pointer) {
	var v int32
	if p := *(**uint64)(unsafe.Pointer(src)); p != nil {
		v = int32(*p)
	}

	p := (**int32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int to int64

func copyIntToInt64(dst, src unsafe.Pointer) {
	*(*int64)(unsafe.Pointer(dst)) = int64(*(*int)(unsafe.Pointer(src)))
}

func copyPIntToInt64(dst, src unsafe.Pointer) {
	var v int64
	if p := *(**int)(unsafe.Pointer(src)); p != nil {
		v = int64(*p)
	}
	*(*int64)(unsafe.Pointer(dst)) = v
}

func copyIntToPInt64(dst, src unsafe.Pointer) {
	v := int64(*(*int)(unsafe.Pointer(src)))
	p := (**int64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPIntToPInt64(dst, src unsafe.Pointer) {
	var v int64
	if p := *(**int)(unsafe.Pointer(src)); p != nil {
		v = int64(*p)
	}

	p := (**int64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int8 to int64

func copyInt8ToInt64(dst, src unsafe.Pointer) {
	*(*int64)(unsafe.Pointer(dst)) = int64(*(*int8)(unsafe.Pointer(src)))
}

func copyPInt8ToInt64(dst, src unsafe.Pointer) {
	var v int64
	if p := *(**int8)(unsafe.Pointer(src)); p != nil {
		v = int64(*p)
	}
	*(*int64)(unsafe.Pointer(dst)) = v
}

func copyInt8ToPInt64(dst, src unsafe.Pointer) {
	v := int64(*(*int8)(unsafe.Pointer(src)))
	p := (**int64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt8ToPInt64(dst, src unsafe.Pointer) {
	var v int64
	if p := *(**int8)(unsafe.Pointer(src)); p != nil {
		v = int64(*p)
	}

	p := (**int64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int16 to int64

func copyInt16ToInt64(dst, src unsafe.Pointer) {
	*(*int64)(unsafe.Pointer(dst)) = int64(*(*int16)(unsafe.Pointer(src)))
}

func copyPInt16ToInt64(dst, src unsafe.Pointer) {
	var v int64
	if p := *(**int16)(unsafe.Pointer(src)); p != nil {
		v = int64(*p)
	}
	*(*int64)(unsafe.Pointer(dst)) = v
}

func copyInt16ToPInt64(dst, src unsafe.Pointer) {
	v := int64(*(*int16)(unsafe.Pointer(src)))
	p := (**int64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt16ToPInt64(dst, src unsafe.Pointer) {
	var v int64
	if p := *(**int16)(unsafe.Pointer(src)); p != nil {
		v = int64(*p)
	}

	p := (**int64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int32 to int64

func copyInt32ToInt64(dst, src unsafe.Pointer) {
	*(*int64)(unsafe.Pointer(dst)) = int64(*(*int32)(unsafe.Pointer(src)))
}

func copyPInt32ToInt64(dst, src unsafe.Pointer) {
	var v int64
	if p := *(**int32)(unsafe.Pointer(src)); p != nil {
		v = int64(*p)
	}
	*(*int64)(unsafe.Pointer(dst)) = v
}

func copyInt32ToPInt64(dst, src unsafe.Pointer) {
	v := int64(*(*int32)(unsafe.Pointer(src)))
	p := (**int64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt32ToPInt64(dst, src unsafe.Pointer) {
	var v int64
	if p := *(**int32)(unsafe.Pointer(src)); p != nil {
		v = int64(*p)
	}

	p := (**int64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int64 to int64

func copyInt64ToInt64(dst, src unsafe.Pointer) {
	*(*int64)(unsafe.Pointer(dst)) = int64(*(*int64)(unsafe.Pointer(src)))
}

func copyPInt64ToInt64(dst, src unsafe.Pointer) {
	var v int64
	if p := *(**int64)(unsafe.Pointer(src)); p != nil {
		v = int64(*p)
	}
	*(*int64)(unsafe.Pointer(dst)) = v
}

func copyInt64ToPInt64(dst, src unsafe.Pointer) {
	v := int64(*(*int64)(unsafe.Pointer(src)))
	p := (**int64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt64ToPInt64(dst, src unsafe.Pointer) {
	var v int64
	if p := *(**int64)(unsafe.Pointer(src)); p != nil {
		v = int64(*p)
	}

	p := (**int64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint to int64

func copyUintToInt64(dst, src unsafe.Pointer) {
	*(*int64)(unsafe.Pointer(dst)) = int64(*(*uint)(unsafe.Pointer(src)))
}

func copyPUintToInt64(dst, src unsafe.Pointer) {
	var v int64
	if p := *(**uint)(unsafe.Pointer(src)); p != nil {
		v = int64(*p)
	}
	*(*int64)(unsafe.Pointer(dst)) = v
}

func copyUintToPInt64(dst, src unsafe.Pointer) {
	v := int64(*(*uint)(unsafe.Pointer(src)))
	p := (**int64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUintToPInt64(dst, src unsafe.Pointer) {
	var v int64
	if p := *(**uint)(unsafe.Pointer(src)); p != nil {
		v = int64(*p)
	}

	p := (**int64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint8 to int64

func copyUint8ToInt64(dst, src unsafe.Pointer) {
	*(*int64)(unsafe.Pointer(dst)) = int64(*(*uint8)(unsafe.Pointer(src)))
}

func copyPUint8ToInt64(dst, src unsafe.Pointer) {
	var v int64
	if p := *(**uint8)(unsafe.Pointer(src)); p != nil {
		v = int64(*p)
	}
	*(*int64)(unsafe.Pointer(dst)) = v
}

func copyUint8ToPInt64(dst, src unsafe.Pointer) {
	v := int64(*(*uint8)(unsafe.Pointer(src)))
	p := (**int64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint8ToPInt64(dst, src unsafe.Pointer) {
	var v int64
	if p := *(**uint8)(unsafe.Pointer(src)); p != nil {
		v = int64(*p)
	}

	p := (**int64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint16 to int64

func copyUint16ToInt64(dst, src unsafe.Pointer) {
	*(*int64)(unsafe.Pointer(dst)) = int64(*(*uint16)(unsafe.Pointer(src)))
}

func copyPUint16ToInt64(dst, src unsafe.Pointer) {
	var v int64
	if p := *(**uint16)(unsafe.Pointer(src)); p != nil {
		v = int64(*p)
	}
	*(*int64)(unsafe.Pointer(dst)) = v
}

func copyUint16ToPInt64(dst, src unsafe.Pointer) {
	v := int64(*(*uint16)(unsafe.Pointer(src)))
	p := (**int64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint16ToPInt64(dst, src unsafe.Pointer) {
	var v int64
	if p := *(**uint16)(unsafe.Pointer(src)); p != nil {
		v = int64(*p)
	}

	p := (**int64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint32 to int64

func copyUint32ToInt64(dst, src unsafe.Pointer) {
	*(*int64)(unsafe.Pointer(dst)) = int64(*(*uint32)(unsafe.Pointer(src)))
}

func copyPUint32ToInt64(dst, src unsafe.Pointer) {
	var v int64
	if p := *(**uint32)(unsafe.Pointer(src)); p != nil {
		v = int64(*p)
	}
	*(*int64)(unsafe.Pointer(dst)) = v
}

func copyUint32ToPInt64(dst, src unsafe.Pointer) {
	v := int64(*(*uint32)(unsafe.Pointer(src)))
	p := (**int64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint32ToPInt64(dst, src unsafe.Pointer) {
	var v int64
	if p := *(**uint32)(unsafe.Pointer(src)); p != nil {
		v = int64(*p)
	}

	p := (**int64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint64 to int64

func copyUint64ToInt64(dst, src unsafe.Pointer) {
	*(*int64)(unsafe.Pointer(dst)) = int64(*(*uint64)(unsafe.Pointer(src)))
}

func copyPUint64ToInt64(dst, src unsafe.Pointer) {
	var v int64
	if p := *(**uint64)(unsafe.Pointer(src)); p != nil {
		v = int64(*p)
	}
	*(*int64)(unsafe.Pointer(dst)) = v
}

func copyUint64ToPInt64(dst, src unsafe.Pointer) {
	v := int64(*(*uint64)(unsafe.Pointer(src)))
	p := (**int64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint64ToPInt64(dst, src unsafe.Pointer) {
	var v int64
	if p := *(**uint64)(unsafe.Pointer(src)); p != nil {
		v = int64(*p)
	}

	p := (**int64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int to uint

func copyIntToUint(dst, src unsafe.Pointer) {
	*(*uint)(unsafe.Pointer(dst)) = uint(*(*int)(unsafe.Pointer(src)))
}

func copyPIntToUint(dst, src unsafe.Pointer) {
	var v uint
	if p := *(**int)(unsafe.Pointer(src)); p != nil {
		v = uint(*p)
	}
	*(*uint)(unsafe.Pointer(dst)) = v
}

func copyIntToPUint(dst, src unsafe.Pointer) {
	v := uint(*(*int)(unsafe.Pointer(src)))
	p := (**uint)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPIntToPUint(dst, src unsafe.Pointer) {
	var v uint
	if p := *(**int)(unsafe.Pointer(src)); p != nil {
		v = uint(*p)
	}

	p := (**uint)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int8 to uint

func copyInt8ToUint(dst, src unsafe.Pointer) {
	*(*uint)(unsafe.Pointer(dst)) = uint(*(*int8)(unsafe.Pointer(src)))
}

func copyPInt8ToUint(dst, src unsafe.Pointer) {
	var v uint
	if p := *(**int8)(unsafe.Pointer(src)); p != nil {
		v = uint(*p)
	}
	*(*uint)(unsafe.Pointer(dst)) = v
}

func copyInt8ToPUint(dst, src unsafe.Pointer) {
	v := uint(*(*int8)(unsafe.Pointer(src)))
	p := (**uint)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt8ToPUint(dst, src unsafe.Pointer) {
	var v uint
	if p := *(**int8)(unsafe.Pointer(src)); p != nil {
		v = uint(*p)
	}

	p := (**uint)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int16 to uint

func copyInt16ToUint(dst, src unsafe.Pointer) {
	*(*uint)(unsafe.Pointer(dst)) = uint(*(*int16)(unsafe.Pointer(src)))
}

func copyPInt16ToUint(dst, src unsafe.Pointer) {
	var v uint
	if p := *(**int16)(unsafe.Pointer(src)); p != nil {
		v = uint(*p)
	}
	*(*uint)(unsafe.Pointer(dst)) = v
}

func copyInt16ToPUint(dst, src unsafe.Pointer) {
	v := uint(*(*int16)(unsafe.Pointer(src)))
	p := (**uint)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt16ToPUint(dst, src unsafe.Pointer) {
	var v uint
	if p := *(**int16)(unsafe.Pointer(src)); p != nil {
		v = uint(*p)
	}

	p := (**uint)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int32 to uint

func copyInt32ToUint(dst, src unsafe.Pointer) {
	*(*uint)(unsafe.Pointer(dst)) = uint(*(*int32)(unsafe.Pointer(src)))
}

func copyPInt32ToUint(dst, src unsafe.Pointer) {
	var v uint
	if p := *(**int32)(unsafe.Pointer(src)); p != nil {
		v = uint(*p)
	}
	*(*uint)(unsafe.Pointer(dst)) = v
}

func copyInt32ToPUint(dst, src unsafe.Pointer) {
	v := uint(*(*int32)(unsafe.Pointer(src)))
	p := (**uint)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt32ToPUint(dst, src unsafe.Pointer) {
	var v uint
	if p := *(**int32)(unsafe.Pointer(src)); p != nil {
		v = uint(*p)
	}

	p := (**uint)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int64 to uint

func copyInt64ToUint(dst, src unsafe.Pointer) {
	*(*uint)(unsafe.Pointer(dst)) = uint(*(*int64)(unsafe.Pointer(src)))
}

func copyPInt64ToUint(dst, src unsafe.Pointer) {
	var v uint
	if p := *(**int64)(unsafe.Pointer(src)); p != nil {
		v = uint(*p)
	}
	*(*uint)(unsafe.Pointer(dst)) = v
}

func copyInt64ToPUint(dst, src unsafe.Pointer) {
	v := uint(*(*int64)(unsafe.Pointer(src)))
	p := (**uint)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt64ToPUint(dst, src unsafe.Pointer) {
	var v uint
	if p := *(**int64)(unsafe.Pointer(src)); p != nil {
		v = uint(*p)
	}

	p := (**uint)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint to uint

func copyUintToUint(dst, src unsafe.Pointer) {
	*(*uint)(unsafe.Pointer(dst)) = uint(*(*uint)(unsafe.Pointer(src)))
}

func copyPUintToUint(dst, src unsafe.Pointer) {
	var v uint
	if p := *(**uint)(unsafe.Pointer(src)); p != nil {
		v = uint(*p)
	}
	*(*uint)(unsafe.Pointer(dst)) = v
}

func copyUintToPUint(dst, src unsafe.Pointer) {
	v := uint(*(*uint)(unsafe.Pointer(src)))
	p := (**uint)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUintToPUint(dst, src unsafe.Pointer) {
	var v uint
	if p := *(**uint)(unsafe.Pointer(src)); p != nil {
		v = uint(*p)
	}

	p := (**uint)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint8 to uint

func copyUint8ToUint(dst, src unsafe.Pointer) {
	*(*uint)(unsafe.Pointer(dst)) = uint(*(*uint8)(unsafe.Pointer(src)))
}

func copyPUint8ToUint(dst, src unsafe.Pointer) {
	var v uint
	if p := *(**uint8)(unsafe.Pointer(src)); p != nil {
		v = uint(*p)
	}
	*(*uint)(unsafe.Pointer(dst)) = v
}

func copyUint8ToPUint(dst, src unsafe.Pointer) {
	v := uint(*(*uint8)(unsafe.Pointer(src)))
	p := (**uint)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint8ToPUint(dst, src unsafe.Pointer) {
	var v uint
	if p := *(**uint8)(unsafe.Pointer(src)); p != nil {
		v = uint(*p)
	}

	p := (**uint)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint16 to uint

func copyUint16ToUint(dst, src unsafe.Pointer) {
	*(*uint)(unsafe.Pointer(dst)) = uint(*(*uint16)(unsafe.Pointer(src)))
}

func copyPUint16ToUint(dst, src unsafe.Pointer) {
	var v uint
	if p := *(**uint16)(unsafe.Pointer(src)); p != nil {
		v = uint(*p)
	}
	*(*uint)(unsafe.Pointer(dst)) = v
}

func copyUint16ToPUint(dst, src unsafe.Pointer) {
	v := uint(*(*uint16)(unsafe.Pointer(src)))
	p := (**uint)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint16ToPUint(dst, src unsafe.Pointer) {
	var v uint
	if p := *(**uint16)(unsafe.Pointer(src)); p != nil {
		v = uint(*p)
	}

	p := (**uint)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint32 to uint

func copyUint32ToUint(dst, src unsafe.Pointer) {
	*(*uint)(unsafe.Pointer(dst)) = uint(*(*uint32)(unsafe.Pointer(src)))
}

func copyPUint32ToUint(dst, src unsafe.Pointer) {
	var v uint
	if p := *(**uint32)(unsafe.Pointer(src)); p != nil {
		v = uint(*p)
	}
	*(*uint)(unsafe.Pointer(dst)) = v
}

func copyUint32ToPUint(dst, src unsafe.Pointer) {
	v := uint(*(*uint32)(unsafe.Pointer(src)))
	p := (**uint)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint32ToPUint(dst, src unsafe.Pointer) {
	var v uint
	if p := *(**uint32)(unsafe.Pointer(src)); p != nil {
		v = uint(*p)
	}

	p := (**uint)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint64 to uint

func copyUint64ToUint(dst, src unsafe.Pointer) {
	*(*uint)(unsafe.Pointer(dst)) = uint(*(*uint64)(unsafe.Pointer(src)))
}

func copyPUint64ToUint(dst, src unsafe.Pointer) {
	var v uint
	if p := *(**uint64)(unsafe.Pointer(src)); p != nil {
		v = uint(*p)
	}
	*(*uint)(unsafe.Pointer(dst)) = v
}

func copyUint64ToPUint(dst, src unsafe.Pointer) {
	v := uint(*(*uint64)(unsafe.Pointer(src)))
	p := (**uint)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint64ToPUint(dst, src unsafe.Pointer) {
	var v uint
	if p := *(**uint64)(unsafe.Pointer(src)); p != nil {
		v = uint(*p)
	}

	p := (**uint)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int to uint8

func copyIntToUint8(dst, src unsafe.Pointer) {
	*(*uint8)(unsafe.Pointer(dst)) = uint8(*(*int)(unsafe.Pointer(src)))
}

func copyPIntToUint8(dst, src unsafe.Pointer) {
	var v uint8
	if p := *(**int)(unsafe.Pointer(src)); p != nil {
		v = uint8(*p)
	}
	*(*uint8)(unsafe.Pointer(dst)) = v
}

func copyIntToPUint8(dst, src unsafe.Pointer) {
	v := uint8(*(*int)(unsafe.Pointer(src)))
	p := (**uint8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPIntToPUint8(dst, src unsafe.Pointer) {
	var v uint8
	if p := *(**int)(unsafe.Pointer(src)); p != nil {
		v = uint8(*p)
	}

	p := (**uint8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int8 to uint8

func copyInt8ToUint8(dst, src unsafe.Pointer) {
	*(*uint8)(unsafe.Pointer(dst)) = uint8(*(*int8)(unsafe.Pointer(src)))
}

func copyPInt8ToUint8(dst, src unsafe.Pointer) {
	var v uint8
	if p := *(**int8)(unsafe.Pointer(src)); p != nil {
		v = uint8(*p)
	}
	*(*uint8)(unsafe.Pointer(dst)) = v
}

func copyInt8ToPUint8(dst, src unsafe.Pointer) {
	v := uint8(*(*int8)(unsafe.Pointer(src)))
	p := (**uint8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt8ToPUint8(dst, src unsafe.Pointer) {
	var v uint8
	if p := *(**int8)(unsafe.Pointer(src)); p != nil {
		v = uint8(*p)
	}

	p := (**uint8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int16 to uint8

func copyInt16ToUint8(dst, src unsafe.Pointer) {
	*(*uint8)(unsafe.Pointer(dst)) = uint8(*(*int16)(unsafe.Pointer(src)))
}

func copyPInt16ToUint8(dst, src unsafe.Pointer) {
	var v uint8
	if p := *(**int16)(unsafe.Pointer(src)); p != nil {
		v = uint8(*p)
	}
	*(*uint8)(unsafe.Pointer(dst)) = v
}

func copyInt16ToPUint8(dst, src unsafe.Pointer) {
	v := uint8(*(*int16)(unsafe.Pointer(src)))
	p := (**uint8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt16ToPUint8(dst, src unsafe.Pointer) {
	var v uint8
	if p := *(**int16)(unsafe.Pointer(src)); p != nil {
		v = uint8(*p)
	}

	p := (**uint8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int32 to uint8

func copyInt32ToUint8(dst, src unsafe.Pointer) {
	*(*uint8)(unsafe.Pointer(dst)) = uint8(*(*int32)(unsafe.Pointer(src)))
}

func copyPInt32ToUint8(dst, src unsafe.Pointer) {
	var v uint8
	if p := *(**int32)(unsafe.Pointer(src)); p != nil {
		v = uint8(*p)
	}
	*(*uint8)(unsafe.Pointer(dst)) = v
}

func copyInt32ToPUint8(dst, src unsafe.Pointer) {
	v := uint8(*(*int32)(unsafe.Pointer(src)))
	p := (**uint8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt32ToPUint8(dst, src unsafe.Pointer) {
	var v uint8
	if p := *(**int32)(unsafe.Pointer(src)); p != nil {
		v = uint8(*p)
	}

	p := (**uint8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int64 to uint8

func copyInt64ToUint8(dst, src unsafe.Pointer) {
	*(*uint8)(unsafe.Pointer(dst)) = uint8(*(*int64)(unsafe.Pointer(src)))
}

func copyPInt64ToUint8(dst, src unsafe.Pointer) {
	var v uint8
	if p := *(**int64)(unsafe.Pointer(src)); p != nil {
		v = uint8(*p)
	}
	*(*uint8)(unsafe.Pointer(dst)) = v
}

func copyInt64ToPUint8(dst, src unsafe.Pointer) {
	v := uint8(*(*int64)(unsafe.Pointer(src)))
	p := (**uint8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt64ToPUint8(dst, src unsafe.Pointer) {
	var v uint8
	if p := *(**int64)(unsafe.Pointer(src)); p != nil {
		v = uint8(*p)
	}

	p := (**uint8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint to uint8

func copyUintToUint8(dst, src unsafe.Pointer) {
	*(*uint8)(unsafe.Pointer(dst)) = uint8(*(*uint)(unsafe.Pointer(src)))
}

func copyPUintToUint8(dst, src unsafe.Pointer) {
	var v uint8
	if p := *(**uint)(unsafe.Pointer(src)); p != nil {
		v = uint8(*p)
	}
	*(*uint8)(unsafe.Pointer(dst)) = v
}

func copyUintToPUint8(dst, src unsafe.Pointer) {
	v := uint8(*(*uint)(unsafe.Pointer(src)))
	p := (**uint8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUintToPUint8(dst, src unsafe.Pointer) {
	var v uint8
	if p := *(**uint)(unsafe.Pointer(src)); p != nil {
		v = uint8(*p)
	}

	p := (**uint8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint8 to uint8

func copyUint8ToUint8(dst, src unsafe.Pointer) {
	*(*uint8)(unsafe.Pointer(dst)) = uint8(*(*uint8)(unsafe.Pointer(src)))
}

func copyPUint8ToUint8(dst, src unsafe.Pointer) {
	var v uint8
	if p := *(**uint8)(unsafe.Pointer(src)); p != nil {
		v = uint8(*p)
	}
	*(*uint8)(unsafe.Pointer(dst)) = v
}

func copyUint8ToPUint8(dst, src unsafe.Pointer) {
	v := uint8(*(*uint8)(unsafe.Pointer(src)))
	p := (**uint8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint8ToPUint8(dst, src unsafe.Pointer) {
	var v uint8
	if p := *(**uint8)(unsafe.Pointer(src)); p != nil {
		v = uint8(*p)
	}

	p := (**uint8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint16 to uint8

func copyUint16ToUint8(dst, src unsafe.Pointer) {
	*(*uint8)(unsafe.Pointer(dst)) = uint8(*(*uint16)(unsafe.Pointer(src)))
}

func copyPUint16ToUint8(dst, src unsafe.Pointer) {
	var v uint8
	if p := *(**uint16)(unsafe.Pointer(src)); p != nil {
		v = uint8(*p)
	}
	*(*uint8)(unsafe.Pointer(dst)) = v
}

func copyUint16ToPUint8(dst, src unsafe.Pointer) {
	v := uint8(*(*uint16)(unsafe.Pointer(src)))
	p := (**uint8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint16ToPUint8(dst, src unsafe.Pointer) {
	var v uint8
	if p := *(**uint16)(unsafe.Pointer(src)); p != nil {
		v = uint8(*p)
	}

	p := (**uint8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint32 to uint8

func copyUint32ToUint8(dst, src unsafe.Pointer) {
	*(*uint8)(unsafe.Pointer(dst)) = uint8(*(*uint32)(unsafe.Pointer(src)))
}

func copyPUint32ToUint8(dst, src unsafe.Pointer) {
	var v uint8
	if p := *(**uint32)(unsafe.Pointer(src)); p != nil {
		v = uint8(*p)
	}
	*(*uint8)(unsafe.Pointer(dst)) = v
}

func copyUint32ToPUint8(dst, src unsafe.Pointer) {
	v := uint8(*(*uint32)(unsafe.Pointer(src)))
	p := (**uint8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint32ToPUint8(dst, src unsafe.Pointer) {
	var v uint8
	if p := *(**uint32)(unsafe.Pointer(src)); p != nil {
		v = uint8(*p)
	}

	p := (**uint8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint64 to uint8

func copyUint64ToUint8(dst, src unsafe.Pointer) {
	*(*uint8)(unsafe.Pointer(dst)) = uint8(*(*uint64)(unsafe.Pointer(src)))
}

func copyPUint64ToUint8(dst, src unsafe.Pointer) {
	var v uint8
	if p := *(**uint64)(unsafe.Pointer(src)); p != nil {
		v = uint8(*p)
	}
	*(*uint8)(unsafe.Pointer(dst)) = v
}

func copyUint64ToPUint8(dst, src unsafe.Pointer) {
	v := uint8(*(*uint64)(unsafe.Pointer(src)))
	p := (**uint8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint64ToPUint8(dst, src unsafe.Pointer) {
	var v uint8
	if p := *(**uint64)(unsafe.Pointer(src)); p != nil {
		v = uint8(*p)
	}

	p := (**uint8)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int to uint16

func copyIntToUint16(dst, src unsafe.Pointer) {
	*(*uint16)(unsafe.Pointer(dst)) = uint16(*(*int)(unsafe.Pointer(src)))
}

func copyPIntToUint16(dst, src unsafe.Pointer) {
	var v uint16
	if p := *(**int)(unsafe.Pointer(src)); p != nil {
		v = uint16(*p)
	}
	*(*uint16)(unsafe.Pointer(dst)) = v
}

func copyIntToPUint16(dst, src unsafe.Pointer) {
	v := uint16(*(*int)(unsafe.Pointer(src)))
	p := (**uint16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPIntToPUint16(dst, src unsafe.Pointer) {
	var v uint16
	if p := *(**int)(unsafe.Pointer(src)); p != nil {
		v = uint16(*p)
	}

	p := (**uint16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int8 to uint16

func copyInt8ToUint16(dst, src unsafe.Pointer) {
	*(*uint16)(unsafe.Pointer(dst)) = uint16(*(*int8)(unsafe.Pointer(src)))
}

func copyPInt8ToUint16(dst, src unsafe.Pointer) {
	var v uint16
	if p := *(**int8)(unsafe.Pointer(src)); p != nil {
		v = uint16(*p)
	}
	*(*uint16)(unsafe.Pointer(dst)) = v
}

func copyInt8ToPUint16(dst, src unsafe.Pointer) {
	v := uint16(*(*int8)(unsafe.Pointer(src)))
	p := (**uint16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt8ToPUint16(dst, src unsafe.Pointer) {
	var v uint16
	if p := *(**int8)(unsafe.Pointer(src)); p != nil {
		v = uint16(*p)
	}

	p := (**uint16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int16 to uint16

func copyInt16ToUint16(dst, src unsafe.Pointer) {
	*(*uint16)(unsafe.Pointer(dst)) = uint16(*(*int16)(unsafe.Pointer(src)))
}

func copyPInt16ToUint16(dst, src unsafe.Pointer) {
	var v uint16
	if p := *(**int16)(unsafe.Pointer(src)); p != nil {
		v = uint16(*p)
	}
	*(*uint16)(unsafe.Pointer(dst)) = v
}

func copyInt16ToPUint16(dst, src unsafe.Pointer) {
	v := uint16(*(*int16)(unsafe.Pointer(src)))
	p := (**uint16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt16ToPUint16(dst, src unsafe.Pointer) {
	var v uint16
	if p := *(**int16)(unsafe.Pointer(src)); p != nil {
		v = uint16(*p)
	}

	p := (**uint16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int32 to uint16

func copyInt32ToUint16(dst, src unsafe.Pointer) {
	*(*uint16)(unsafe.Pointer(dst)) = uint16(*(*int32)(unsafe.Pointer(src)))
}

func copyPInt32ToUint16(dst, src unsafe.Pointer) {
	var v uint16
	if p := *(**int32)(unsafe.Pointer(src)); p != nil {
		v = uint16(*p)
	}
	*(*uint16)(unsafe.Pointer(dst)) = v
}

func copyInt32ToPUint16(dst, src unsafe.Pointer) {
	v := uint16(*(*int32)(unsafe.Pointer(src)))
	p := (**uint16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt32ToPUint16(dst, src unsafe.Pointer) {
	var v uint16
	if p := *(**int32)(unsafe.Pointer(src)); p != nil {
		v = uint16(*p)
	}

	p := (**uint16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int64 to uint16

func copyInt64ToUint16(dst, src unsafe.Pointer) {
	*(*uint16)(unsafe.Pointer(dst)) = uint16(*(*int64)(unsafe.Pointer(src)))
}

func copyPInt64ToUint16(dst, src unsafe.Pointer) {
	var v uint16
	if p := *(**int64)(unsafe.Pointer(src)); p != nil {
		v = uint16(*p)
	}
	*(*uint16)(unsafe.Pointer(dst)) = v
}

func copyInt64ToPUint16(dst, src unsafe.Pointer) {
	v := uint16(*(*int64)(unsafe.Pointer(src)))
	p := (**uint16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt64ToPUint16(dst, src unsafe.Pointer) {
	var v uint16
	if p := *(**int64)(unsafe.Pointer(src)); p != nil {
		v = uint16(*p)
	}

	p := (**uint16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint to uint16

func copyUintToUint16(dst, src unsafe.Pointer) {
	*(*uint16)(unsafe.Pointer(dst)) = uint16(*(*uint)(unsafe.Pointer(src)))
}

func copyPUintToUint16(dst, src unsafe.Pointer) {
	var v uint16
	if p := *(**uint)(unsafe.Pointer(src)); p != nil {
		v = uint16(*p)
	}
	*(*uint16)(unsafe.Pointer(dst)) = v
}

func copyUintToPUint16(dst, src unsafe.Pointer) {
	v := uint16(*(*uint)(unsafe.Pointer(src)))
	p := (**uint16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUintToPUint16(dst, src unsafe.Pointer) {
	var v uint16
	if p := *(**uint)(unsafe.Pointer(src)); p != nil {
		v = uint16(*p)
	}

	p := (**uint16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint8 to uint16

func copyUint8ToUint16(dst, src unsafe.Pointer) {
	*(*uint16)(unsafe.Pointer(dst)) = uint16(*(*uint8)(unsafe.Pointer(src)))
}

func copyPUint8ToUint16(dst, src unsafe.Pointer) {
	var v uint16
	if p := *(**uint8)(unsafe.Pointer(src)); p != nil {
		v = uint16(*p)
	}
	*(*uint16)(unsafe.Pointer(dst)) = v
}

func copyUint8ToPUint16(dst, src unsafe.Pointer) {
	v := uint16(*(*uint8)(unsafe.Pointer(src)))
	p := (**uint16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint8ToPUint16(dst, src unsafe.Pointer) {
	var v uint16
	if p := *(**uint8)(unsafe.Pointer(src)); p != nil {
		v = uint16(*p)
	}

	p := (**uint16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint16 to uint16

func copyUint16ToUint16(dst, src unsafe.Pointer) {
	*(*uint16)(unsafe.Pointer(dst)) = uint16(*(*uint16)(unsafe.Pointer(src)))
}

func copyPUint16ToUint16(dst, src unsafe.Pointer) {
	var v uint16
	if p := *(**uint16)(unsafe.Pointer(src)); p != nil {
		v = uint16(*p)
	}
	*(*uint16)(unsafe.Pointer(dst)) = v
}

func copyUint16ToPUint16(dst, src unsafe.Pointer) {
	v := uint16(*(*uint16)(unsafe.Pointer(src)))
	p := (**uint16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint16ToPUint16(dst, src unsafe.Pointer) {
	var v uint16
	if p := *(**uint16)(unsafe.Pointer(src)); p != nil {
		v = uint16(*p)
	}

	p := (**uint16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint32 to uint16

func copyUint32ToUint16(dst, src unsafe.Pointer) {
	*(*uint16)(unsafe.Pointer(dst)) = uint16(*(*uint32)(unsafe.Pointer(src)))
}

func copyPUint32ToUint16(dst, src unsafe.Pointer) {
	var v uint16
	if p := *(**uint32)(unsafe.Pointer(src)); p != nil {
		v = uint16(*p)
	}
	*(*uint16)(unsafe.Pointer(dst)) = v
}

func copyUint32ToPUint16(dst, src unsafe.Pointer) {
	v := uint16(*(*uint32)(unsafe.Pointer(src)))
	p := (**uint16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint32ToPUint16(dst, src unsafe.Pointer) {
	var v uint16
	if p := *(**uint32)(unsafe.Pointer(src)); p != nil {
		v = uint16(*p)
	}

	p := (**uint16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint64 to uint16

func copyUint64ToUint16(dst, src unsafe.Pointer) {
	*(*uint16)(unsafe.Pointer(dst)) = uint16(*(*uint64)(unsafe.Pointer(src)))
}

func copyPUint64ToUint16(dst, src unsafe.Pointer) {
	var v uint16
	if p := *(**uint64)(unsafe.Pointer(src)); p != nil {
		v = uint16(*p)
	}
	*(*uint16)(unsafe.Pointer(dst)) = v
}

func copyUint64ToPUint16(dst, src unsafe.Pointer) {
	v := uint16(*(*uint64)(unsafe.Pointer(src)))
	p := (**uint16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint64ToPUint16(dst, src unsafe.Pointer) {
	var v uint16
	if p := *(**uint64)(unsafe.Pointer(src)); p != nil {
		v = uint16(*p)
	}

	p := (**uint16)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int to uint32

func copyIntToUint32(dst, src unsafe.Pointer) {
	*(*uint32)(unsafe.Pointer(dst)) = uint32(*(*int)(unsafe.Pointer(src)))
}

func copyPIntToUint32(dst, src unsafe.Pointer) {
	var v uint32
	if p := *(**int)(unsafe.Pointer(src)); p != nil {
		v = uint32(*p)
	}
	*(*uint32)(unsafe.Pointer(dst)) = v
}

func copyIntToPUint32(dst, src unsafe.Pointer) {
	v := uint32(*(*int)(unsafe.Pointer(src)))
	p := (**uint32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPIntToPUint32(dst, src unsafe.Pointer) {
	var v uint32
	if p := *(**int)(unsafe.Pointer(src)); p != nil {
		v = uint32(*p)
	}

	p := (**uint32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int8 to uint32

func copyInt8ToUint32(dst, src unsafe.Pointer) {
	*(*uint32)(unsafe.Pointer(dst)) = uint32(*(*int8)(unsafe.Pointer(src)))
}

func copyPInt8ToUint32(dst, src unsafe.Pointer) {
	var v uint32
	if p := *(**int8)(unsafe.Pointer(src)); p != nil {
		v = uint32(*p)
	}
	*(*uint32)(unsafe.Pointer(dst)) = v
}

func copyInt8ToPUint32(dst, src unsafe.Pointer) {
	v := uint32(*(*int8)(unsafe.Pointer(src)))
	p := (**uint32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt8ToPUint32(dst, src unsafe.Pointer) {
	var v uint32
	if p := *(**int8)(unsafe.Pointer(src)); p != nil {
		v = uint32(*p)
	}

	p := (**uint32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int16 to uint32

func copyInt16ToUint32(dst, src unsafe.Pointer) {
	*(*uint32)(unsafe.Pointer(dst)) = uint32(*(*int16)(unsafe.Pointer(src)))
}

func copyPInt16ToUint32(dst, src unsafe.Pointer) {
	var v uint32
	if p := *(**int16)(unsafe.Pointer(src)); p != nil {
		v = uint32(*p)
	}
	*(*uint32)(unsafe.Pointer(dst)) = v
}

func copyInt16ToPUint32(dst, src unsafe.Pointer) {
	v := uint32(*(*int16)(unsafe.Pointer(src)))
	p := (**uint32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt16ToPUint32(dst, src unsafe.Pointer) {
	var v uint32
	if p := *(**int16)(unsafe.Pointer(src)); p != nil {
		v = uint32(*p)
	}

	p := (**uint32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int32 to uint32

func copyInt32ToUint32(dst, src unsafe.Pointer) {
	*(*uint32)(unsafe.Pointer(dst)) = uint32(*(*int32)(unsafe.Pointer(src)))
}

func copyPInt32ToUint32(dst, src unsafe.Pointer) {
	var v uint32
	if p := *(**int32)(unsafe.Pointer(src)); p != nil {
		v = uint32(*p)
	}
	*(*uint32)(unsafe.Pointer(dst)) = v
}

func copyInt32ToPUint32(dst, src unsafe.Pointer) {
	v := uint32(*(*int32)(unsafe.Pointer(src)))
	p := (**uint32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt32ToPUint32(dst, src unsafe.Pointer) {
	var v uint32
	if p := *(**int32)(unsafe.Pointer(src)); p != nil {
		v = uint32(*p)
	}

	p := (**uint32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int64 to uint32

func copyInt64ToUint32(dst, src unsafe.Pointer) {
	*(*uint32)(unsafe.Pointer(dst)) = uint32(*(*int64)(unsafe.Pointer(src)))
}

func copyPInt64ToUint32(dst, src unsafe.Pointer) {
	var v uint32
	if p := *(**int64)(unsafe.Pointer(src)); p != nil {
		v = uint32(*p)
	}
	*(*uint32)(unsafe.Pointer(dst)) = v
}

func copyInt64ToPUint32(dst, src unsafe.Pointer) {
	v := uint32(*(*int64)(unsafe.Pointer(src)))
	p := (**uint32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt64ToPUint32(dst, src unsafe.Pointer) {
	var v uint32
	if p := *(**int64)(unsafe.Pointer(src)); p != nil {
		v = uint32(*p)
	}

	p := (**uint32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint to uint32

func copyUintToUint32(dst, src unsafe.Pointer) {
	*(*uint32)(unsafe.Pointer(dst)) = uint32(*(*uint)(unsafe.Pointer(src)))
}

func copyPUintToUint32(dst, src unsafe.Pointer) {
	var v uint32
	if p := *(**uint)(unsafe.Pointer(src)); p != nil {
		v = uint32(*p)
	}
	*(*uint32)(unsafe.Pointer(dst)) = v
}

func copyUintToPUint32(dst, src unsafe.Pointer) {
	v := uint32(*(*uint)(unsafe.Pointer(src)))
	p := (**uint32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUintToPUint32(dst, src unsafe.Pointer) {
	var v uint32
	if p := *(**uint)(unsafe.Pointer(src)); p != nil {
		v = uint32(*p)
	}

	p := (**uint32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint8 to uint32

func copyUint8ToUint32(dst, src unsafe.Pointer) {
	*(*uint32)(unsafe.Pointer(dst)) = uint32(*(*uint8)(unsafe.Pointer(src)))
}

func copyPUint8ToUint32(dst, src unsafe.Pointer) {
	var v uint32
	if p := *(**uint8)(unsafe.Pointer(src)); p != nil {
		v = uint32(*p)
	}
	*(*uint32)(unsafe.Pointer(dst)) = v
}

func copyUint8ToPUint32(dst, src unsafe.Pointer) {
	v := uint32(*(*uint8)(unsafe.Pointer(src)))
	p := (**uint32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint8ToPUint32(dst, src unsafe.Pointer) {
	var v uint32
	if p := *(**uint8)(unsafe.Pointer(src)); p != nil {
		v = uint32(*p)
	}

	p := (**uint32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint16 to uint32

func copyUint16ToUint32(dst, src unsafe.Pointer) {
	*(*uint32)(unsafe.Pointer(dst)) = uint32(*(*uint16)(unsafe.Pointer(src)))
}

func copyPUint16ToUint32(dst, src unsafe.Pointer) {
	var v uint32
	if p := *(**uint16)(unsafe.Pointer(src)); p != nil {
		v = uint32(*p)
	}
	*(*uint32)(unsafe.Pointer(dst)) = v
}

func copyUint16ToPUint32(dst, src unsafe.Pointer) {
	v := uint32(*(*uint16)(unsafe.Pointer(src)))
	p := (**uint32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint16ToPUint32(dst, src unsafe.Pointer) {
	var v uint32
	if p := *(**uint16)(unsafe.Pointer(src)); p != nil {
		v = uint32(*p)
	}

	p := (**uint32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint32 to uint32

func copyUint32ToUint32(dst, src unsafe.Pointer) {
	*(*uint32)(unsafe.Pointer(dst)) = uint32(*(*uint32)(unsafe.Pointer(src)))
}

func copyPUint32ToUint32(dst, src unsafe.Pointer) {
	var v uint32
	if p := *(**uint32)(unsafe.Pointer(src)); p != nil {
		v = uint32(*p)
	}
	*(*uint32)(unsafe.Pointer(dst)) = v
}

func copyUint32ToPUint32(dst, src unsafe.Pointer) {
	v := uint32(*(*uint32)(unsafe.Pointer(src)))
	p := (**uint32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint32ToPUint32(dst, src unsafe.Pointer) {
	var v uint32
	if p := *(**uint32)(unsafe.Pointer(src)); p != nil {
		v = uint32(*p)
	}

	p := (**uint32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint64 to uint32

func copyUint64ToUint32(dst, src unsafe.Pointer) {
	*(*uint32)(unsafe.Pointer(dst)) = uint32(*(*uint64)(unsafe.Pointer(src)))
}

func copyPUint64ToUint32(dst, src unsafe.Pointer) {
	var v uint32
	if p := *(**uint64)(unsafe.Pointer(src)); p != nil {
		v = uint32(*p)
	}
	*(*uint32)(unsafe.Pointer(dst)) = v
}

func copyUint64ToPUint32(dst, src unsafe.Pointer) {
	v := uint32(*(*uint64)(unsafe.Pointer(src)))
	p := (**uint32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint64ToPUint32(dst, src unsafe.Pointer) {
	var v uint32
	if p := *(**uint64)(unsafe.Pointer(src)); p != nil {
		v = uint32(*p)
	}

	p := (**uint32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int to uint64

func copyIntToUint64(dst, src unsafe.Pointer) {
	*(*uint64)(unsafe.Pointer(dst)) = uint64(*(*int)(unsafe.Pointer(src)))
}

func copyPIntToUint64(dst, src unsafe.Pointer) {
	var v uint64
	if p := *(**int)(unsafe.Pointer(src)); p != nil {
		v = uint64(*p)
	}
	*(*uint64)(unsafe.Pointer(dst)) = v
}

func copyIntToPUint64(dst, src unsafe.Pointer) {
	v := uint64(*(*int)(unsafe.Pointer(src)))
	p := (**uint64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPIntToPUint64(dst, src unsafe.Pointer) {
	var v uint64
	if p := *(**int)(unsafe.Pointer(src)); p != nil {
		v = uint64(*p)
	}

	p := (**uint64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int8 to uint64

func copyInt8ToUint64(dst, src unsafe.Pointer) {
	*(*uint64)(unsafe.Pointer(dst)) = uint64(*(*int8)(unsafe.Pointer(src)))
}

func copyPInt8ToUint64(dst, src unsafe.Pointer) {
	var v uint64
	if p := *(**int8)(unsafe.Pointer(src)); p != nil {
		v = uint64(*p)
	}
	*(*uint64)(unsafe.Pointer(dst)) = v
}

func copyInt8ToPUint64(dst, src unsafe.Pointer) {
	v := uint64(*(*int8)(unsafe.Pointer(src)))
	p := (**uint64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt8ToPUint64(dst, src unsafe.Pointer) {
	var v uint64
	if p := *(**int8)(unsafe.Pointer(src)); p != nil {
		v = uint64(*p)
	}

	p := (**uint64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int16 to uint64

func copyInt16ToUint64(dst, src unsafe.Pointer) {
	*(*uint64)(unsafe.Pointer(dst)) = uint64(*(*int16)(unsafe.Pointer(src)))
}

func copyPInt16ToUint64(dst, src unsafe.Pointer) {
	var v uint64
	if p := *(**int16)(unsafe.Pointer(src)); p != nil {
		v = uint64(*p)
	}
	*(*uint64)(unsafe.Pointer(dst)) = v
}

func copyInt16ToPUint64(dst, src unsafe.Pointer) {
	v := uint64(*(*int16)(unsafe.Pointer(src)))
	p := (**uint64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt16ToPUint64(dst, src unsafe.Pointer) {
	var v uint64
	if p := *(**int16)(unsafe.Pointer(src)); p != nil {
		v = uint64(*p)
	}

	p := (**uint64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int32 to uint64

func copyInt32ToUint64(dst, src unsafe.Pointer) {
	*(*uint64)(unsafe.Pointer(dst)) = uint64(*(*int32)(unsafe.Pointer(src)))
}

func copyPInt32ToUint64(dst, src unsafe.Pointer) {
	var v uint64
	if p := *(**int32)(unsafe.Pointer(src)); p != nil {
		v = uint64(*p)
	}
	*(*uint64)(unsafe.Pointer(dst)) = v
}

func copyInt32ToPUint64(dst, src unsafe.Pointer) {
	v := uint64(*(*int32)(unsafe.Pointer(src)))
	p := (**uint64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt32ToPUint64(dst, src unsafe.Pointer) {
	var v uint64
	if p := *(**int32)(unsafe.Pointer(src)); p != nil {
		v = uint64(*p)
	}

	p := (**uint64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// int64 to uint64

func copyInt64ToUint64(dst, src unsafe.Pointer) {
	*(*uint64)(unsafe.Pointer(dst)) = uint64(*(*int64)(unsafe.Pointer(src)))
}

func copyPInt64ToUint64(dst, src unsafe.Pointer) {
	var v uint64
	if p := *(**int64)(unsafe.Pointer(src)); p != nil {
		v = uint64(*p)
	}
	*(*uint64)(unsafe.Pointer(dst)) = v
}

func copyInt64ToPUint64(dst, src unsafe.Pointer) {
	v := uint64(*(*int64)(unsafe.Pointer(src)))
	p := (**uint64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPInt64ToPUint64(dst, src unsafe.Pointer) {
	var v uint64
	if p := *(**int64)(unsafe.Pointer(src)); p != nil {
		v = uint64(*p)
	}

	p := (**uint64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint to uint64

func copyUintToUint64(dst, src unsafe.Pointer) {
	*(*uint64)(unsafe.Pointer(dst)) = uint64(*(*uint)(unsafe.Pointer(src)))
}

func copyPUintToUint64(dst, src unsafe.Pointer) {
	var v uint64
	if p := *(**uint)(unsafe.Pointer(src)); p != nil {
		v = uint64(*p)
	}
	*(*uint64)(unsafe.Pointer(dst)) = v
}

func copyUintToPUint64(dst, src unsafe.Pointer) {
	v := uint64(*(*uint)(unsafe.Pointer(src)))
	p := (**uint64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUintToPUint64(dst, src unsafe.Pointer) {
	var v uint64
	if p := *(**uint)(unsafe.Pointer(src)); p != nil {
		v = uint64(*p)
	}

	p := (**uint64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint8 to uint64

func copyUint8ToUint64(dst, src unsafe.Pointer) {
	*(*uint64)(unsafe.Pointer(dst)) = uint64(*(*uint8)(unsafe.Pointer(src)))
}

func copyPUint8ToUint64(dst, src unsafe.Pointer) {
	var v uint64
	if p := *(**uint8)(unsafe.Pointer(src)); p != nil {
		v = uint64(*p)
	}
	*(*uint64)(unsafe.Pointer(dst)) = v
}

func copyUint8ToPUint64(dst, src unsafe.Pointer) {
	v := uint64(*(*uint8)(unsafe.Pointer(src)))
	p := (**uint64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint8ToPUint64(dst, src unsafe.Pointer) {
	var v uint64
	if p := *(**uint8)(unsafe.Pointer(src)); p != nil {
		v = uint64(*p)
	}

	p := (**uint64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint16 to uint64

func copyUint16ToUint64(dst, src unsafe.Pointer) {
	*(*uint64)(unsafe.Pointer(dst)) = uint64(*(*uint16)(unsafe.Pointer(src)))
}

func copyPUint16ToUint64(dst, src unsafe.Pointer) {
	var v uint64
	if p := *(**uint16)(unsafe.Pointer(src)); p != nil {
		v = uint64(*p)
	}
	*(*uint64)(unsafe.Pointer(dst)) = v
}

func copyUint16ToPUint64(dst, src unsafe.Pointer) {
	v := uint64(*(*uint16)(unsafe.Pointer(src)))
	p := (**uint64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint16ToPUint64(dst, src unsafe.Pointer) {
	var v uint64
	if p := *(**uint16)(unsafe.Pointer(src)); p != nil {
		v = uint64(*p)
	}

	p := (**uint64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint32 to uint64

func copyUint32ToUint64(dst, src unsafe.Pointer) {
	*(*uint64)(unsafe.Pointer(dst)) = uint64(*(*uint32)(unsafe.Pointer(src)))
}

func copyPUint32ToUint64(dst, src unsafe.Pointer) {
	var v uint64
	if p := *(**uint32)(unsafe.Pointer(src)); p != nil {
		v = uint64(*p)
	}
	*(*uint64)(unsafe.Pointer(dst)) = v
}

func copyUint32ToPUint64(dst, src unsafe.Pointer) {
	v := uint64(*(*uint32)(unsafe.Pointer(src)))
	p := (**uint64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint32ToPUint64(dst, src unsafe.Pointer) {
	var v uint64
	if p := *(**uint32)(unsafe.Pointer(src)); p != nil {
		v = uint64(*p)
	}

	p := (**uint64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// uint64 to uint64

func copyUint64ToUint64(dst, src unsafe.Pointer) {
	*(*uint64)(unsafe.Pointer(dst)) = uint64(*(*uint64)(unsafe.Pointer(src)))
}

func copyPUint64ToUint64(dst, src unsafe.Pointer) {
	var v uint64
	if p := *(**uint64)(unsafe.Pointer(src)); p != nil {
		v = uint64(*p)
	}
	*(*uint64)(unsafe.Pointer(dst)) = v
}

func copyUint64ToPUint64(dst, src unsafe.Pointer) {
	v := uint64(*(*uint64)(unsafe.Pointer(src)))
	p := (**uint64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPUint64ToPUint64(dst, src unsafe.Pointer) {
	var v uint64
	if p := *(**uint64)(unsafe.Pointer(src)); p != nil {
		v = uint64(*p)
	}

	p := (**uint64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// float32 to float32

func copyFloat32ToFloat32(dst, src unsafe.Pointer) {
	*(*float32)(unsafe.Pointer(dst)) = float32(*(*float32)(unsafe.Pointer(src)))
}

func copyPFloat32ToFloat32(dst, src unsafe.Pointer) {
	var v float32
	if p := *(**float32)(unsafe.Pointer(src)); p != nil {
		v = float32(*p)
	}
	*(*float32)(unsafe.Pointer(dst)) = v
}

func copyFloat32ToPFloat32(dst, src unsafe.Pointer) {
	v := float32(*(*float32)(unsafe.Pointer(src)))
	p := (**float32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPFloat32ToPFloat32(dst, src unsafe.Pointer) {
	var v float32
	if p := *(**float32)(unsafe.Pointer(src)); p != nil {
		v = float32(*p)
	}

	p := (**float32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// float64 to float32

func copyFloat64ToFloat32(dst, src unsafe.Pointer) {
	*(*float32)(unsafe.Pointer(dst)) = float32(*(*float64)(unsafe.Pointer(src)))
}

func copyPFloat64ToFloat32(dst, src unsafe.Pointer) {
	var v float32
	if p := *(**float64)(unsafe.Pointer(src)); p != nil {
		v = float32(*p)
	}
	*(*float32)(unsafe.Pointer(dst)) = v
}

func copyFloat64ToPFloat32(dst, src unsafe.Pointer) {
	v := float32(*(*float64)(unsafe.Pointer(src)))
	p := (**float32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPFloat64ToPFloat32(dst, src unsafe.Pointer) {
	var v float32
	if p := *(**float64)(unsafe.Pointer(src)); p != nil {
		v = float32(*p)
	}

	p := (**float32)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// float32 to float64

func copyFloat32ToFloat64(dst, src unsafe.Pointer) {
	*(*float64)(unsafe.Pointer(dst)) = float64(*(*float32)(unsafe.Pointer(src)))
}

func copyPFloat32ToFloat64(dst, src unsafe.Pointer) {
	var v float64
	if p := *(**float32)(unsafe.Pointer(src)); p != nil {
		v = float64(*p)
	}
	*(*float64)(unsafe.Pointer(dst)) = v
}

func copyFloat32ToPFloat64(dst, src unsafe.Pointer) {
	v := float64(*(*float32)(unsafe.Pointer(src)))
	p := (**float64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPFloat32ToPFloat64(dst, src unsafe.Pointer) {
	var v float64
	if p := *(**float32)(unsafe.Pointer(src)); p != nil {
		v = float64(*p)
	}

	p := (**float64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// float64 to float64

func copyFloat64ToFloat64(dst, src unsafe.Pointer) {
	*(*float64)(unsafe.Pointer(dst)) = float64(*(*float64)(unsafe.Pointer(src)))
}

func copyPFloat64ToFloat64(dst, src unsafe.Pointer) {
	var v float64
	if p := *(**float64)(unsafe.Pointer(src)); p != nil {
		v = float64(*p)
	}
	*(*float64)(unsafe.Pointer(dst)) = v
}

func copyFloat64ToPFloat64(dst, src unsafe.Pointer) {
	v := float64(*(*float64)(unsafe.Pointer(src)))
	p := (**float64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPFloat64ToPFloat64(dst, src unsafe.Pointer) {
	var v float64
	if p := *(**float64)(unsafe.Pointer(src)); p != nil {
		v = float64(*p)
	}

	p := (**float64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// bool to bool

func copyBoolToBool(dst, src unsafe.Pointer) {
	*(*bool)(unsafe.Pointer(dst)) = bool(*(*bool)(unsafe.Pointer(src)))
}

func copyPBoolToBool(dst, src unsafe.Pointer) {
	var v bool
	if p := *(**bool)(unsafe.Pointer(src)); p != nil {
		v = bool(*p)
	}
	*(*bool)(unsafe.Pointer(dst)) = v
}

func copyBoolToPBool(dst, src unsafe.Pointer) {
	v := bool(*(*bool)(unsafe.Pointer(src)))
	p := (**bool)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPBoolToPBool(dst, src unsafe.Pointer) {
	var v bool
	if p := *(**bool)(unsafe.Pointer(src)); p != nil {
		v = bool(*p)
	}

	p := (**bool)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// complex64 to complex64

func copyComplex64ToComplex64(dst, src unsafe.Pointer) {
	*(*complex64)(unsafe.Pointer(dst)) = complex64(*(*complex64)(unsafe.Pointer(src)))
}

func copyPComplex64ToComplex64(dst, src unsafe.Pointer) {
	var v complex64
	if p := *(**complex64)(unsafe.Pointer(src)); p != nil {
		v = complex64(*p)
	}
	*(*complex64)(unsafe.Pointer(dst)) = v
}

func copyComplex64ToPComplex64(dst, src unsafe.Pointer) {
	v := complex64(*(*complex64)(unsafe.Pointer(src)))
	p := (**complex64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPComplex64ToPComplex64(dst, src unsafe.Pointer) {
	var v complex64
	if p := *(**complex64)(unsafe.Pointer(src)); p != nil {
		v = complex64(*p)
	}

	p := (**complex64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// complex128 to complex64

func copyComplex128ToComplex64(dst, src unsafe.Pointer) {
	*(*complex64)(unsafe.Pointer(dst)) = complex64(*(*complex128)(unsafe.Pointer(src)))
}

func copyPComplex128ToComplex64(dst, src unsafe.Pointer) {
	var v complex64
	if p := *(**complex128)(unsafe.Pointer(src)); p != nil {
		v = complex64(*p)
	}
	*(*complex64)(unsafe.Pointer(dst)) = v
}

func copyComplex128ToPComplex64(dst, src unsafe.Pointer) {
	v := complex64(*(*complex128)(unsafe.Pointer(src)))
	p := (**complex64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPComplex128ToPComplex64(dst, src unsafe.Pointer) {
	var v complex64
	if p := *(**complex128)(unsafe.Pointer(src)); p != nil {
		v = complex64(*p)
	}

	p := (**complex64)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// complex64 to complex128

func copyComplex64ToComplex128(dst, src unsafe.Pointer) {
	*(*complex128)(unsafe.Pointer(dst)) = complex128(*(*complex64)(unsafe.Pointer(src)))
}

func copyPComplex64ToComplex128(dst, src unsafe.Pointer) {
	var v complex128
	if p := *(**complex64)(unsafe.Pointer(src)); p != nil {
		v = complex128(*p)
	}
	*(*complex128)(unsafe.Pointer(dst)) = v
}

func copyComplex64ToPComplex128(dst, src unsafe.Pointer) {
	v := complex128(*(*complex64)(unsafe.Pointer(src)))
	p := (**complex128)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPComplex64ToPComplex128(dst, src unsafe.Pointer) {
	var v complex128
	if p := *(**complex64)(unsafe.Pointer(src)); p != nil {
		v = complex128(*p)
	}

	p := (**complex128)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// complex128 to complex128

func copyComplex128ToComplex128(dst, src unsafe.Pointer) {
	*(*complex128)(unsafe.Pointer(dst)) = complex128(*(*complex128)(unsafe.Pointer(src)))
}

func copyPComplex128ToComplex128(dst, src unsafe.Pointer) {
	var v complex128
	if p := *(**complex128)(unsafe.Pointer(src)); p != nil {
		v = complex128(*p)
	}
	*(*complex128)(unsafe.Pointer(dst)) = v
}

func copyComplex128ToPComplex128(dst, src unsafe.Pointer) {
	v := complex128(*(*complex128)(unsafe.Pointer(src)))
	p := (**complex128)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPComplex128ToPComplex128(dst, src unsafe.Pointer) {
	var v complex128
	if p := *(**complex128)(unsafe.Pointer(src)); p != nil {
		v = complex128(*p)
	}

	p := (**complex128)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// string to string

func copyStringToString(dst, src unsafe.Pointer) {
	*(*string)(unsafe.Pointer(dst)) = string(*(*string)(unsafe.Pointer(src)))
}

func copyPStringToString(dst, src unsafe.Pointer) {
	var v string
	if p := *(**string)(unsafe.Pointer(src)); p != nil {
		v = string(*p)
	}
	*(*string)(unsafe.Pointer(dst)) = v
}

func copyStringToPString(dst, src unsafe.Pointer) {
	v := string(*(*string)(unsafe.Pointer(src)))
	p := (**string)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPStringToPString(dst, src unsafe.Pointer) {
	var v string
	if p := *(**string)(unsafe.Pointer(src)); p != nil {
		v = string(*p)
	}

	p := (**string)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// []byte to string

func copyBytesToString(dst, src unsafe.Pointer) {
	*(*string)(unsafe.Pointer(dst)) = string(*(*[]byte)(unsafe.Pointer(src)))
}

func copyPBytesToString(dst, src unsafe.Pointer) {
	var v string
	if p := *(**[]byte)(unsafe.Pointer(src)); p != nil {
		v = string(*p)
	}
	*(*string)(unsafe.Pointer(dst)) = v
}

func copyBytesToPString(dst, src unsafe.Pointer) {
	v := string(*(*[]byte)(unsafe.Pointer(src)))
	p := (**string)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPBytesToPString(dst, src unsafe.Pointer) {
	var v string
	if p := *(**[]byte)(unsafe.Pointer(src)); p != nil {
		v = string(*p)
	}

	p := (**string)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// string to []byte

func copyStringToBytes(dst, src unsafe.Pointer) {
	*(*[]byte)(unsafe.Pointer(dst)) = []byte(*(*string)(unsafe.Pointer(src)))
}

func copyPStringToBytes(dst, src unsafe.Pointer) {
	var v []byte
	if p := *(**string)(unsafe.Pointer(src)); p != nil {
		v = []byte(*p)
	}
	*(*[]byte)(unsafe.Pointer(dst)) = v
}

func copyStringToPBytes(dst, src unsafe.Pointer) {
	v := []byte(*(*string)(unsafe.Pointer(src)))
	p := (**[]byte)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPStringToPBytes(dst, src unsafe.Pointer) {
	var v []byte
	if p := *(**string)(unsafe.Pointer(src)); p != nil {
		v = []byte(*p)
	}

	p := (**[]byte)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// []byte to []byte

func copyBytesToBytes(dst, src unsafe.Pointer) {
	*(*[]byte)(unsafe.Pointer(dst)) = []byte(*(*[]byte)(unsafe.Pointer(src)))
}

func copyPBytesToBytes(dst, src unsafe.Pointer) {
	var v []byte
	if p := *(**[]byte)(unsafe.Pointer(src)); p != nil {
		v = []byte(*p)
	}
	*(*[]byte)(unsafe.Pointer(dst)) = v
}

func copyBytesToPBytes(dst, src unsafe.Pointer) {
	v := []byte(*(*[]byte)(unsafe.Pointer(src)))
	p := (**[]byte)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPBytesToPBytes(dst, src unsafe.Pointer) {
	var v []byte
	if p := *(**[]byte)(unsafe.Pointer(src)); p != nil {
		v = []byte(*p)
	}

	p := (**[]byte)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// time.Time to time.Time

func copyTimeToTime(dst, src unsafe.Pointer) {
	*(*time.Time)(unsafe.Pointer(dst)) = time.Time(*(*time.Time)(unsafe.Pointer(src)))
}

func copyPTimeToTime(dst, src unsafe.Pointer) {
	var v time.Time
	if p := *(**time.Time)(unsafe.Pointer(src)); p != nil {
		v = time.Time(*p)
	}
	*(*time.Time)(unsafe.Pointer(dst)) = v
}

func copyTimeToPTime(dst, src unsafe.Pointer) {
	v := time.Time(*(*time.Time)(unsafe.Pointer(src)))
	p := (**time.Time)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPTimeToPTime(dst, src unsafe.Pointer) {
	var v time.Time
	if p := *(**time.Time)(unsafe.Pointer(src)); p != nil {
		v = time.Time(*p)
	}

	p := (**time.Time)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// time.Duration to time.Duration

func copyDurationToDuration(dst, src unsafe.Pointer) {
	*(*time.Duration)(unsafe.Pointer(dst)) = time.Duration(*(*time.Duration)(unsafe.Pointer(src)))
}

func copyPDurationToDuration(dst, src unsafe.Pointer) {
	var v time.Duration
	if p := *(**time.Duration)(unsafe.Pointer(src)); p != nil {
		v = time.Duration(*p)
	}
	*(*time.Duration)(unsafe.Pointer(dst)) = v
}

func copyDurationToPDuration(dst, src unsafe.Pointer) {
	v := time.Duration(*(*time.Duration)(unsafe.Pointer(src)))
	p := (**time.Duration)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

func copyPDurationToPDuration(dst, src unsafe.Pointer) {
	var v time.Duration
	if p := *(**time.Duration)(unsafe.Pointer(src)); p != nil {
		v = time.Duration(*p)
	}

	p := (**time.Duration)(unsafe.Pointer(dst))
	if p := *p; p != nil {
		*p = v
		return
	}
	*p = &v
}

// Memcopy funcs.
func copy1(dst, src unsafe.Pointer) {
	*(*[1]byte)(unsafe.Pointer(dst)) = *(*[1]byte)(unsafe.Pointer(src))
}

func copy2(dst, src unsafe.Pointer) {
	*(*[2]byte)(unsafe.Pointer(dst)) = *(*[2]byte)(unsafe.Pointer(src))
}

func copy3(dst, src unsafe.Pointer) {
	*(*[3]byte)(unsafe.Pointer(dst)) = *(*[3]byte)(unsafe.Pointer(src))
}

func copy4(dst, src unsafe.Pointer) {
	*(*[4]byte)(unsafe.Pointer(dst)) = *(*[4]byte)(unsafe.Pointer(src))
}

func copy5(dst, src unsafe.Pointer) {
	*(*[5]byte)(unsafe.Pointer(dst)) = *(*[5]byte)(unsafe.Pointer(src))
}

func copy6(dst, src unsafe.Pointer) {
	*(*[6]byte)(unsafe.Pointer(dst)) = *(*[6]byte)(unsafe.Pointer(src))
}

func copy7(dst, src unsafe.Pointer) {
	*(*[7]byte)(unsafe.Pointer(dst)) = *(*[7]byte)(unsafe.Pointer(src))
}

func copy8(dst, src unsafe.Pointer) {
	*(*[8]byte)(unsafe.Pointer(dst)) = *(*[8]byte)(unsafe.Pointer(src))
}

func copy9(dst, src unsafe.Pointer) {
	*(*[9]byte)(unsafe.Pointer(dst)) = *(*[9]byte)(unsafe.Pointer(src))
}

func copy10(dst, src unsafe.Pointer) {
	*(*[10]byte)(unsafe.Pointer(dst)) = *(*[10]byte)(unsafe.Pointer(src))
}

func copy11(dst, src unsafe.Pointer) {
	*(*[11]byte)(unsafe.Pointer(dst)) = *(*[11]byte)(unsafe.Pointer(src))
}

func copy12(dst, src unsafe.Pointer) {
	*(*[12]byte)(unsafe.Pointer(dst)) = *(*[12]byte)(unsafe.Pointer(src))
}

func copy13(dst, src unsafe.Pointer) {
	*(*[13]byte)(unsafe.Pointer(dst)) = *(*[13]byte)(unsafe.Pointer(src))
}

func copy14(dst, src unsafe.Pointer) {
	*(*[14]byte)(unsafe.Pointer(dst)) = *(*[14]byte)(unsafe.Pointer(src))
}

func copy15(dst, src unsafe.Pointer) {
	*(*[15]byte)(unsafe.Pointer(dst)) = *(*[15]byte)(unsafe.Pointer(src))
}

func copy16(dst, src unsafe.Pointer) {
	*(*[16]byte)(unsafe.Pointer(dst)) = *(*[16]byte)(unsafe.Pointer(src))
}

func copy17(dst, src unsafe.Pointer) {
	*(*[17]byte)(unsafe.Pointer(dst)) = *(*[17]byte)(unsafe.Pointer(src))
}

func copy18(dst, src unsafe.Pointer) {
	*(*[18]byte)(unsafe.Pointer(dst)) = *(*[18]byte)(unsafe.Pointer(src))
}

func copy19(dst, src unsafe.Pointer) {
	*(*[19]byte)(unsafe.Pointer(dst)) = *(*[19]byte)(unsafe.Pointer(src))
}

func copy20(dst, src unsafe.Pointer) {
	*(*[20]byte)(unsafe.Pointer(dst)) = *(*[20]byte)(unsafe.Pointer(src))
}

func copy21(dst, src unsafe.Pointer) {
	*(*[21]byte)(unsafe.Pointer(dst)) = *(*[21]byte)(unsafe.Pointer(src))
}

func copy22(dst, src unsafe.Pointer) {
	*(*[22]byte)(unsafe.Pointer(dst)) = *(*[22]byte)(unsafe.Pointer(src))
}

func copy23(dst, src unsafe.Pointer) {
	*(*[23]byte)(unsafe.Pointer(dst)) = *(*[23]byte)(unsafe.Pointer(src))
}

func copy24(dst, src unsafe.Pointer) {
	*(*[24]byte)(unsafe.Pointer(dst)) = *(*[24]byte)(unsafe.Pointer(src))
}

func copy25(dst, src unsafe.Pointer) {
	*(*[25]byte)(unsafe.Pointer(dst)) = *(*[25]byte)(unsafe.Pointer(src))
}

func copy26(dst, src unsafe.Pointer) {
	*(*[26]byte)(unsafe.Pointer(dst)) = *(*[26]byte)(unsafe.Pointer(src))
}

func copy27(dst, src unsafe.Pointer) {
	*(*[27]byte)(unsafe.Pointer(dst)) = *(*[27]byte)(unsafe.Pointer(src))
}

func copy28(dst, src unsafe.Pointer) {
	*(*[28]byte)(unsafe.Pointer(dst)) = *(*[28]byte)(unsafe.Pointer(src))
}

func copy29(dst, src unsafe.Pointer) {
	*(*[29]byte)(unsafe.Pointer(dst)) = *(*[29]byte)(unsafe.Pointer(src))
}

func copy30(dst, src unsafe.Pointer) {
	*(*[30]byte)(unsafe.Pointer(dst)) = *(*[30]byte)(unsafe.Pointer(src))
}

func copy31(dst, src unsafe.Pointer) {
	*(*[31]byte)(unsafe.Pointer(dst)) = *(*[31]byte)(unsafe.Pointer(src))
}

func copy32(dst, src unsafe.Pointer) {
	*(*[32]byte)(unsafe.Pointer(dst)) = *(*[32]byte)(unsafe.Pointer(src))
}

func copy33(dst, src unsafe.Pointer) {
	*(*[33]byte)(unsafe.Pointer(dst)) = *(*[33]byte)(unsafe.Pointer(src))
}

func copy34(dst, src unsafe.Pointer) {
	*(*[34]byte)(unsafe.Pointer(dst)) = *(*[34]byte)(unsafe.Pointer(src))
}

func copy35(dst, src unsafe.Pointer) {
	*(*[35]byte)(unsafe.Pointer(dst)) = *(*[35]byte)(unsafe.Pointer(src))
}

func copy36(dst, src unsafe.Pointer) {
	*(*[36]byte)(unsafe.Pointer(dst)) = *(*[36]byte)(unsafe.Pointer(src))
}

func copy37(dst, src unsafe.Pointer) {
	*(*[37]byte)(unsafe.Pointer(dst)) = *(*[37]byte)(unsafe.Pointer(src))
}

func copy38(dst, src unsafe.Pointer) {
	*(*[38]byte)(unsafe.Pointer(dst)) = *(*[38]byte)(unsafe.Pointer(src))
}

func copy39(dst, src unsafe.Pointer) {
	*(*[39]byte)(unsafe.Pointer(dst)) = *(*[39]byte)(unsafe.Pointer(src))
}

func copy40(dst, src unsafe.Pointer) {
	*(*[40]byte)(unsafe.Pointer(dst)) = *(*[40]byte)(unsafe.Pointer(src))
}

func copy41(dst, src unsafe.Pointer) {
	*(*[41]byte)(unsafe.Pointer(dst)) = *(*[41]byte)(unsafe.Pointer(src))
}

func copy42(dst, src unsafe.Pointer) {
	*(*[42]byte)(unsafe.Pointer(dst)) = *(*[42]byte)(unsafe.Pointer(src))
}

func copy43(dst, src unsafe.Pointer) {
	*(*[43]byte)(unsafe.Pointer(dst)) = *(*[43]byte)(unsafe.Pointer(src))
}

func copy44(dst, src unsafe.Pointer) {
	*(*[44]byte)(unsafe.Pointer(dst)) = *(*[44]byte)(unsafe.Pointer(src))
}

func copy45(dst, src unsafe.Pointer) {
	*(*[45]byte)(unsafe.Pointer(dst)) = *(*[45]byte)(unsafe.Pointer(src))
}

func copy46(dst, src unsafe.Pointer) {
	*(*[46]byte)(unsafe.Pointer(dst)) = *(*[46]byte)(unsafe.Pointer(src))
}

func copy47(dst, src unsafe.Pointer) {
	*(*[47]byte)(unsafe.Pointer(dst)) = *(*[47]byte)(unsafe.Pointer(src))
}

func copy48(dst, src unsafe.Pointer) {
	*(*[48]byte)(unsafe.Pointer(dst)) = *(*[48]byte)(unsafe.Pointer(src))
}

func copy49(dst, src unsafe.Pointer) {
	*(*[49]byte)(unsafe.Pointer(dst)) = *(*[49]byte)(unsafe.Pointer(src))
}

func copy50(dst, src unsafe.Pointer) {
	*(*[50]byte)(unsafe.Pointer(dst)) = *(*[50]byte)(unsafe.Pointer(src))
}

func copy51(dst, src unsafe.Pointer) {
	*(*[51]byte)(unsafe.Pointer(dst)) = *(*[51]byte)(unsafe.Pointer(src))
}

func copy52(dst, src unsafe.Pointer) {
	*(*[52]byte)(unsafe.Pointer(dst)) = *(*[52]byte)(unsafe.Pointer(src))
}

func copy53(dst, src unsafe.Pointer) {
	*(*[53]byte)(unsafe.Pointer(dst)) = *(*[53]byte)(unsafe.Pointer(src))
}

func copy54(dst, src unsafe.Pointer) {
	*(*[54]byte)(unsafe.Pointer(dst)) = *(*[54]byte)(unsafe.Pointer(src))
}

func copy55(dst, src unsafe.Pointer) {
	*(*[55]byte)(unsafe.Pointer(dst)) = *(*[55]byte)(unsafe.Pointer(src))
}

func copy56(dst, src unsafe.Pointer) {
	*(*[56]byte)(unsafe.Pointer(dst)) = *(*[56]byte)(unsafe.Pointer(src))
}

func copy57(dst, src unsafe.Pointer) {
	*(*[57]byte)(unsafe.Pointer(dst)) = *(*[57]byte)(unsafe.Pointer(src))
}

func copy58(dst, src unsafe.Pointer) {
	*(*[58]byte)(unsafe.Pointer(dst)) = *(*[58]byte)(unsafe.Pointer(src))
}

func copy59(dst, src unsafe.Pointer) {
	*(*[59]byte)(unsafe.Pointer(dst)) = *(*[59]byte)(unsafe.Pointer(src))
}

func copy60(dst, src unsafe.Pointer) {
	*(*[60]byte)(unsafe.Pointer(dst)) = *(*[60]byte)(unsafe.Pointer(src))
}

func copy61(dst, src unsafe.Pointer) {
	*(*[61]byte)(unsafe.Pointer(dst)) = *(*[61]byte)(unsafe.Pointer(src))
}

func copy62(dst, src unsafe.Pointer) {
	*(*[62]byte)(unsafe.Pointer(dst)) = *(*[62]byte)(unsafe.Pointer(src))
}

func copy63(dst, src unsafe.Pointer) {
	*(*[63]byte)(unsafe.Pointer(dst)) = *(*[63]byte)(unsafe.Pointer(src))
}

func copy64(dst, src unsafe.Pointer) {
	*(*[64]byte)(unsafe.Pointer(dst)) = *(*[64]byte)(unsafe.Pointer(src))
}

func copy65(dst, src unsafe.Pointer) {
	*(*[65]byte)(unsafe.Pointer(dst)) = *(*[65]byte)(unsafe.Pointer(src))
}

func copy66(dst, src unsafe.Pointer) {
	*(*[66]byte)(unsafe.Pointer(dst)) = *(*[66]byte)(unsafe.Pointer(src))
}

func copy67(dst, src unsafe.Pointer) {
	*(*[67]byte)(unsafe.Pointer(dst)) = *(*[67]byte)(unsafe.Pointer(src))
}

func copy68(dst, src unsafe.Pointer) {
	*(*[68]byte)(unsafe.Pointer(dst)) = *(*[68]byte)(unsafe.Pointer(src))
}

func copy69(dst, src unsafe.Pointer) {
	*(*[69]byte)(unsafe.Pointer(dst)) = *(*[69]byte)(unsafe.Pointer(src))
}

func copy70(dst, src unsafe.Pointer) {
	*(*[70]byte)(unsafe.Pointer(dst)) = *(*[70]byte)(unsafe.Pointer(src))
}

func copy71(dst, src unsafe.Pointer) {
	*(*[71]byte)(unsafe.Pointer(dst)) = *(*[71]byte)(unsafe.Pointer(src))
}

func copy72(dst, src unsafe.Pointer) {
	*(*[72]byte)(unsafe.Pointer(dst)) = *(*[72]byte)(unsafe.Pointer(src))
}

func copy73(dst, src unsafe.Pointer) {
	*(*[73]byte)(unsafe.Pointer(dst)) = *(*[73]byte)(unsafe.Pointer(src))
}

func copy74(dst, src unsafe.Pointer) {
	*(*[74]byte)(unsafe.Pointer(dst)) = *(*[74]byte)(unsafe.Pointer(src))
}

func copy75(dst, src unsafe.Pointer) {
	*(*[75]byte)(unsafe.Pointer(dst)) = *(*[75]byte)(unsafe.Pointer(src))
}

func copy76(dst, src unsafe.Pointer) {
	*(*[76]byte)(unsafe.Pointer(dst)) = *(*[76]byte)(unsafe.Pointer(src))
}

func copy77(dst, src unsafe.Pointer) {
	*(*[77]byte)(unsafe.Pointer(dst)) = *(*[77]byte)(unsafe.Pointer(src))
}

func copy78(dst, src unsafe.Pointer) {
	*(*[78]byte)(unsafe.Pointer(dst)) = *(*[78]byte)(unsafe.Pointer(src))
}

func copy79(dst, src unsafe.Pointer) {
	*(*[79]byte)(unsafe.Pointer(dst)) = *(*[79]byte)(unsafe.Pointer(src))
}

func copy80(dst, src unsafe.Pointer) {
	*(*[80]byte)(unsafe.Pointer(dst)) = *(*[80]byte)(unsafe.Pointer(src))
}

func copy81(dst, src unsafe.Pointer) {
	*(*[81]byte)(unsafe.Pointer(dst)) = *(*[81]byte)(unsafe.Pointer(src))
}

func copy82(dst, src unsafe.Pointer) {
	*(*[82]byte)(unsafe.Pointer(dst)) = *(*[82]byte)(unsafe.Pointer(src))
}

func copy83(dst, src unsafe.Pointer) {
	*(*[83]byte)(unsafe.Pointer(dst)) = *(*[83]byte)(unsafe.Pointer(src))
}

func copy84(dst, src unsafe.Pointer) {
	*(*[84]byte)(unsafe.Pointer(dst)) = *(*[84]byte)(unsafe.Pointer(src))
}

func copy85(dst, src unsafe.Pointer) {
	*(*[85]byte)(unsafe.Pointer(dst)) = *(*[85]byte)(unsafe.Pointer(src))
}

func copy86(dst, src unsafe.Pointer) {
	*(*[86]byte)(unsafe.Pointer(dst)) = *(*[86]byte)(unsafe.Pointer(src))
}

func copy87(dst, src unsafe.Pointer) {
	*(*[87]byte)(unsafe.Pointer(dst)) = *(*[87]byte)(unsafe.Pointer(src))
}

func copy88(dst, src unsafe.Pointer) {
	*(*[88]byte)(unsafe.Pointer(dst)) = *(*[88]byte)(unsafe.Pointer(src))
}

func copy89(dst, src unsafe.Pointer) {
	*(*[89]byte)(unsafe.Pointer(dst)) = *(*[89]byte)(unsafe.Pointer(src))
}

func copy90(dst, src unsafe.Pointer) {
	*(*[90]byte)(unsafe.Pointer(dst)) = *(*[90]byte)(unsafe.Pointer(src))
}

func copy91(dst, src unsafe.Pointer) {
	*(*[91]byte)(unsafe.Pointer(dst)) = *(*[91]byte)(unsafe.Pointer(src))
}

func copy92(dst, src unsafe.Pointer) {
	*(*[92]byte)(unsafe.Pointer(dst)) = *(*[92]byte)(unsafe.Pointer(src))
}

func copy93(dst, src unsafe.Pointer) {
	*(*[93]byte)(unsafe.Pointer(dst)) = *(*[93]byte)(unsafe.Pointer(src))
}

func copy94(dst, src unsafe.Pointer) {
	*(*[94]byte)(unsafe.Pointer(dst)) = *(*[94]byte)(unsafe.Pointer(src))
}

func copy95(dst, src unsafe.Pointer) {
	*(*[95]byte)(unsafe.Pointer(dst)) = *(*[95]byte)(unsafe.Pointer(src))
}

func copy96(dst, src unsafe.Pointer) {
	*(*[96]byte)(unsafe.Pointer(dst)) = *(*[96]byte)(unsafe.Pointer(src))
}

func copy97(dst, src unsafe.Pointer) {
	*(*[97]byte)(unsafe.Pointer(dst)) = *(*[97]byte)(unsafe.Pointer(src))
}

func copy98(dst, src unsafe.Pointer) {
	*(*[98]byte)(unsafe.Pointer(dst)) = *(*[98]byte)(unsafe.Pointer(src))
}

func copy99(dst, src unsafe.Pointer) {
	*(*[99]byte)(unsafe.Pointer(dst)) = *(*[99]byte)(unsafe.Pointer(src))
}

func copy100(dst, src unsafe.Pointer) {
	*(*[100]byte)(unsafe.Pointer(dst)) = *(*[100]byte)(unsafe.Pointer(src))
}

func copy101(dst, src unsafe.Pointer) {
	*(*[101]byte)(unsafe.Pointer(dst)) = *(*[101]byte)(unsafe.Pointer(src))
}

func copy102(dst, src unsafe.Pointer) {
	*(*[102]byte)(unsafe.Pointer(dst)) = *(*[102]byte)(unsafe.Pointer(src))
}

func copy103(dst, src unsafe.Pointer) {
	*(*[103]byte)(unsafe.Pointer(dst)) = *(*[103]byte)(unsafe.Pointer(src))
}

func copy104(dst, src unsafe.Pointer) {
	*(*[104]byte)(unsafe.Pointer(dst)) = *(*[104]byte)(unsafe.Pointer(src))
}

func copy105(dst, src unsafe.Pointer) {
	*(*[105]byte)(unsafe.Pointer(dst)) = *(*[105]byte)(unsafe.Pointer(src))
}

func copy106(dst, src unsafe.Pointer) {
	*(*[106]byte)(unsafe.Pointer(dst)) = *(*[106]byte)(unsafe.Pointer(src))
}

func copy107(dst, src unsafe.Pointer) {
	*(*[107]byte)(unsafe.Pointer(dst)) = *(*[107]byte)(unsafe.Pointer(src))
}

func copy108(dst, src unsafe.Pointer) {
	*(*[108]byte)(unsafe.Pointer(dst)) = *(*[108]byte)(unsafe.Pointer(src))
}

func copy109(dst, src unsafe.Pointer) {
	*(*[109]byte)(unsafe.Pointer(dst)) = *(*[109]byte)(unsafe.Pointer(src))
}

func copy110(dst, src unsafe.Pointer) {
	*(*[110]byte)(unsafe.Pointer(dst)) = *(*[110]byte)(unsafe.Pointer(src))
}

func copy111(dst, src unsafe.Pointer) {
	*(*[111]byte)(unsafe.Pointer(dst)) = *(*[111]byte)(unsafe.Pointer(src))
}

func copy112(dst, src unsafe.Pointer) {
	*(*[112]byte)(unsafe.Pointer(dst)) = *(*[112]byte)(unsafe.Pointer(src))
}

func copy113(dst, src unsafe.Pointer) {
	*(*[113]byte)(unsafe.Pointer(dst)) = *(*[113]byte)(unsafe.Pointer(src))
}

func copy114(dst, src unsafe.Pointer) {
	*(*[114]byte)(unsafe.Pointer(dst)) = *(*[114]byte)(unsafe.Pointer(src))
}

func copy115(dst, src unsafe.Pointer) {
	*(*[115]byte)(unsafe.Pointer(dst)) = *(*[115]byte)(unsafe.Pointer(src))
}

func copy116(dst, src unsafe.Pointer) {
	*(*[116]byte)(unsafe.Pointer(dst)) = *(*[116]byte)(unsafe.Pointer(src))
}

func copy117(dst, src unsafe.Pointer) {
	*(*[117]byte)(unsafe.Pointer(dst)) = *(*[117]byte)(unsafe.Pointer(src))
}

func copy118(dst, src unsafe.Pointer) {
	*(*[118]byte)(unsafe.Pointer(dst)) = *(*[118]byte)(unsafe.Pointer(src))
}

func copy119(dst, src unsafe.Pointer) {
	*(*[119]byte)(unsafe.Pointer(dst)) = *(*[119]byte)(unsafe.Pointer(src))
}

func copy120(dst, src unsafe.Pointer) {
	*(*[120]byte)(unsafe.Pointer(dst)) = *(*[120]byte)(unsafe.Pointer(src))
}

func copy121(dst, src unsafe.Pointer) {
	*(*[121]byte)(unsafe.Pointer(dst)) = *(*[121]byte)(unsafe.Pointer(src))
}

func copy122(dst, src unsafe.Pointer) {
	*(*[122]byte)(unsafe.Pointer(dst)) = *(*[122]byte)(unsafe.Pointer(src))
}

func copy123(dst, src unsafe.Pointer) {
	*(*[123]byte)(unsafe.Pointer(dst)) = *(*[123]byte)(unsafe.Pointer(src))
}

func copy124(dst, src unsafe.Pointer) {
	*(*[124]byte)(unsafe.Pointer(dst)) = *(*[124]byte)(unsafe.Pointer(src))
}

func copy125(dst, src unsafe.Pointer) {
	*(*[125]byte)(unsafe.Pointer(dst)) = *(*[125]byte)(unsafe.Pointer(src))
}

func copy126(dst, src unsafe.Pointer) {
	*(*[126]byte)(unsafe.Pointer(dst)) = *(*[126]byte)(unsafe.Pointer(src))
}

func copy127(dst, src unsafe.Pointer) {
	*(*[127]byte)(unsafe.Pointer(dst)) = *(*[127]byte)(unsafe.Pointer(src))
}

func copy128(dst, src unsafe.Pointer) {
	*(*[128]byte)(unsafe.Pointer(dst)) = *(*[128]byte)(unsafe.Pointer(src))
}

func copy129(dst, src unsafe.Pointer) {
	*(*[129]byte)(unsafe.Pointer(dst)) = *(*[129]byte)(unsafe.Pointer(src))
}

func copy130(dst, src unsafe.Pointer) {
	*(*[130]byte)(unsafe.Pointer(dst)) = *(*[130]byte)(unsafe.Pointer(src))
}

func copy131(dst, src unsafe.Pointer) {
	*(*[131]byte)(unsafe.Pointer(dst)) = *(*[131]byte)(unsafe.Pointer(src))
}

func copy132(dst, src unsafe.Pointer) {
	*(*[132]byte)(unsafe.Pointer(dst)) = *(*[132]byte)(unsafe.Pointer(src))
}

func copy133(dst, src unsafe.Pointer) {
	*(*[133]byte)(unsafe.Pointer(dst)) = *(*[133]byte)(unsafe.Pointer(src))
}

func copy134(dst, src unsafe.Pointer) {
	*(*[134]byte)(unsafe.Pointer(dst)) = *(*[134]byte)(unsafe.Pointer(src))
}

func copy135(dst, src unsafe.Pointer) {
	*(*[135]byte)(unsafe.Pointer(dst)) = *(*[135]byte)(unsafe.Pointer(src))
}

func copy136(dst, src unsafe.Pointer) {
	*(*[136]byte)(unsafe.Pointer(dst)) = *(*[136]byte)(unsafe.Pointer(src))
}

func copy137(dst, src unsafe.Pointer) {
	*(*[137]byte)(unsafe.Pointer(dst)) = *(*[137]byte)(unsafe.Pointer(src))
}

func copy138(dst, src unsafe.Pointer) {
	*(*[138]byte)(unsafe.Pointer(dst)) = *(*[138]byte)(unsafe.Pointer(src))
}

func copy139(dst, src unsafe.Pointer) {
	*(*[139]byte)(unsafe.Pointer(dst)) = *(*[139]byte)(unsafe.Pointer(src))
}

func copy140(dst, src unsafe.Pointer) {
	*(*[140]byte)(unsafe.Pointer(dst)) = *(*[140]byte)(unsafe.Pointer(src))
}

func copy141(dst, src unsafe.Pointer) {
	*(*[141]byte)(unsafe.Pointer(dst)) = *(*[141]byte)(unsafe.Pointer(src))
}

func copy142(dst, src unsafe.Pointer) {
	*(*[142]byte)(unsafe.Pointer(dst)) = *(*[142]byte)(unsafe.Pointer(src))
}

func copy143(dst, src unsafe.Pointer) {
	*(*[143]byte)(unsafe.Pointer(dst)) = *(*[143]byte)(unsafe.Pointer(src))
}

func copy144(dst, src unsafe.Pointer) {
	*(*[144]byte)(unsafe.Pointer(dst)) = *(*[144]byte)(unsafe.Pointer(src))
}

func copy145(dst, src unsafe.Pointer) {
	*(*[145]byte)(unsafe.Pointer(dst)) = *(*[145]byte)(unsafe.Pointer(src))
}

func copy146(dst, src unsafe.Pointer) {
	*(*[146]byte)(unsafe.Pointer(dst)) = *(*[146]byte)(unsafe.Pointer(src))
}

func copy147(dst, src unsafe.Pointer) {
	*(*[147]byte)(unsafe.Pointer(dst)) = *(*[147]byte)(unsafe.Pointer(src))
}

func copy148(dst, src unsafe.Pointer) {
	*(*[148]byte)(unsafe.Pointer(dst)) = *(*[148]byte)(unsafe.Pointer(src))
}

func copy149(dst, src unsafe.Pointer) {
	*(*[149]byte)(unsafe.Pointer(dst)) = *(*[149]byte)(unsafe.Pointer(src))
}

func copy150(dst, src unsafe.Pointer) {
	*(*[150]byte)(unsafe.Pointer(dst)) = *(*[150]byte)(unsafe.Pointer(src))
}

func copy151(dst, src unsafe.Pointer) {
	*(*[151]byte)(unsafe.Pointer(dst)) = *(*[151]byte)(unsafe.Pointer(src))
}

func copy152(dst, src unsafe.Pointer) {
	*(*[152]byte)(unsafe.Pointer(dst)) = *(*[152]byte)(unsafe.Pointer(src))
}

func copy153(dst, src unsafe.Pointer) {
	*(*[153]byte)(unsafe.Pointer(dst)) = *(*[153]byte)(unsafe.Pointer(src))
}

func copy154(dst, src unsafe.Pointer) {
	*(*[154]byte)(unsafe.Pointer(dst)) = *(*[154]byte)(unsafe.Pointer(src))
}

func copy155(dst, src unsafe.Pointer) {
	*(*[155]byte)(unsafe.Pointer(dst)) = *(*[155]byte)(unsafe.Pointer(src))
}

func copy156(dst, src unsafe.Pointer) {
	*(*[156]byte)(unsafe.Pointer(dst)) = *(*[156]byte)(unsafe.Pointer(src))
}

func copy157(dst, src unsafe.Pointer) {
	*(*[157]byte)(unsafe.Pointer(dst)) = *(*[157]byte)(unsafe.Pointer(src))
}

func copy158(dst, src unsafe.Pointer) {
	*(*[158]byte)(unsafe.Pointer(dst)) = *(*[158]byte)(unsafe.Pointer(src))
}

func copy159(dst, src unsafe.Pointer) {
	*(*[159]byte)(unsafe.Pointer(dst)) = *(*[159]byte)(unsafe.Pointer(src))
}

func copy160(dst, src unsafe.Pointer) {
	*(*[160]byte)(unsafe.Pointer(dst)) = *(*[160]byte)(unsafe.Pointer(src))
}

func copy161(dst, src unsafe.Pointer) {
	*(*[161]byte)(unsafe.Pointer(dst)) = *(*[161]byte)(unsafe.Pointer(src))
}

func copy162(dst, src unsafe.Pointer) {
	*(*[162]byte)(unsafe.Pointer(dst)) = *(*[162]byte)(unsafe.Pointer(src))
}

func copy163(dst, src unsafe.Pointer) {
	*(*[163]byte)(unsafe.Pointer(dst)) = *(*[163]byte)(unsafe.Pointer(src))
}

func copy164(dst, src unsafe.Pointer) {
	*(*[164]byte)(unsafe.Pointer(dst)) = *(*[164]byte)(unsafe.Pointer(src))
}

func copy165(dst, src unsafe.Pointer) {
	*(*[165]byte)(unsafe.Pointer(dst)) = *(*[165]byte)(unsafe.Pointer(src))
}

func copy166(dst, src unsafe.Pointer) {
	*(*[166]byte)(unsafe.Pointer(dst)) = *(*[166]byte)(unsafe.Pointer(src))
}

func copy167(dst, src unsafe.Pointer) {
	*(*[167]byte)(unsafe.Pointer(dst)) = *(*[167]byte)(unsafe.Pointer(src))
}

func copy168(dst, src unsafe.Pointer) {
	*(*[168]byte)(unsafe.Pointer(dst)) = *(*[168]byte)(unsafe.Pointer(src))
}

func copy169(dst, src unsafe.Pointer) {
	*(*[169]byte)(unsafe.Pointer(dst)) = *(*[169]byte)(unsafe.Pointer(src))
}

func copy170(dst, src unsafe.Pointer) {
	*(*[170]byte)(unsafe.Pointer(dst)) = *(*[170]byte)(unsafe.Pointer(src))
}

func copy171(dst, src unsafe.Pointer) {
	*(*[171]byte)(unsafe.Pointer(dst)) = *(*[171]byte)(unsafe.Pointer(src))
}

func copy172(dst, src unsafe.Pointer) {
	*(*[172]byte)(unsafe.Pointer(dst)) = *(*[172]byte)(unsafe.Pointer(src))
}

func copy173(dst, src unsafe.Pointer) {
	*(*[173]byte)(unsafe.Pointer(dst)) = *(*[173]byte)(unsafe.Pointer(src))
}

func copy174(dst, src unsafe.Pointer) {
	*(*[174]byte)(unsafe.Pointer(dst)) = *(*[174]byte)(unsafe.Pointer(src))
}

func copy175(dst, src unsafe.Pointer) {
	*(*[175]byte)(unsafe.Pointer(dst)) = *(*[175]byte)(unsafe.Pointer(src))
}

func copy176(dst, src unsafe.Pointer) {
	*(*[176]byte)(unsafe.Pointer(dst)) = *(*[176]byte)(unsafe.Pointer(src))
}

func copy177(dst, src unsafe.Pointer) {
	*(*[177]byte)(unsafe.Pointer(dst)) = *(*[177]byte)(unsafe.Pointer(src))
}

func copy178(dst, src unsafe.Pointer) {
	*(*[178]byte)(unsafe.Pointer(dst)) = *(*[178]byte)(unsafe.Pointer(src))
}

func copy179(dst, src unsafe.Pointer) {
	*(*[179]byte)(unsafe.Pointer(dst)) = *(*[179]byte)(unsafe.Pointer(src))
}

func copy180(dst, src unsafe.Pointer) {
	*(*[180]byte)(unsafe.Pointer(dst)) = *(*[180]byte)(unsafe.Pointer(src))
}

func copy181(dst, src unsafe.Pointer) {
	*(*[181]byte)(unsafe.Pointer(dst)) = *(*[181]byte)(unsafe.Pointer(src))
}

func copy182(dst, src unsafe.Pointer) {
	*(*[182]byte)(unsafe.Pointer(dst)) = *(*[182]byte)(unsafe.Pointer(src))
}

func copy183(dst, src unsafe.Pointer) {
	*(*[183]byte)(unsafe.Pointer(dst)) = *(*[183]byte)(unsafe.Pointer(src))
}

func copy184(dst, src unsafe.Pointer) {
	*(*[184]byte)(unsafe.Pointer(dst)) = *(*[184]byte)(unsafe.Pointer(src))
}

func copy185(dst, src unsafe.Pointer) {
	*(*[185]byte)(unsafe.Pointer(dst)) = *(*[185]byte)(unsafe.Pointer(src))
}

func copy186(dst, src unsafe.Pointer) {
	*(*[186]byte)(unsafe.Pointer(dst)) = *(*[186]byte)(unsafe.Pointer(src))
}

func copy187(dst, src unsafe.Pointer) {
	*(*[187]byte)(unsafe.Pointer(dst)) = *(*[187]byte)(unsafe.Pointer(src))
}

func copy188(dst, src unsafe.Pointer) {
	*(*[188]byte)(unsafe.Pointer(dst)) = *(*[188]byte)(unsafe.Pointer(src))
}

func copy189(dst, src unsafe.Pointer) {
	*(*[189]byte)(unsafe.Pointer(dst)) = *(*[189]byte)(unsafe.Pointer(src))
}

func copy190(dst, src unsafe.Pointer) {
	*(*[190]byte)(unsafe.Pointer(dst)) = *(*[190]byte)(unsafe.Pointer(src))
}

func copy191(dst, src unsafe.Pointer) {
	*(*[191]byte)(unsafe.Pointer(dst)) = *(*[191]byte)(unsafe.Pointer(src))
}

func copy192(dst, src unsafe.Pointer) {
	*(*[192]byte)(unsafe.Pointer(dst)) = *(*[192]byte)(unsafe.Pointer(src))
}

func copy193(dst, src unsafe.Pointer) {
	*(*[193]byte)(unsafe.Pointer(dst)) = *(*[193]byte)(unsafe.Pointer(src))
}

func copy194(dst, src unsafe.Pointer) {
	*(*[194]byte)(unsafe.Pointer(dst)) = *(*[194]byte)(unsafe.Pointer(src))
}

func copy195(dst, src unsafe.Pointer) {
	*(*[195]byte)(unsafe.Pointer(dst)) = *(*[195]byte)(unsafe.Pointer(src))
}

func copy196(dst, src unsafe.Pointer) {
	*(*[196]byte)(unsafe.Pointer(dst)) = *(*[196]byte)(unsafe.Pointer(src))
}

func copy197(dst, src unsafe.Pointer) {
	*(*[197]byte)(unsafe.Pointer(dst)) = *(*[197]byte)(unsafe.Pointer(src))
}

func copy198(dst, src unsafe.Pointer) {
	*(*[198]byte)(unsafe.Pointer(dst)) = *(*[198]byte)(unsafe.Pointer(src))
}

func copy199(dst, src unsafe.Pointer) {
	*(*[199]byte)(unsafe.Pointer(dst)) = *(*[199]byte)(unsafe.Pointer(src))
}

func copy200(dst, src unsafe.Pointer) {
	*(*[200]byte)(unsafe.Pointer(dst)) = *(*[200]byte)(unsafe.Pointer(src))
}

func copy201(dst, src unsafe.Pointer) {
	*(*[201]byte)(unsafe.Pointer(dst)) = *(*[201]byte)(unsafe.Pointer(src))
}

func copy202(dst, src unsafe.Pointer) {
	*(*[202]byte)(unsafe.Pointer(dst)) = *(*[202]byte)(unsafe.Pointer(src))
}

func copy203(dst, src unsafe.Pointer) {
	*(*[203]byte)(unsafe.Pointer(dst)) = *(*[203]byte)(unsafe.Pointer(src))
}

func copy204(dst, src unsafe.Pointer) {
	*(*[204]byte)(unsafe.Pointer(dst)) = *(*[204]byte)(unsafe.Pointer(src))
}

func copy205(dst, src unsafe.Pointer) {
	*(*[205]byte)(unsafe.Pointer(dst)) = *(*[205]byte)(unsafe.Pointer(src))
}

func copy206(dst, src unsafe.Pointer) {
	*(*[206]byte)(unsafe.Pointer(dst)) = *(*[206]byte)(unsafe.Pointer(src))
}

func copy207(dst, src unsafe.Pointer) {
	*(*[207]byte)(unsafe.Pointer(dst)) = *(*[207]byte)(unsafe.Pointer(src))
}

func copy208(dst, src unsafe.Pointer) {
	*(*[208]byte)(unsafe.Pointer(dst)) = *(*[208]byte)(unsafe.Pointer(src))
}

func copy209(dst, src unsafe.Pointer) {
	*(*[209]byte)(unsafe.Pointer(dst)) = *(*[209]byte)(unsafe.Pointer(src))
}

func copy210(dst, src unsafe.Pointer) {
	*(*[210]byte)(unsafe.Pointer(dst)) = *(*[210]byte)(unsafe.Pointer(src))
}

func copy211(dst, src unsafe.Pointer) {
	*(*[211]byte)(unsafe.Pointer(dst)) = *(*[211]byte)(unsafe.Pointer(src))
}

func copy212(dst, src unsafe.Pointer) {
	*(*[212]byte)(unsafe.Pointer(dst)) = *(*[212]byte)(unsafe.Pointer(src))
}

func copy213(dst, src unsafe.Pointer) {
	*(*[213]byte)(unsafe.Pointer(dst)) = *(*[213]byte)(unsafe.Pointer(src))
}

func copy214(dst, src unsafe.Pointer) {
	*(*[214]byte)(unsafe.Pointer(dst)) = *(*[214]byte)(unsafe.Pointer(src))
}

func copy215(dst, src unsafe.Pointer) {
	*(*[215]byte)(unsafe.Pointer(dst)) = *(*[215]byte)(unsafe.Pointer(src))
}

func copy216(dst, src unsafe.Pointer) {
	*(*[216]byte)(unsafe.Pointer(dst)) = *(*[216]byte)(unsafe.Pointer(src))
}

func copy217(dst, src unsafe.Pointer) {
	*(*[217]byte)(unsafe.Pointer(dst)) = *(*[217]byte)(unsafe.Pointer(src))
}

func copy218(dst, src unsafe.Pointer) {
	*(*[218]byte)(unsafe.Pointer(dst)) = *(*[218]byte)(unsafe.Pointer(src))
}

func copy219(dst, src unsafe.Pointer) {
	*(*[219]byte)(unsafe.Pointer(dst)) = *(*[219]byte)(unsafe.Pointer(src))
}

func copy220(dst, src unsafe.Pointer) {
	*(*[220]byte)(unsafe.Pointer(dst)) = *(*[220]byte)(unsafe.Pointer(src))
}

func copy221(dst, src unsafe.Pointer) {
	*(*[221]byte)(unsafe.Pointer(dst)) = *(*[221]byte)(unsafe.Pointer(src))
}

func copy222(dst, src unsafe.Pointer) {
	*(*[222]byte)(unsafe.Pointer(dst)) = *(*[222]byte)(unsafe.Pointer(src))
}

func copy223(dst, src unsafe.Pointer) {
	*(*[223]byte)(unsafe.Pointer(dst)) = *(*[223]byte)(unsafe.Pointer(src))
}

func copy224(dst, src unsafe.Pointer) {
	*(*[224]byte)(unsafe.Pointer(dst)) = *(*[224]byte)(unsafe.Pointer(src))
}

func copy225(dst, src unsafe.Pointer) {
	*(*[225]byte)(unsafe.Pointer(dst)) = *(*[225]byte)(unsafe.Pointer(src))
}

func copy226(dst, src unsafe.Pointer) {
	*(*[226]byte)(unsafe.Pointer(dst)) = *(*[226]byte)(unsafe.Pointer(src))
}

func copy227(dst, src unsafe.Pointer) {
	*(*[227]byte)(unsafe.Pointer(dst)) = *(*[227]byte)(unsafe.Pointer(src))
}

func copy228(dst, src unsafe.Pointer) {
	*(*[228]byte)(unsafe.Pointer(dst)) = *(*[228]byte)(unsafe.Pointer(src))
}

func copy229(dst, src unsafe.Pointer) {
	*(*[229]byte)(unsafe.Pointer(dst)) = *(*[229]byte)(unsafe.Pointer(src))
}

func copy230(dst, src unsafe.Pointer) {
	*(*[230]byte)(unsafe.Pointer(dst)) = *(*[230]byte)(unsafe.Pointer(src))
}

func copy231(dst, src unsafe.Pointer) {
	*(*[231]byte)(unsafe.Pointer(dst)) = *(*[231]byte)(unsafe.Pointer(src))
}

func copy232(dst, src unsafe.Pointer) {
	*(*[232]byte)(unsafe.Pointer(dst)) = *(*[232]byte)(unsafe.Pointer(src))
}

func copy233(dst, src unsafe.Pointer) {
	*(*[233]byte)(unsafe.Pointer(dst)) = *(*[233]byte)(unsafe.Pointer(src))
}

func copy234(dst, src unsafe.Pointer) {
	*(*[234]byte)(unsafe.Pointer(dst)) = *(*[234]byte)(unsafe.Pointer(src))
}

func copy235(dst, src unsafe.Pointer) {
	*(*[235]byte)(unsafe.Pointer(dst)) = *(*[235]byte)(unsafe.Pointer(src))
}

func copy236(dst, src unsafe.Pointer) {
	*(*[236]byte)(unsafe.Pointer(dst)) = *(*[236]byte)(unsafe.Pointer(src))
}

func copy237(dst, src unsafe.Pointer) {
	*(*[237]byte)(unsafe.Pointer(dst)) = *(*[237]byte)(unsafe.Pointer(src))
}

func copy238(dst, src unsafe.Pointer) {
	*(*[238]byte)(unsafe.Pointer(dst)) = *(*[238]byte)(unsafe.Pointer(src))
}

func copy239(dst, src unsafe.Pointer) {
	*(*[239]byte)(unsafe.Pointer(dst)) = *(*[239]byte)(unsafe.Pointer(src))
}

func copy240(dst, src unsafe.Pointer) {
	*(*[240]byte)(unsafe.Pointer(dst)) = *(*[240]byte)(unsafe.Pointer(src))
}

func copy241(dst, src unsafe.Pointer) {
	*(*[241]byte)(unsafe.Pointer(dst)) = *(*[241]byte)(unsafe.Pointer(src))
}

func copy242(dst, src unsafe.Pointer) {
	*(*[242]byte)(unsafe.Pointer(dst)) = *(*[242]byte)(unsafe.Pointer(src))
}

func copy243(dst, src unsafe.Pointer) {
	*(*[243]byte)(unsafe.Pointer(dst)) = *(*[243]byte)(unsafe.Pointer(src))
}

func copy244(dst, src unsafe.Pointer) {
	*(*[244]byte)(unsafe.Pointer(dst)) = *(*[244]byte)(unsafe.Pointer(src))
}

func copy245(dst, src unsafe.Pointer) {
	*(*[245]byte)(unsafe.Pointer(dst)) = *(*[245]byte)(unsafe.Pointer(src))
}

func copy246(dst, src unsafe.Pointer) {
	*(*[246]byte)(unsafe.Pointer(dst)) = *(*[246]byte)(unsafe.Pointer(src))
}

func copy247(dst, src unsafe.Pointer) {
	*(*[247]byte)(unsafe.Pointer(dst)) = *(*[247]byte)(unsafe.Pointer(src))
}

func copy248(dst, src unsafe.Pointer) {
	*(*[248]byte)(unsafe.Pointer(dst)) = *(*[248]byte)(unsafe.Pointer(src))
}

func copy249(dst, src unsafe.Pointer) {
	*(*[249]byte)(unsafe.Pointer(dst)) = *(*[249]byte)(unsafe.Pointer(src))
}

func copy250(dst, src unsafe.Pointer) {
	*(*[250]byte)(unsafe.Pointer(dst)) = *(*[250]byte)(unsafe.Pointer(src))
}

func copy251(dst, src unsafe.Pointer) {
	*(*[251]byte)(unsafe.Pointer(dst)) = *(*[251]byte)(unsafe.Pointer(src))
}

func copy252(dst, src unsafe.Pointer) {
	*(*[252]byte)(unsafe.Pointer(dst)) = *(*[252]byte)(unsafe.Pointer(src))
}

func copy253(dst, src unsafe.Pointer) {
	*(*[253]byte)(unsafe.Pointer(dst)) = *(*[253]byte)(unsafe.Pointer(src))
}

func copy254(dst, src unsafe.Pointer) {
	*(*[254]byte)(unsafe.Pointer(dst)) = *(*[254]byte)(unsafe.Pointer(src))
}

func copy255(dst, src unsafe.Pointer) {
	*(*[255]byte)(unsafe.Pointer(dst)) = *(*[255]byte)(unsafe.Pointer(src))
}
