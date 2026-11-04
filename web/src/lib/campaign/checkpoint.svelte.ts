import { campaignAPI } from './api';
import type { Action, Campaign } from './types';
import type { Platform } from './platforms';

export class CampaignCheckpoint {
  publicationPlatform = $state<Platform>('facebook');
  campaign = $state<Campaign>();
  partNumber = $state(1);
  passageIndex = $state(0);
  busy = $state(false);
  recording = $state(false);
  error = $state('');
  pending = $state<{ file: string; duration: number }>();

  get part() {
    return this.campaign?.parts[this.partNumber - 1];
  }

  async act(action: Action) {
    if (!this.campaign) {
      return;
    }

    this.campaign = await campaignAPI.update(this.campaign, this.partNumber, action);
  }

  async saveTake() {
    const pending = this.pending;
    const campaign = this.campaign;
    const takeId = this.part?.takeId;
    if (!pending || !campaign) {
      return;
    }

    try {
      await this.act({
        action: 'complete',
        takeId,
        ...pending,
      });
      this.pending = undefined;
    } catch (error) {
      const latest = await campaignAPI.get(campaign.id).catch(() => undefined);
      const part = latest?.parts[this.partNumber - 1];
      if (latest && part && part.file === pending.file && part.takeId === takeId) {
        this.campaign = latest;
        this.pending = undefined;

        return;
      }

      if (latest) {
        this.campaign = latest;
      }

      throw error;
    }
  }
}
