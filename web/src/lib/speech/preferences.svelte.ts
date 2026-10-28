export type SpeechProvider = 'browser' | 'google' | 'elevenlabs' | 'screen';

class SpeechPreferences {
  provider = $state<SpeechProvider>('browser');
  voice = $state('');
  voices = $state<Partial<Record<SpeechProvider, string>>>({});
  language = $state('en-US');
  model = $state('eleven_multilingual_v2');
  status = $state('Ready to read');
  error = $state('');

  load() {
    try {
      const saved = JSON.parse(localStorage.getItem('chronicle-speech') ?? '{}');
      if (['browser', 'google', 'elevenlabs', 'screen'].includes(saved.provider)) {
        this.provider = saved.provider;
      }

      for (const provider of ['browser', 'google', 'elevenlabs', 'screen'] as SpeechProvider[]) {
        if (typeof saved.voices?.[provider] === 'string') {
          this.voices[provider] = saved.voices[provider];
        }
      }

      this.voice = typeof saved.voice === 'string' ? saved.voice : '';
      this.voices[this.provider] = this.voice;
      this.model = typeof saved.model === 'string' ? saved.model : 'eleven_multilingual_v2';
      this.language = typeof saved.language === 'string' ? saved.language : 'en-US';
    } catch {
      /* Use local defaults when preferences are unavailable. */
    }
  }

  chooseProvider() {
    this.voice = this.voices[this.provider] ?? '';
    this.save();
  }

  save() {
    this.voices[this.provider] = this.voice;

    localStorage.setItem(
      'chronicle-speech',
      JSON.stringify({
        provider: this.provider,
        voice: this.voice,
        voices: this.voices,
        language: this.language,
        model: this.model,
      }),
    );
  }
}

export const speech = new SpeechPreferences();
