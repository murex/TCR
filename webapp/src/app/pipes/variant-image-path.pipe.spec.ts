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

import {VariantImagePathPipe} from './variant-image-path.pipe';

describe('VariantImagePathPipe', () => {
  let pipe: VariantImagePathPipe;

  beforeEach(() => {
    pipe = new VariantImagePathPipe();
  });

  [
    {input: 'relaxed', expected: 'assets/images/variant-relaxed.png'},
    {input: 'btcr', expected: 'assets/images/variant-btcr.png'},
    {input: 'original', expected: 'assets/images/variant-original.png'},
    {input: 'introspective', expected: 'assets/images/variant-introspective.png'},
    {input: null, expected: ''},
    {input: undefined, expected: ''},
    {input: '', expected: ''}
  ].forEach(testCase => {
    it(`should return "${testCase.expected}" when the variant value is "${testCase.input}"`, () => {
      expect(pipe.transform(testCase.input)).toEqual(testCase.expected);
    });
  });
});
