import { afterEach, describe, expect, it, vi } from 'vitest';
import { passageClock } from './passage-clock';
import { readingSession } from './session';

afterEach(() => vi.useRealTimers());

describe('passage presentation clock', () => {
  it('freezes while paused and accounts for narration speed', () => {
    vi.useFakeTimers();

    const progress = vi.fn();
    const clock = passageClock({
      text: 'Synthetic timing fixture.',
      voice: '',
      language: 'en',
      speed: 2,
      progress,
      done: vi.fn(),
      fail: vi.fn(),
    });
    clock.resume();
    vi.advanceTimersByTime(1500);
    clock.pause();
    expect(progress.mock.lastCall?.[0]).toBeCloseTo(3);

    const calls = progress.mock.calls.length;
    vi.advanceTimersByTime(10000);
    expect(progress).toHaveBeenCalledTimes(calls);
    clock.resume();
    vi.advanceTimersByTime(500);
    clock.pause();
    expect(progress.mock.lastCall?.[0]).toBeCloseTo(4);
  });

  it('cancelling a screen passage prevents all later visual progress and completion', () => {
    vi.useFakeTimers();

    const progress = vi.fn();
    const done = vi.fn();
    const session = readingSession({
      text: 'Synthetic reading fixture.',
      voice: '',
      language: 'en',
      speed: 1,
      progress,
      done,
      fail: vi.fn(),
    });
    vi.advanceTimersByTime(500);
    session.cancel();

    const calls = progress.mock.calls.length;
    vi.advanceTimersByTime(30000);
    expect(progress).toHaveBeenCalledTimes(calls);
    expect(done).not.toHaveBeenCalled();
  });
});
