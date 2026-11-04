import { afterEach, describe, expect, it, vi } from 'vitest';
import { googleSession } from './google-session';

class SyntheticAudio {
  static current: SyntheticAudio;
  currentTime = 0;
  duration = 60;
  src = '';
  playbackRate = 1;
  onended = () => {};
  onerror = () => {};
  ontimeupdate = () => {};
  play = vi.fn().mockResolvedValue(undefined);
  pause = vi.fn();
  removeAttribute = vi.fn();

  constructor() {
    SyntheticAudio.current = this;
  }
}

afterEach(() => {
  vi.useRealTimers();
  vi.unstubAllGlobals();
});

describe('saved-audio image synchronization', () => {
  it('counts actual media time, ignores loading time, and stops reporting after cancellation', async () => {
    vi.useFakeTimers();
    vi.stubGlobal('Audio', SyntheticAudio);

    let deliver!: (value: unknown) => void;
    vi.stubGlobal(
      'fetch',
      vi.fn().mockReturnValue(
        new Promise((resolve) => {
          deliver = resolve;
        }),
      ),
    );

    const progress = vi.fn();
    const session = googleSession({
      text: 'Synthetic audio clock fixture.',
      voice: 'test',
      language: 'en',
      speed: 1,
      progress,
      done: vi.fn(),
      fail: vi.fn(),
    });
    await vi.advanceTimersByTimeAsync(5000);
    expect(progress).not.toHaveBeenCalled();
    session.pause();
    deliver({
      ok: true,
      json: async () => ({
        available: true,
        url: '/synthetic-audio',
      }),
    });
    await vi.advanceTimersByTimeAsync(1000);
    expect(SyntheticAudio.current.play).not.toHaveBeenCalled();
    expect(progress).not.toHaveBeenCalled();
    SyntheticAudio.current.currentTime = 6;
    session.resume();
    await vi.advanceTimersByTimeAsync(100);
    expect(progress.mock.lastCall?.[0]).toBe(6);
    session.pause();

    const calls = progress.mock.calls.length;
    await vi.advanceTimersByTimeAsync(3000);
    expect(progress).toHaveBeenCalledTimes(calls);
    session.resume();
    await vi.advanceTimersByTimeAsync(100);
    session.cancel();

    const finalCalls = progress.mock.calls.length;
    await vi.advanceTimersByTimeAsync(30000);
    expect(progress).toHaveBeenCalledTimes(finalCalls);
  });
});
