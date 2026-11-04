import { prepareCaptureFrame } from './frame';
import { speech } from '../speech/preferences.svelte';
import type { StoryPlayer } from '../speech/player';
import type { Studio } from './studio.svelte';

export function readPart(studio: Studio, player: StoryPlayer, key: string, current: () => boolean) {
  const passage = studio.part?.passages[studio.passageIndex];
  if (!passage || !current()) {
    return;
  }

  const fail = (message: string) => {
    studio.error = message;
    void studio.abort();
  };

  player.update(`${key}-${studio.partNumber}-${studio.passageIndex}`, true, speech.provider, {
    text: passage.text,
    voice: speech.voice,
    language: speech.language,
    model: speech.model,
    speed: 1,
    done: () => {
      if (!current()) {
        return;
      }

      if (studio.passageIndex + 1 < (studio.part?.passages.length ?? 0)) {
        studio.passageIndex++;
      } else if (
        studio.previewing &&
        studio.previewAll &&
        studio.partNumber < (studio.campaign?.parts.length ?? 0)
      ) {
        studio.partNumber++;
        studio.passageIndex = 0;
      } else if (studio.previewing) {
        studio.previewing = false;

        return;
      } else {
        void studio.finish();

        return;
      }

      prepareCaptureFrame()
        .then(() => readPart(studio, player, key, current))
        .catch((error) => fail((error as Error).message));
    },
    fail,
  });
}
