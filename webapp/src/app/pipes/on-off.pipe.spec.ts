import {OnOffPipe} from './on-off.pipe';

const MARKER_ON = "✅";
const MARKER_OFF = "❌";

describe('OnOffPipe', () => {
  it('create an instance', () => {
    expect(new OnOffPipe()).toBeTruthy();
  });

  it(`should return "${MARKER_ON}" when the input value is truthy`, () => {
    expect(new OnOffPipe().transform(true)).toEqual(MARKER_ON);
  });

  it(`should return "${MARKER_OFF}" when the input value is falsy`, () => {
    expect(new OnOffPipe().transform(false)).toEqual(MARKER_OFF);
  });

  it(`should return "${MARKER_ON}" when the input value is a non-empty string`, () => {
    expect(new OnOffPipe().transform("non-empty string")).toEqual(MARKER_ON);
  });

  it(`should return "${MARKER_OFF}" when the input value is an empty string`, () => {
    expect(new OnOffPipe().transform("")).toEqual(MARKER_OFF);
  });

  it(`should return "${MARKER_OFF}" when the input value is null`, () => {
    expect(new OnOffPipe().transform(null)).toEqual(MARKER_OFF);
  });

  it(`should return "${MARKER_OFF}" when the input value is undefined`, () => {
    expect(new OnOffPipe().transform(undefined)).toEqual(MARKER_OFF);
  });
});
