'use strict';

const sayHello = require("../src/helloWorld");

describe("Hello World", () => {
    it("should say hello to Joe", () => {
        expect(sayHello("Joe")).toBe("Hello Joe!");
    });
});
