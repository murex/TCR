from hello_world.hello_world import say_hello


class TestHelloWorld:

    def test_say_hello(self) -> None:
        assert 'Hello Sue!' == say_hello('Sue')
