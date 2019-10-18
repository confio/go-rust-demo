mod db;
mod error;
mod memory;

pub use memory::{free_rust, Buffer};
pub use crate::db::{db_t, DB};

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

// This loads key from DB and then
#[no_mangle]
pub extern "C" fn db_access(db: DB, key: Buffer, err: Option<&mut Buffer>) -> Buffer {
    let r = catch_unwind(|| do_db_access(db, key)).unwrap_or(Err("Caught panic".to_string()));
    let v = handle_c_error(r, err);
    Buffer::from_vec(v)
}

fn do_db_access(db: DB, key: Buffer) -> Result<Vec<u8>, String> {
    let vkey = key.read().ok_or("no input".to_string())?.to_vec();
    let val = db.get(vkey.clone()).ok_or("no data".to_string())?;
    let mut res = b"Got value: ".to_vec();
    res.extend_from_slice(&val);
    Ok(res)
}

