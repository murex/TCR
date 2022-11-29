import org.junit.Test;
import static org.junit.Assert.assertEquals;

public class DummyTest {
    @Test
    public void acceptance_test() {
        assertEquals(42, Dummy.doSomething());
    }
}
