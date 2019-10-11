# Go-Rust Demo

This is a simple library showing my learnings in how to combine Go and Rust.
It assumes there is an existing rust library with the desired functionality,
and the purpose is to produce a nice Go interface for this.

### Structure

The approach taken is to create one join Go and Rust repo that can then be
imported by the Go project. On the Rust side, it will produce `extern "C"` bindings
to the library and compile down to a `cdylib`, that is, you get a `*.so, *.dll, etc`
that exposes a C API to the Rust code. On the Go side, you use `cgo` to create
a Go bridge to this code. The full build step
involves compiling rust -> C library, and linking that library to the Go code.
For ergonomics of the user, we will include pre-compiled libraries to easily
link with, and Go developers should just be able to import this directly.

## Gotchas

Beyond learning the intracacies of both cgo and rust ffi, there are two points
that required a bit more research. So it is nice to look into how they are solved here.
I don't pretend that these are perfect solutions, but they cover the cases that
matter to me. I would love feedback on how to improve as well.

### Passing Strings / Bytes

When there is a large chunk of memory that cannot fit on the stack, passing
over ffi provides a challenge. The solution in C is usually to pass a desciptor
to the memory `ptr: *mut u8, len: usize` and then reconstruct a slice inside
the rust library. Passing out memory involves allocating some memory, creating
a reference to it and "forgetting" it for Rust's deallocator. Hoping that the
caller will later make another call to the library to free it.

We define some helpers in `memory.rs` and use a `Buffer` type for passing the pointer
and size together as an argument or a return value.

### Errors and Panics

We happily use `Result<T, E>` all over the rust code, but that cannot cross the FFI
boundary. However, C does provide us access to `errno` to signal a numerical error
code, but not a message. We add some logic to set `errno` upon an error and
store the error message in a singleton, which can be queried from the called
by `get_last_error()`. `cgo` provides support in that it will set the `_, err`
return value if `errno` is set. We can detect this, and then load the custom
message if we detect something went wrong.

An uncaught panic that hits the FFI boundary may also crash the calling process,
which is very bad behavior. To avoid this, we must `catch_unwind()` on any code
that may panic, and then convert it to a normal error that we can report as above.
Note, that this also requires `panic = "unwind"` to be set in `profile.release` in
`Cargo.toml`.

You can look at the helpers in `errors.go` and an example using this (with panics)
in `lib.rs:may_panic`.

## References

This example arose from learnings when producing go bindings to the cosmwasm smart
contract framework. You can look at [go-cosmwasm](https://github.com/confio/go-cosmwasm)
to see a full API build using these techniques. I will try to port any future
learnings to this demo module as well.