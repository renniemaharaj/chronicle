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
  utterance.onend = () => {
    if (!cancelled) {
      input.done();
    }
  };
  utterance.onerror = (event) => {
    if (!cancelled) {
      input.fail(`Speech stopped: ${event.error}`);
    }
  };
  speechSynthesis.speak(utterance);
  return {
    pause: () => speechSynthesis.pause(),
    resume: () => speechSynthesis.resume(),
    cancel: () => {
      cancelled = true;
      speechSynthesis.cancel();
    },
  };
}

export function readingSession(input: SpeechInput): SpeechSession {
  let remaining = Math.max(3000, (input.text.split(/\s+/).length * 400) / input.speed);
  let started = 0;
  let timer: ReturnType<typeof setTimeout>;
  const resume = () => {
    started = Date.now();
    timer = setTimeout(input.done, remaining);
  };
  resume();
  return {
    pause: () => {
      clearTimeout(timer);
      remaining = Math.max(0, remaining - (Date.now() - started));
    },
    resume,
    cancel: () => clearTimeout(timer),
  };
}
