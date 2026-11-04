import { tick } from 'svelte';
import { campaignAPI } from './api';
import { Studio } from './studio.svelte';
import { workspace as ws } from '../state/workspace.svelte';
import type { Campaign } from './types';
import type { Platform } from './platforms';

export class CampaignWorkspace {
  studio = new Studio();
  campaigns = $state<Campaign[]>([]);
  tab = $state('Series');
  section = $state('Episodes');
  createOpen = $state(false);
  studioOpen = $state(false);
  search = $state('');
  filter = $state('all');

  get available() {
    return this.campaigns.filter((campaign) => campaign.projectId === ws.project.id);
  }

  get locked() {
    return this.studio.busy || this.studio.recording || this.studio.previewing || !!this.studio.pending;
  }

  async perform(action: () => Promise<void>) {
    this.studio.busy = true;
    this.studio.error = '';
    try {
      await action();
    } catch (error) {
      this.studio.error = (error as Error).message;
    } finally {
      this.studio.busy = false;
    }
  }

  private restore(value: Campaign) {
    this.studio.campaign = value;
    this.studio.partNumber = value.parts.find((part) => part.status !== 'published')?.number ?? 1;
    this.studio.passageIndex = 0;
  }

  async reload() {
    await this.perform(async () => {
      this.campaigns = await campaignAPI.list();

      const current =
        this.available.find((item) => item.id === this.studio.campaign?.id) ?? this.available[0];
      if (current) {
        this.restore(current);
      }
    });
  }

  async choose(id: string) {
    if (this.locked) {
      return;
    }

    await this.perform(async () => this.restore(await campaignAPI.get(id)));
  }

  async create(seconds: number) {
    await this.perform(async () => {
      const value = await campaignAPI.create(ws.project.id, seconds);
      this.campaigns = [value, ...this.campaigns];
      this.restore(value);
      this.createOpen = false;
      this.studioOpen = false;
      this.tab = 'Series';
    });
  }

  select(number: number) {
    if (this.locked) {
      return;
    }

    this.studio.partNumber = number;
    this.studio.passageIndex = 0;
    ws.rightOpen = true;
  }

  openStudio(number = this.studio.partNumber, platform: Platform = 'facebook') {
    if (this.locked) {
      return;
    }

    this.select(number);
    this.studio.publicationPlatform = platform;
    this.studioOpen = true;
  }

  async previewSeries() {
    this.openStudio(1);
    await tick();
    if (this.studio.campaign?.engine === 'ffmpeg') {
      const video = document.querySelector<HTMLVideoElement>('.reel-rendered-video');
      await video?.play();

      return;
    }

    await this.studio.preview(true);
  }
}
