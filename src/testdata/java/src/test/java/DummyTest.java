import static org.junit.jupiter.api.Assertions.assertEquals;
import org.junit.jupiter.api.Test;

public class DummyTest {
    @Test
    public void acceptance_test() {
        assertEquals(42, Dummy.doSomething());
    }
}
