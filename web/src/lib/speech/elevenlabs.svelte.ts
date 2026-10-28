import { speechRequest } from './audio-api';

export interface ElevenVoice {
  voice_id: string;
  name: string;
  category: string;
  description: string;
  labels: Record<string, string>;
}

export interface ElevenModel {
  model_id: string;
  name: string;
  maximum_text_length_per_request: number;
}

export interface ElevenUsage {
  tier: string;
  character_count: number;
  character_limit: number;
  next_character_count_reset_unix: number;
  status: string;
}

export class ElevenLabsAccount {
  settings = $state({
    hasApiKey: false,
    keySource: 'none',
  });
  voices = $state<ElevenVoice[]>([]);
  models = $state<ElevenModel[]>([]);
  usage = $state<ElevenUsage | null>(null);
  usageError = $state('');
  error = $state('');
  loading = $state(false);
  checked = $state(false);

  async load() {
    this.loading = true;
    this.error = '';
    try {
      this.settings = await speechRequest('settings');
      this.checked = true;
      if (this.settings.hasApiKey) {
        const results = await Promise.allSettled([
          speechRequest<ElevenVoice[]>('voices?provider=elevenlabs'),
          speechRequest<ElevenModel[]>('models'),
          this.loadUsage(),
        ]);
        const [voices, models] = results;
        if (voices.status === 'fulfilled') {
          this.voices = voices.value;
        } else {
          this.voices = [];
          this.error = voices.reason.message;
        }

        if (models.status === 'fulfilled') {
          this.models = models.value;
        } else {
          this.models = [];
          this.error = [this.error, models.reason.message].filter(Boolean).join(' ');
        }
      } else {
        this.voices = [];
        this.models = [];
        this.usage = null;
      }
    } catch (cause) {
      this.error = (cause as Error).message;
    } finally {
      this.loading = false;
    }
  }

  async loadUsage() {
    this.usageError = '';
    try {
      this.usage = await speechRequest<ElevenUsage>('usage');
    } catch (cause) {
      this.usage = null;
      this.usageError = `Usage unavailable: ${(cause as Error).message}`;
    }
  }

  async save(apiKey: string) {
    this.loading = true;
    this.error = '';
    try {
      await speechRequest('settings', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ apiKey }),
      });
      await this.load();
    } catch (cause) {
      this.error = (cause as Error).message;
    } finally {
      this.loading = false;
    }
  }
}
