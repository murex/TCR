import {OnOffPipe} from './on-off.pipe';

describe('OnOffPipe', () => {
  let pipe: OnOffPipe;

  beforeEach(() => {
    pipe = new OnOffPipe();
  });

  const markerOn = "✅";
  const markerOff = "❌";

  [
    {input: true, expected: markerOn},
    {input: false, expected: markerOff},
    {input: "non-empty string", expected: markerOn},
    {input: "", expected: markerOff},
    {input: null, expected: markerOff},
    {input: undefined, expected: markerOff}
  ].forEach(testCase => {
    it(`should return "${testCase.expected}" when the input value is ${testCase.input}`, () => {
      expect(pipe.transform(testCase.input)).toEqual(testCase.expected);
    });
  });
});
