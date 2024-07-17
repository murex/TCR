/*
Copyright (c) 2024 Murex

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

import {FormatTimerPipe} from './format-timer.pipe';

describe('FormatTimerPipe', () => {
  let pipe: FormatTimerPipe;

  beforeEach(() => {
    pipe = new FormatTimerPipe();
  });

  const defaultFormat = '--:--';

  [
    {input: 'abc', expected: defaultFormat},
    {input: '', expected: defaultFormat},
    {input: undefined, expected: defaultFormat},
    {input: null, expected: defaultFormat},
    {input: '0', expected: '00:00'},
    {input: '1', expected: '00:01'},
    {input: '60', expected: '01:00'},
    {input: '119', expected: '01:59'},
    {input: '3600', expected: '1:00:00'},
    {input: '3661', expected: '1:01:01'},
    {input: '-1', expected: '-00:01'},
    {input: '-60', expected: '-01:00'},
    {input: '-119', expected: '-01:59'},
    {input: '-3600', expected: '-1:00:00'},
  ].forEach(testCase => {
    it(`should return "${testCase.expected}" when input is ${testCase.input}`, () => {
      expect(pipe.transform(testCase.input)).toEqual(testCase.expected);
    });
  });
});
