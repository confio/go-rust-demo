package api

/*
#cgo CFLAGS: -I .
#cgo LDFLAGS: -Wl,-rpath,${SRCDIR} -L${SRCDIR} -lgo_rust_demo

#include "bindings.h"
*/
import "C"

import (
	"fmt"
)

// nice aliases to the rust names
type i32 = C.int32_t
type i64 = C.int64_t
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

func DemoDBAccess(kv KVStore, key []byte) ([]byte, error) {
	db := buildDB(kv)
	buf := sendSlice(key)
	msg := C.Buffer{}
	res, err := C.db_access(db, buf, &msg)
	freeAfterSend(buf)
	if err != nil {
		return nil, errorWithMessage(err, msg)
	}
	return receiveSlice(res), nil
}

/**** To error module ***/

func errorWithMessage(err error, b C.Buffer) error {
	msg := receiveSlice(b)
	if msg == nil {
		return err
	}
	return fmt.Errorf("%s", string(msg))
}
