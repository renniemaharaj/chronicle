import { afterEach, describe, expect, it, vi } from 'vitest';
import { motionLoop } from './story-motion-loop';

afterEach(() => {
  vi.restoreAllMocks();
  vi.unstubAllGlobals();
});

describe('smooth presentation clock', () => {
  it('draws between narration samples and cancels immediately when paused', () => {
    let callback: FrameRequestCallback | undefined;
    const cancel = vi.fn(() => {
      callback = undefined;
    });
    vi.stubGlobal(
      'requestAnimationFrame',
      vi.fn((next) => {
        callback = next;

        return 1;
      }),
    );
    vi.stubGlobal('cancelAnimationFrame', cancel);
    vi.spyOn(performance, 'now').mockReturnValue(0);

    const frames: number[] = [];
    const loop = motionLoop((seconds) => frames.push(seconds));
    const effects = {
      pulse: true,
      camera: true,
    };
    loop.update(10, 1, true, effects);
    callback!(16);
    callback!(32);
    callback!(48);
    expect(frames).toEqual([10.016, 10.032, 10.048]);
    loop.update(10.05, 1, false, effects);
    expect(callback).toBeUndefined();
    expect(cancel).toHaveBeenCalled();
    expect(frames.at(-1)).toBe(10.05);
  });

  it('honors playback speed and stops extrapolating during a stalled audio clock', () => {
    let callback: FrameRequestCallback | undefined;
    vi.stubGlobal(
      'requestAnimationFrame',
      vi.fn((next) => {
        callback = next;

        return 1;
      }),
    );
    vi.stubGlobal('cancelAnimationFrame', vi.fn());
    vi.spyOn(performance, 'now').mockReturnValue(0);

    const frames: number[] = [];
    const loop = motionLoop((seconds) => frames.push(seconds));
    loop.update(10, 2, true, {
      pulse: true,
      camera: false,
    });
    callback!(25);
    callback!(1000);
    expect(frames).toEqual([10.05, 10.2]);
    loop.stop();
  });
});
