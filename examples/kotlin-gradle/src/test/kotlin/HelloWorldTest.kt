import org.junit.Test
import org.junit.Assert.assertEquals

class HelloWorldTest {
    @Test
    fun sayHello() {
        val helloWorld = HelloWorld()
        assertEquals("Hello Sue!", helloWorld.sayHello("Sue"))
    }
}
