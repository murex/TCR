import org.junit.Test;

import static org.junit.Assert.assertEquals;

public class HelloWorldTest {

    @Test
    public void sayHello() {
        assertEquals("Hello Sue!", HelloWorld.sayHello("Sue"));
    }
}
