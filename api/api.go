package api

// #cgo LDFLAGS: -Wl,-rpath,${SRCDIR} -L${SRCDIR} -lgo_rust_demo
// #include <stdlib.h>
// #include "bindings.h"
import "C"

import "fmt"
import "unsafe"

// nice aliases to the rust names
type i32 = C.int32_t
type u8 = C.uint8_t
type u8_ptr = *C.uint8_t
type usize = C.uintptr_t
type cint = C.int


func Add(a int32, b int32) int32 {
	return (int32)(C.add(i32(a), i32(b)))
}

func Greet(name []byte) []byte {
	buf := sendSlice(name)
	raw := C.greet(buf)
	// make sure to free after call
	freeAfterSend(buf)

	return receiveSlice(raw)
}

func Divide(a, b int32) (int32, error) {
    buf := C.Buffer{}
	res, err := C.divide(i32(a), i32(b), &buf)
	if err != nil {
		return 0, errorWithMessage(err, buf)
	}
	return int32(res), nil
}

func RandomMessage(guess int32) (string, error) {
    buf := C.Buffer{}
	res, err := C.may_panic(i32(guess), &buf)
	if err != nil {
		return "", errorWithMessage(err, buf)
	}
	return string(receiveSlice(res)), nil
}


/**** To error module ***/

func errorWithMessage(err error, b C.Buffer) error {
	msg := receiveSlice(b)
	if msg == nil {
		return err
	}
	return fmt.Errorf("%s", string(msg))
}

/*** To memory module **/

func sendSlice(s []byte) C.Buffer {
	if s == nil {
		return C.Buffer{ptr: u8_ptr(nil), len: usize(0), cap: usize(0)};
	}
	return C.Buffer{
		ptr: u8_ptr(C.CBytes(s)),
		len: usize(len(s)),
		cap: usize(len(s)),
	}
}

func receiveSlice(b C.Buffer) []byte {
	if emptyBuf(b) {
		return nil
	}
	res := C.GoBytes(unsafe.Pointer(b.ptr), cint(b.len))
	C.free_rust(b)
	return res
}

func freeAfterSend(b C.Buffer) {
	if !emptyBuf(b) {
		C.free(unsafe.Pointer(b.ptr))
	}
}

func emptyBuf(b C.Buffer) bool {
	return b.ptr == u8_ptr(nil) || b.len == usize(0) || b.cap == usize(0)
}


