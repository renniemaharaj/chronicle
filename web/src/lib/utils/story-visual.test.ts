import { describe, expect, it } from 'vitest';
import { imageURL, validVisual, visualOpacity } from './story-visual';
import type { MediaAsset } from '../types/media';

const cue = {
  mediaId: 'synthetic',
  startSeconds: 4,
  durationSeconds: 8,
  transitionSeconds: 1,
};

describe('passage image timing', () => {
  it('returns to the map outside the cue and fades at both boundaries', () => {
    expect(visualOpacity(cue, 3)).toBe(0);
    expect(visualOpacity(cue, 4)).toBe(0);
    expect(visualOpacity(cue, 4.5)).toBe(0.5);
    expect(visualOpacity(cue, 7)).toBe(1);
    expect(visualOpacity(cue, 11.5)).toBe(0.5);
    expect(visualOpacity(cue, 12)).toBe(0);
    expect(visualOpacity(cue, 80)).toBe(0);
  });

  it('rejects impossible timings and keeps instant transitions bounded', () => {
    expect(
      validVisual({
        ...cue,
        durationSeconds: 1,
      }),
    ).toBe(false);
    expect(
      validVisual({
        ...cue,
        startSeconds: NaN,
      }),
    ).toBe(false);
    expect(
      validVisual({
        ...cue,
        transitionSeconds: -1,
      }),
    ).toBe(false);
    expect(
      visualOpacity(
        {
          ...cue,
          transitionSeconds: 0,
        },
        4,
      ),
    ).toBe(1);
    expect(
      visualOpacity(
        {
          ...cue,
          transitionSeconds: 0,
        },
        12,
      ),
    ).toBe(0);
  });

  it('allows only imported local generated image paths', () => {
    const generated = {
      generated: {},
      url: `/api/media/images/${'a'.repeat(64)}.png`,
    } as MediaAsset;
    expect(imageURL(generated)).toBe(generated.url);
    expect(
      imageURL({
        ...generated,
        url: '/api/config',
      }),
    ).toBeUndefined();
    expect(
      imageURL({
        ...generated,
        url: 'https://example.org/fake.png',
      }),
    ).toBeUndefined();
    expect(
      imageURL({
        ...generated,
        url: '/api/media/images/../../config',
      }),
    ).toBeUndefined();
  });
});
