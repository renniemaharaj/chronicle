import { prepareCaptureFrame } from './frame';
import { readPart } from './playback';
import { CampaignCheckpoint } from './checkpoint.svelte';
import { Recorder } from './obs';
import { StoryPlayer } from '../speech/player';
import { speech } from '../speech/preferences.svelte';

export class Studio extends CampaignCheckpoint {
  previewing = $state(false);
  previewAll = false;
  private recorder = new Recorder();
  private player = new StoryPlayer();
  private generation = 0;
  private activeTakeId = '';

  async capture(password: string, sceneName: string) {
    this.busy = true;
    this.error = '';
    this.passageIndex = 0;

    const generation = ++this.generation;
    try {
      speech.load();
      if (speech.provider === 'screen') {
        throw new Error('Choose an audio voice in Voice & playback before recording.');
      }

      await this.recorder.connect(password, sceneName);
      if (generation !== this.generation) {
        await this.recorder.disconnect();

        return;
      }

      await this.act({ action: 'begin' });
      this.activeTakeId = this.part?.takeId ?? '';
      await prepareCaptureFrame();
      if (generation !== this.generation) {
        await this.abort();

        return;
      }

      await this.recorder.start();
      if (generation !== this.generation) {
        await this.abort();

        return;
      }

      this.recording = true;
      this.read(generation);
    } catch (error) {
      this.error = (error as Error).message;
      await this.abort();
    } finally {
      this.busy = false;
    }
  }

  private read(generation: number) {
    readPart(this, this.player, String(generation), () => generation === this.generation);
  }

  async preview(all = false) {
    this.error = '';
    this.passageIndex = 0;
    this.previewAll = all;
    this.previewing = true;

    const generation = ++this.generation;
    try {
      speech.load();
      await prepareCaptureFrame();
      this.read(generation);
    } catch (error) {
      this.error = (error as Error).message;
      await this.abort();
    }
  }

  async finish() {
    this.busy = true;
    this.player.stop();
    try {
      this.pending = await this.recorder.stop();
      await this.saveTake();
    } catch (error) {
      this.error = (error as Error).message;
    } finally {
      this.recording = false;
      this.busy = false;
      await this.recorder.disconnect();
    }
  }

  async abort() {
    this.previewing = false;
    this.generation++;
    this.player.stop();
    try {
      await this.recorder.discard();
      if (!this.pending && this.part?.status === 'recording' && this.part.takeId === this.activeTakeId) {
        await this.act({
          action: 'interrupt',
          takeId: this.part.takeId,
        });
      }
    } catch (error) {
      this.error = `${this.error} Recovery needed: ${(error as Error).message}`;
    } finally {
      this.recording = false;
      await this.recorder.disconnect();
    }
  }
}
