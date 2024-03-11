import {ShowEmptyPipe} from './show-empty.pipe';

const NOT_SET = '[not set]';

describe('ShowEmptyPipe', () => {
  it('create an instance', () => {
    expect(new ShowEmptyPipe()).toBeTruthy();
  });

  it('should return the input value when it is not empty', () => {
    const input = 'input value';
    expect(new ShowEmptyPipe().transform(input)).toEqual(input);
  });

  it('should return "[not set]" when the input value is null', () => {
    expect(new ShowEmptyPipe().transform(null)).toEqual(NOT_SET);
  });

  it('should return "[not set]" when the input value is undefined', () => {
    expect(new ShowEmptyPipe().transform(undefined)).toEqual(NOT_SET);
  });

  it('should return "[not set]" when the input value is an empty string', () => {
    expect(new ShowEmptyPipe().transform('')).toEqual(NOT_SET);
  });
});
