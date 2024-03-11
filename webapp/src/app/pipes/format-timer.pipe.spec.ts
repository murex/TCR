import {FormatTimerPipe} from './format-timer.pipe';

describe('FormatTimerPipe', () => {
  let pipe: FormatTimerPipe;

  beforeEach(() => {
    pipe = new FormatTimerPipe();
  });

  const defaultFormat = '--:--';
  const testCases = [
    {input: 'abc', expected: defaultFormat},
    {input: '', expected: defaultFormat},
    {input: '0', expected: '00:00'},
    {input: '1', expected: '00:01'},
    {input: '60', expected: '01:00'},
    {input: '119', expected: '01:59'},
    {input: '-1', expected: '-00:01'},
    {input: '-60', expected: '-01:00'},
    {input: '-119', expected: '-01:59'}
  ];

  for (const testCase of testCases) {
    it(`should return "${testCase.expected}" when input is ${testCase.input}`, () => {
      expect(pipe.transform(testCase.input)).toEqual(testCase.expected);
    });
  }
});
