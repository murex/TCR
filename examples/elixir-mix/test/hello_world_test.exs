defmodule HelloWorldTest do
  use ExUnit.Case

  test "say hello" do
    assert HelloWorld.say_hello("Sue") == "Hello Sue!"
  end

end
