mod error;
mod memory;

pub use memory::{free_rust, Buffer};

use crate::error::{handle_c_error, set_error};
use std::panic::catch_unwind;

#[no_mangle]
pub extern "C" fn add(a: i32, b: i32) -> i32 {
    a + b
}

#[no_mangle]
pub extern "C" fn greet(name: Buffer) -> Buffer {
    let rname = name.read().unwrap_or(b"<nil>");
    let mut v = b"Hello, ".to_vec();
    v.extend_from_slice(rname);
    Buffer::from_vec(v)
}

/// divide returns the rounded (i32) result, returns a C error if div == 0
#[no_mangle]
pub extern "C" fn divide(num: i32, div: i32, err: Option<&mut Buffer>) -> i32 {
    if div == 0 {
        set_error("Cannot divide by zero".to_string(), err);
        return 0;
    }
    num / div
}

#[no_mangle]
pub extern "C" fn may_panic(guess: i32, err: Option<&mut Buffer>) -> Buffer {
    let r = catch_unwind(|| do_may_panic(guess)).unwrap_or(Err("Caught panic".to_string()));
    let v = handle_c_error(r, err).into_bytes();
    Buffer::from_vec(v)
}

fn do_may_panic(guess: i32) -> Result<String, String> {
    if guess == 0 {
        panic!("Must be negative or positive")
    } else if guess < 17 {
        Err("Too low".to_owned())
    } else {
        Ok("You are a winner!".to_owned())
    }
}

// this represents something passed in from the caller side of FFI
#[repr(C)]
pub struct db_t { }

#[repr(C)]
pub struct DB {
    pub state: *mut db_t,
    pub c_get: extern fn(*mut db_t, Buffer, Buffer) -> i64,
    pub c_set: extern fn(*mut db_t, Buffer, Buffer),
}

impl DB {
    pub fn get(&self, key: Vec<u8>) -> Option<Vec<u8>> {
        let buf = Buffer::from_vec(key);
        // TODO: dynamic size
        let mut buf2 = Buffer::from_vec(vec![0u8; 2000]);
        let res = (self.c_get)(self.state, buf, buf2);

        // read in the number of bytes returned
        if res < 0 {
            // TODO
            panic!("val was not big enough for data");
        }
        if res == 0 {
            return None
        }
        buf2.len = res as usize;
        unsafe { Some(buf2.consume()) }
    }

    pub fn set(&self, key: Vec<u8>, value: Vec<u8>) {
        let buf = Buffer::from_vec(key);
        let buf2 = Buffer::from_vec(value);
        // caller will free input
        (self.c_set)(self.state, buf, buf2);
    }
}

// This loads key from DB and then
#[no_mangle]
pub extern "C" fn db_access(db: DB, key: Buffer) -> Buffer {
    let vkey = key.read().unwrap_or(b"no-key").to_vec();
    let val = db.get(vkey.clone()).unwrap_or(b"<nil>".to_vec());
    let mut res = b"Got value: ".to_vec();
    res.extend_from_slice(&val);
    Buffer::from_vec(res)
}
