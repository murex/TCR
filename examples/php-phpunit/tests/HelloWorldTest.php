<?php

declare(strict_types=1);

namespace HelloWorld\Tests;

use HelloWorld\HelloWorld;
use PHPUnit\Framework\TestCase;

final class HelloWorldTest extends TestCase
{
    public function test_say_hello(): void
    {
        $helloWorld = new HelloWorld;
        $this->assertEquals("Hello Sue!", $helloWorld->sayHello("Sue"));
    }
}
