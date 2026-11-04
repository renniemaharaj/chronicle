import { passageClock } from './passage-clock';

export interface SpeechSession {
  pause(): void;
  resume(): void;
  cancel(): void;
}

export interface SpeechInput {
  text: string;
  voice: string;
  language: string;
  speed: number;
  model?: string;
  progress?(seconds: number): void;
  duration?(seconds: number): void;
  done(): void;
  fail(message: string): void;
}

export function browserSession(input: SpeechInput): SpeechSession {
  if (!('speechSynthesis' in window)) {
    throw new Error('This browser has no speech engine. Choose Google Cloud or read on screen.');
  }

  const utterance = new SpeechSynthesisUtterance(input.text);
  const voice = speechSynthesis.getVoices().find((item) => item.voiceURI === input.voice);
  if (input.voice && !voice) {
    throw new Error('The selected browser voice is unavailable. Choose another voice.');
  }

  if (voice) {
    utterance.voice = voice;
  }

  utterance.lang = input.language;
  utterance.rate = input.speed;

  let cancelled = false;
  let paused = false;
  let started = false;
  const clock = passageClock(input);
  utterance.onstart = () => {
    started = true;
    if (!paused) {
      clock.resume();
    }
  };
  utterance.onend = () => {
    clock.pause();
    if (!cancelled) {
      input.done();
    }
  };
  utterance.onerror = (event) => {
    clock.pause();
    if (!cancelled) {
      input.fail(`Speech stopped: ${event.error}`);
    }
  };
  speechSynthesis.speak(utterance);

  return {
    pause: () => {
      paused = true;
      clock.pause();
      speechSynthesis.pause();
    },
    resume: () => {
      paused = false;
      if (started) {
        clock.resume();
      }

      speechSynthesis.resume();
    },
    cancel: () => {
      clock.pause();
      cancelled = true;
      speechSynthesis.cancel();
    },
  };
}

export function readingSession(input: SpeechInput): SpeechSession {
  const clock = passageClock(input);
  let remaining = Math.max(3000, (input.text.split(/\s+/).length * 400) / input.speed);
  input.duration?.((remaining / 1000) * input.speed);

  let started = 0;
  let timer: ReturnType<typeof setTimeout>;

  const resume = () => {
    clock.resume();
    started = Date.now();
    timer = setTimeout(input.done, remaining);
  };

  resume();

  return {
    pause: () => {
      clock.pause();
      clearTimeout(timer);
      remaining = Math.max(0, remaining - (Date.now() - started));
    },
    resume,
    cancel: () => {
      clock.pause();
      clearTimeout(timer);
    },
  };
}
