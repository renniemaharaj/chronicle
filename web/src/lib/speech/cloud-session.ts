import type { SpeechInput, SpeechSession } from './session';
import { speechChunks } from './chunks';
import { prepareAudio, type CloudProvider } from './audio-api';

export function cloudSession(input: SpeechInput, provider: CloudProvider): SpeechSession {
  const request = new AbortController();
  const audio = new Audio();
  const chunks = speechChunks(input.text);
  let index = 0;
  let cancelled = false;
  let paused = false;
  let loaded = false;

  const play = () =>
    audio.play().catch(() => {
      if (!cancelled && !paused) {
        input.fail('Audio playback was blocked. Press Play to try again.');
      }
    });

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
      input.fail('Saved narration could not be played. Check the server connection.');
    }
  };

  async function loadNext() {
    loaded = false;
    try {
      const saved = await prepareAudio(chunks[index++], provider, input, false, request.signal);
      if (cancelled) {
        return;
      }

      if (!saved.available) {
        throw new Error('This passage needs audio. Open Voice & playback and generate missing audio first.');
      }

      audio.src = saved.url;
      audio.playbackRate = input.speed;
      loaded = true;
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
      if (loaded) {
        void play();
      }
    },
    cancel: () => {
      cancelled = true;
      request.abort();
      audio.pause();
      audio.removeAttribute('src');
    },
  };
}
