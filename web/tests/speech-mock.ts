import type { Page } from '@playwright/test';

export async function mockSpeech(page: Page) {
  await page.addInitScript(() => {
    const state = {
      texts: [] as string[],
      current: null as { onend?: () => void } | null,
      pauses: 0,
      resumes: 0,
      cancels: 0,
    };
    Object.assign(window, { __speech: state });
    Object.defineProperty(window, 'SpeechSynthesisUtterance', {
      value: class {
        constructor(public text: string) {}
      },
      configurable: true,
    });
    Object.defineProperty(window, 'speechSynthesis', {
      configurable: true,
      value: {
        getVoices: () => [
          { voiceURI: 'test-a', name: 'Test voice A', lang: 'en-US' },
          { voiceURI: 'test-b', name: 'Test voice B', lang: 'en-US' },
        ],
        addEventListener() {},
        removeEventListener() {},
        speak(utterance: { text: string; onend?: () => void }) {
          state.current = utterance;
          state.texts.push(utterance.text);
        },
        pause() {
          state.pauses++;
        },
        resume() {
          state.resumes++;
        },
        cancel() {
          state.cancels++;
        },
      },
    });
  });
}

export async function speechState(page: Page) {
  return page.evaluate(
    () => (window as unknown as { __speech: { texts: string[]; pauses: number; resumes: number } }).__speech,
  );
}

export async function finishReading(page: Page) {
  await page.evaluate(() =>
    (window as unknown as { __speech: { current: { onend(): void } } }).__speech.current.onend(),
  );
}
