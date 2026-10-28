import { afterEach, describe, expect, it, vi } from 'vitest';
import { googleSession } from './google-session';
import { cloudSession } from './cloud-session';
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

function mockAudio(available = true) {
  vi.stubGlobal('Audio', TestAudio);

  const request = vi.fn().mockResolvedValue({
    ok: true,
    json: async () => ({
      available,
      url: '/api/speech/audio/test',
    }),
  });
  vi.stubGlobal('fetch', request);

  return request;
}

describe('saved cloud audio playback', () => {
  it('bounds Unicode chunks without splitting characters', () => {
    const text = 'Céline 🎵 '.repeat(900);
    const chunks = speechChunks(text);
    expect(chunks.length).toBeGreaterThan(1);
    expect(chunks.every((chunk) => new TextEncoder().encode(chunk).length <= 4500)).toBe(true);
    expect(chunks.join(' ').replaceAll(' ', '')).toBe(text.replaceAll(' ', ''));
    expect(chunks.join('')).not.toContain('\uFFFD');
  });

  it('keeps pending audio paused, advances all chunks, and releases audio on cancel', async () => {
    const request = mockAudio();
    const done = vi.fn();
    const session = googleSession({
      text: 'A researched event. '.repeat(500),
      voice: 'en-US-Test',
      language: 'en-US',
      speed: 1.5,
      done,
      fail: vi.fn(),
    });
    session.pause();
    await vi.waitFor(() => expect(TestAudio.current.src).toBe('/api/speech/audio/test'));
    expect(TestAudio.current.play).not.toHaveBeenCalled();
    expect(TestAudio.current.playbackRate).toBe(1.5);
    session.resume();
    expect(TestAudio.current.play).toHaveBeenCalledOnce();

    const count = speechChunks('A researched event. '.repeat(500)).length;
    for (let index = 1; index <= count; index++) {
      TestAudio.current.onended();
      if (index < count) {
        await vi.waitFor(() => expect(request).toHaveBeenCalledTimes(index + 1));
      }

      await Promise.resolve();
      await Promise.resolve();
    }

    await vi.waitFor(() => expect(done).toHaveBeenCalledOnce());
    session.cancel();
    expect(TestAudio.current.removeAttribute).toHaveBeenCalledWith('src');
    expect(request.mock.calls[0][1].signal.aborted).toBe(true);
    expect(JSON.parse(request.mock.calls[0][1].body).generate).toBe(false);
  });

  it('does not generate missing audio implicitly or advance the story', async () => {
    const request = mockAudio(false);
    const done = vi.fn();
    const fail = vi.fn();
    cloudSession(
      {
        text: 'Synthetic narration.',
        voice: 'test-voice',
        language: 'en-US',
        model: 'test-model',
        speed: 1,
        done,
        fail,
      },
      'elevenlabs',
    );
    await vi.waitFor(() => expect(fail).toHaveBeenCalledOnce());
    expect(fail.mock.calls[0][0]).toContain('generate missing audio');
    expect(done).not.toHaveBeenCalled();
    expect(TestAudio.current.play).not.toHaveBeenCalled();
    expect(JSON.parse(request.mock.calls[0][1].body)).toMatchObject({
      provider: 'elevenlabs',
      modelId: 'test-model',
      generate: false,
    });
  });
});
