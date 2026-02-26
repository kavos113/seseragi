use std::{thread, time};

fn main() {
    println!("Hello, this is rust!");
    thread::sleep(time::Duration::from_secs(5));
    println!("Task Completed!")
}
