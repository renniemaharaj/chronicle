import { beforeEach, describe, expect, it, vi } from 'vitest';
import type { SpeechInput } from './session';

const mock = vi.hoisted(() => ({
  input: undefined as SpeechInput | undefined,
  pause: vi.fn(),
  resume: vi.fn(),
  cancel: vi.fn(),
  start: vi.fn(),
}));
vi.mock('./session', () => ({
  browserSession: (input: SpeechInput) => {
    mock.input = input;
    mock.start();
    return mock;
  },
  readingSession: (input: SpeechInput) => {
    mock.input = input;
    return mock;
  },
}));
vi.mock('./google-session', () => ({
  googleSession: (input: SpeechInput) => {
    mock.input = input;
    return mock;
  },
}));
import { StoryPlayer } from './player';

beforeEach(() => vi.clearAllMocks());

describe('deterministic story playback', () => {
  it('pauses and resumes without recreating the utterance', () => {
    const player = new StoryPlayer();
    const input = { text: 'A sourced event.', voice: '', language: 'en-US', speed: 1, done: vi.fn(), fail: vi.fn() };
    player.update('event-a', true, 'browser', input);
    player.update('event-a', false, 'browser', input);
    expect(mock.pause).toHaveBeenCalledOnce();
    player.update('event-a', true, 'browser', input);
    expect(mock.resume).toHaveBeenCalledOnce();
    expect(mock.start).toHaveBeenCalledOnce();
    mock.input?.done();
    expect(input.done).toHaveBeenCalledOnce();
    player.stop();
  });

  it('cancels skipped events and guards stale completions after a voice change', () => {
    const player = new StoryPlayer();
    const input = { text: 'First.', voice: '', language: 'en-US', speed: 1, done: vi.fn(), fail: vi.fn() };
    player.update('first-voice-a', true, 'browser', input);
    const old = mock.input;
    player.update('first-voice-b', true, 'browser', { ...input, voice: 'voice-b' });
    expect(mock.cancel).toHaveBeenCalledOnce();
    old?.done();
    old?.fail('Old error');
    expect(input.done).not.toHaveBeenCalled();
    expect(input.fail).not.toHaveBeenCalled();
    mock.input?.done();
    expect(input.done).toHaveBeenCalledOnce();
  });
});
