import { afterEach, describe, expect, it, vi } from 'vitest';
import { googleSession } from './google-session';
import { speechChunks } from './chunks';

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

describe('Google audio playback', () => {
  it('bounds Unicode chunks without splitting characters', () => {
    const text = 'Céline 🎵 '.repeat(900);
    const chunks = speechChunks(text);
    expect(chunks.length).toBeGreaterThan(1);
    expect(chunks.every((chunk) => new TextEncoder().encode(chunk).length <= 4500)).toBe(true);
    expect(chunks.join(' ').replaceAll(' ', '')).toBe(text.replaceAll(' ', ''));
    expect(chunks.join('')).not.toContain('\uFFFD');
  });

  it('keeps pending audio paused, advances all chunks, and releases audio on cancel', async () => {
    vi.stubGlobal('Audio', TestAudio);
    const request = vi.fn().mockResolvedValue({ ok: true, blob: async () => new Blob(['mock-mp3']) });
    vi.stubGlobal('fetch', request);
    vi.spyOn(URL, 'createObjectURL').mockReturnValue('blob:test');
    const revoke = vi.spyOn(URL, 'revokeObjectURL').mockImplementation(() => {});
    const done = vi.fn();
    const session = googleSession({
      text: 'A researched event. '.repeat(500),
      voice: 'en-US-Test',
      language: 'en-US',
      speed: 1,
      done,
      fail: vi.fn(),
    });
    session.pause();
    await vi.waitFor(() => expect(TestAudio.current.src).toBe('blob:test'));
    expect(TestAudio.current.play).not.toHaveBeenCalled();
    session.resume();
    expect(TestAudio.current.play).toHaveBeenCalledOnce();
    const count = speechChunks('A researched event. '.repeat(500)).length;
    for (let index = 1; index <= count; index++) {
      TestAudio.current.onended();
      if (index < count) await vi.waitFor(() => expect(request).toHaveBeenCalledTimes(index + 1));
      await Promise.resolve();
      await Promise.resolve();
    }
    await vi.waitFor(() => expect(done).toHaveBeenCalledOnce());
    session.cancel();
    expect(revoke).toHaveBeenCalled();
    const options = request.mock.calls[0][1];
    expect(options.signal.aborted).toBe(true);
  });
});
