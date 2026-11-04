import { describe, expect, it } from 'vitest';
import { cameraMotion } from '../components/map/story-motion';
import { activeStoryCaption, storyCaptions } from './story-captions';
import { floatingImageSize } from './story-image-fit';

describe('story production geometry and timing', () => {
  it('keeps camera drift bounded around one geographic anchor across long passages', () => {
    expect(cameraMotion(0)).toEqual({
      x: 0,
      y: 0,
      zoom: 0,
    });
    expect(cameraMotion(NaN)).toEqual(cameraMotion(0));
    for (let second = 0; second < 3600; second += 7) {
      const position = cameraMotion(second);
      expect(Math.abs(position.x)).toBeLessThanOrEqual(8);
      expect(Math.abs(position.y)).toBeLessThanOrEqual(5);
      expect(Math.abs(position.zoom)).toBeLessThanOrEqual(0.07);
    }

    expect(cameraMotion(15)).not.toEqual(cameraMotion(0));
    expect(cameraMotion(15)).toEqual(cameraMotion(15));
  });

  it('fits portrait and landscape images without creating a wider blank panel', () => {
    for (const ratio of [9 / 16, 1, 16 / 9, NaN]) {
      const size = floatingImageSize(900, 500, ratio, 400);
      expect(size.width).toBeLessThanOrEqual(352);
      expect(size.height).toBeLessThanOrEqual(352);
      expect(size.width / size.height).toBeCloseTo(Number.isNaN(ratio) ? 9 / 16 : ratio);
    }

    expect(floatingImageSize(300, 700, 9 / 16, 150).width).toBeLessThanOrEqual(108);
    expect(floatingImageSize(1400, 1000, 9 / 16, 400).width).toBe(352);
  });

  it('preserves every narrated word and its qualifications in short caption groups', () => {
    const text =
      'Sources disagree about part of this account. This is an interpretation, not an attested route. More details remain uncertain.';
    const track = storyCaptions(text);
    expect(track.captions.map((cue) => cue.text).join(' ')).toBe(text);
    expect(track.captions.every((cue) => cue.end - cue.start <= 12)).toBe(true);
    expect(activeStoryCaption(text, 0, 20)).toBe(track.captions[0].text);
    expect(activeStoryCaption(text, 19, 20)).toBe(track.captions.at(-1)?.text);
    expect(activeStoryCaption('', 0)).toBe('');
    expect(activeStoryCaption(text, 200, 20)).toBe(track.captions.at(-1)?.text);
  });
});
