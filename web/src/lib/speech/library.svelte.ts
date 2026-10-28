import { prepareAudio, SpeechRequestError, type CloudProvider } from './audio-api';

export interface AudioPart {
  text: string;
  voice: string;
  language: string;
  model: string;
  passageId?: string;
  title?: string;
}

const activePreparations = new Set<NarrationLibrary>();

export class NarrationLibrary {
  saved = $state<boolean[]>([]);
  busy = $state(false);
  generating = $state(false);
  checked = $state(false);
  error = $state('');
  message = $state('');
  stopRequested = $state(false);
  blocked = $state('');
  private disposed = false;

  constructor(
    readonly provider: CloudProvider,
    readonly parts: AudioPart[],
  ) {}

  get missing() {
    return this.parts.filter((_, index) => !this.saved[index]);
  }

  get characters() {
    return this.missing.reduce((total, part) => total + Array.from(part.text).length, 0);
  }

  get completed() {
    return this.parts.length - this.missing.length;
  }

  get percent() {
    if (!this.parts.length) {
      return 0;
    }

    return Math.floor((this.completed / this.parts.length) * 100);
  }

  async check() {
    this.busy = true;
    this.checked = false;
    this.error = '';
    try {
      for (const [index, part] of this.parts.entries()) {
        if (this.disposed) {
          return;
        }

        const result = await prepareAudio(part.text, this.provider, part, false);
        this.saved[index] = result.available;
      }

      this.checked = true;
    } catch (cause) {
      this.error = (cause as Error).message;
    } finally {
      this.busy = false;
    }
  }

  async generate() {
    activePreparations.add(this);
    this.busy = true;
    this.generating = true;
    this.stopRequested = false;
    this.error = '';
    this.message = '';
    this.blocked = '';
    try {
      for (const [index, part] of this.parts.entries()) {
        if (this.stopRequested || this.disposed) {
          break;
        }

        if (this.saved[index]) {
          continue;
        }

        const result = await prepareAudio(part.text, this.provider, part, true);
        this.saved[index] = result.available;
      }

      this.message =
        this.stopRequested || this.missing.length
          ? 'Generation stopped. Completed files are saved.'
          : 'Audio is ready to play.';
    } catch (cause) {
      this.blocked = cause instanceof SpeechRequestError ? cause.code : 'request_failed';
      this.error = `${(cause as Error).message} Completed files are saved; retry to generate only missing parts.`;
    } finally {
      this.busy = false;
      this.generating = false;
      activePreparations.delete(this);
    }
  }

  stop() {
    // Let the current paid request finish and save its file before stopping.
    this.stopRequested = true;
  }

  dispose() {
    this.disposed = true;
    this.stop();
  }
}

export function stopNarrationPreparation() {
  for (const library of activePreparations) {
    library.stop();
  }
}
