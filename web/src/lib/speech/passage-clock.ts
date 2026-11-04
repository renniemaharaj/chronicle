import type { SpeechInput } from './session';

// Used for screen reading and browser speech; saved audio reports its actual media time.
export function passageClock(input: SpeechInput) {
  let elapsed = 0;
  let started = 0;
  let timer: ReturnType<typeof setInterval> | undefined;

  const report = () => {
    const now = Date.now();
    elapsed += ((now - started) / 1000) * input.speed;
    started = now;
    input.progress?.(elapsed);
  };

  return {
    resume: () => {
      if (timer || !input.progress) {
        return;
      }

      started = Date.now();
      timer = setInterval(report, 50);
    },
    pause: () => {
      if (timer) {
        report();
        clearInterval(timer);
        timer = undefined;
      }
    },
  };
}
