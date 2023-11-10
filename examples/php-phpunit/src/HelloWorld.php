<?php

declare(strict_types=1);

namespace HelloWorld;

class HelloWorld
{
    public function sayHello(String $name): String
    {
        return sprintf("Hello %s!", $name);
    }
}
