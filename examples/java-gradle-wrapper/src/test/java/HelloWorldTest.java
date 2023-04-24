import static org.junit.jupiter.api.Assertions.assertEquals;
import org.junit.jupiter.api.Test;

public class HelloWorldTest {

    @Test
    public void sayHello() {
        assertEquals("Hello Sue!", HelloWorld.sayHello("Sue"));
    }
}
