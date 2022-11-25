#include <HelloWorld.hpp>

#include <gtest/gtest.h>

TEST(HelloWorld, sayHello)
{
    HelloWorld helloWorld;
	 EXPECT_EQ("Hello Sue!", helloWorld.sayHello("Sue"));
}
