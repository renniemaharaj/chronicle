import type { SpeechInput, SpeechSession } from './session';
import { speechChunks } from './chunks';

export function googleSession(input: SpeechInput): SpeechSession {
  const request = new AbortController();
  const audio = new Audio();
  const chunks = speechChunks(input.text);
  let index = 0;
  let cancelled = false;
  let paused = false;
  let url = '';
  audio.playbackRate = input.speed;
  audio.onended = () => {
    if (cancelled) {
      return;
    }
    if (index < chunks.length) {
      void loadNext();
    } else {
      input.done();
    }
  };
  audio.onerror = () => {
    if (!cancelled) {
      input.fail('Google speech audio could not be played.');
    }
  };
  const play = () =>
    audio.play().catch(() => {
      if (!cancelled && !paused) {
        input.fail('Audio playback was blocked. Press Play to try again.');
      }
    });

  async function loadNext() {
    try {
      const response = await fetch('/api/speech/synthesize', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        signal: request.signal,
        body: JSON.stringify({ text: chunks[index++], voiceName: input.voice, languageCode: input.language }),
      });
      if (!response.ok) {
        const body = await response.json().catch(() => ({}));
        throw new Error(body.error || 'Google speech request failed.');
      }
      const blob = await response.blob();
      if (cancelled) {
        return;
      }
      if (url) {
        URL.revokeObjectURL(url);
      }
      url = URL.createObjectURL(blob);
      audio.src = url;
      audio.playbackRate = input.speed;
      if (!paused) {
        await play();
      }
    } catch (error) {
      if (!cancelled) {
        input.fail((error as Error).message);
      }
    }
  }

  void loadNext();
  return {
    pause: () => {
      paused = true;
      audio.pause();
    },
    resume: () => {
      paused = false;
      if (url) {
        void play();
      }
    },
    cancel: () => {
      cancelled = true;
      request.abort();
      audio.pause();
      audio.removeAttribute('src');
      if (url) {
        URL.revokeObjectURL(url);
      }
    },
  };
}
