import Test.Hspec        (Spec, it, shouldBe)
import Test.Hspec.Runner (configFailFast, defaultConfig, hspecWith)

import HelloWorld (sayHello)

main :: IO ()
main = hspecWith defaultConfig {configFailFast = True} specs

specs :: Spec
specs = it "say hello" $
          sayHello "Sue" `shouldBe` "Hello, Sue!"
