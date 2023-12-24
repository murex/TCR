import { ShowEmptyPipe } from './show-empty.pipe';

describe('ShowEmptyPipe', () => {
  it('create an instance', () => {
    const pipe = new ShowEmptyPipe();
    expect(pipe).toBeTruthy();
  });
});
