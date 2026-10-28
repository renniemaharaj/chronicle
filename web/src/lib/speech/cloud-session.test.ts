import { afterEach, describe, expect, it, vi } from 'vitest';
import { googleSession } from './google-session';

class TestAudio {
  static current: TestAudio;
  src = '';
  playbackRate = 1;
  onended = () => {};
  onerror = () => {};
  play = vi.fn().mockResolvedValue(undefined);
  pause = vi.fn();
  removeAttribute = vi.fn();
  constructor() {
    TestAudio.current = this;
  }
}

afterEach(() => {
  vi.unstubAllGlobals();
  vi.restoreAllMocks();
});

describe('cancelled saved audio', () => {
  it('ignores pending audio when the user skips the passage', async () => {
    vi.stubGlobal('Audio', TestAudio);

    let finish!: (value: unknown) => void;
    vi.stubGlobal(
      'fetch',
      vi.fn().mockReturnValue(
        new Promise((resolve) => {
          finish = resolve;
        }),
      ),
    );

    const done = vi.fn();
    const fail = vi.fn();
    const session = googleSession({
      text: 'Test.',
      voice: 'test',
      language: 'en',
      speed: 1,
      done,
      fail,
    });
    session.cancel();
    finish({
      ok: true,
      json: async () => ({
        available: true,
        url: '/saved',
      }),
    });
    await new Promise((resolve) => setTimeout(resolve, 0));
    expect(TestAudio.current.src).toBe('');
    expect(TestAudio.current.play).not.toHaveBeenCalled();
    expect(done).not.toHaveBeenCalled();
    expect(fail).not.toHaveBeenCalled();
  });
});
