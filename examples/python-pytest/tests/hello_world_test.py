import unittest

from hello_world.hello_world import say_hello


class HelloWorldTest(unittest.TestCase):

    def test_say_hello(self):
        self.assertEqual('Hello Sue!', say_hello('Sue'))


if __name__ == "__main__":
    unittest.main()
