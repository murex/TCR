import { sayHello } from "../src/HelloWorld";

describe("Hello World", () => {
    it("should say hello to Joe", () => {
        expect(sayHello("Joe")).toBe("Hello Joe!");
    });
});
