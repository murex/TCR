use hello_world;

#[test]
fn say_hello_test() {
    assert_eq!("Hello Sue!", hello_world::say_hello("Sue"));
}

