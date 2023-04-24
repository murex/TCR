import kotlin.test.Test
import kotlin.test.assertEquals

class HelloWorldTest {
    @Test
    fun sayHello() {
        val helloWorld = HelloWorld()
        assertEquals("Hello Sue!", helloWorld.sayHello("Sue"))
    }
}
