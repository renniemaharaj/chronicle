import { browserSession, readingSession, type SpeechInput, type SpeechSession } from './session';
import { googleSession } from './google-session';
import type { SpeechProvider } from './preferences.svelte';

export class StoryPlayer {
  private session?: SpeechSession;
  private key = '';
  private paused = true;
  private generation = 0;

  update(key: string, playing: boolean, provider: SpeechProvider, input: SpeechInput) {
    if (key !== this.key) {
      this.stop();
      this.key = key;
    }
    if (!playing) {
      if (!this.paused) {
        this.session?.pause();
      }
      this.paused = true;
      return;
    }
    if (this.session) {
      if (this.paused) {
        this.session.resume();
      }
      this.paused = false;
      return;
    }
    const generation = ++this.generation;
    const guarded = {
      ...input,
      done: () => {
        if (generation !== this.generation) {
          return;
        }
        this.stop();
        input.done();
      },
      fail: (message: string) => {
        if (generation !== this.generation) {
          return;
        }
        this.stop();
        input.fail(message);
      },
    };
    try {
      this.paused = false;
      this.session =
        provider === 'screen'
          ? readingSession(guarded)
          : provider === 'google'
            ? googleSession(guarded)
            : browserSession(guarded);
    } catch (error) {
      guarded.fail((error as Error).message);
    }
  }

  stop() {
    this.generation++;
    this.session?.cancel();
    this.session = undefined;
    this.key = '';
    this.paused = true;
  }
}
