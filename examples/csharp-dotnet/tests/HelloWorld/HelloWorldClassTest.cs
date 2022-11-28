using Xunit;

namespace HelloWorld;

public class HelloWorldClassTest
{
    [Fact]
    public void SayHello()
    {
        Assert.Equal("Hello Sue!", HelloWorldClass.SayHello("Sue"));
    }
}
