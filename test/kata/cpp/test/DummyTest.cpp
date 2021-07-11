
#include <kata/dummy/Dummy.hpp>

#include <gtest/gtest.h>

TEST(Dummy, acceptance_test)
{
	EXPECT_EQ(42, kata::dummy::doSomething());
}
