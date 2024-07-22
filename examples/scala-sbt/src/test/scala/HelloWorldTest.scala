import org.scalatest.funsuite.AnyFunSuite
import org.scalatest.matchers.should.Matchers


/** @version 1.1.0 */
class HelloWorldTest extends AnyFunSuite with Matchers {

  test("Say Hello") {
    HelloWorld.sayHello("Sue") should be ("Hello, Sue!")
  }
}
