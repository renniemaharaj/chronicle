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
  let completedSeconds = 0;
  let progressTimer: ReturnType<typeof setInterval> | undefined;

  const progress = () => {
    if (!cancelled && !paused && loaded && Number.isFinite(audio.currentTime)) {
      input.progress?.(completedSeconds + audio.currentTime);
    }
  };

  audio.ontimeupdate = progress;
  audio.onloadedmetadata = () => {
    if (!cancelled && chunks.length === 1 && Number.isFinite(audio.duration)) {
      input.duration?.(audio.duration);
    }
  };

  const reportProgress = () => {
    if (!cancelled && !paused && loaded && !progressTimer) {
      progressTimer = setInterval(progress, 50);
    }
  };

  const stopProgress = () => {
    clearInterval(progressTimer);
    progressTimer = undefined;
  };

  const play = () =>
    audio
      .play()
      .then(reportProgress)
      .catch(() => {
        if (!cancelled && !paused) {
          input.fail('Audio playback was blocked. Press Play to try again.');
        }
      });

  audio.onended = () => {
    stopProgress();
    if (cancelled) {
      return;
    }

    completedSeconds += Number.isFinite(audio.duration) ? audio.duration : audio.currentTime;
    if (index < chunks.length) {
      void loadNext();
    } else {
      input.done();
    }
  };
  audio.onerror = () => {
    stopProgress();
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
      stopProgress();
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
      stopProgress();
      cancelled = true;
      request.abort();
      audio.pause();
      audio.removeAttribute('src');
    },
  };
}
