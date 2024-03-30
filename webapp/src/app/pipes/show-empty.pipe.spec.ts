import {ShowEmptyPipe} from './show-empty.pipe';

describe('ShowEmptyPipe', () => {
  let pipe: ShowEmptyPipe;

  beforeEach(() => {
    pipe = new ShowEmptyPipe();
  });

  const notSet = '[not set]';

  [
    {input: 'input value', expected: 'input value'},
    {input: null, expected: notSet},
    {input: undefined, expected: notSet},
    {input: '', expected: notSet}
  ].forEach(testCase => {
    it(`should return "${testCase.expected}" when the input value is ${testCase.input}`, () => {
      expect(pipe.transform(testCase.input)).toEqual(testCase.expected);
    });
  });
});
