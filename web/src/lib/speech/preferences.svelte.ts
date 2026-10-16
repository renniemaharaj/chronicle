export type SpeechProvider = 'browser' | 'google' | 'screen';

class SpeechPreferences {
  provider = $state<SpeechProvider>('browser');
  voice = $state('');
  language = $state('en-US');
  status = $state('Ready to read');
  error = $state('');

  load() {
    try {
      const saved = JSON.parse(localStorage.getItem('chronicle-speech') ?? '{}');
      if (['browser', 'google', 'screen'].includes(saved.provider)) {
        this.provider = saved.provider;
      }
      this.voice = typeof saved.voice === 'string' ? saved.voice : '';
      this.language = typeof saved.language === 'string' ? saved.language : 'en-US';
    } catch {
      /* Use local defaults when preferences are unavailable. */
    }
  }

  save() {
    localStorage.setItem(
      'chronicle-speech',
      JSON.stringify({
        provider: this.provider,
        voice: this.voice,
        language: this.language,
      }),
    );
  }
}

export const speech = new SpeechPreferences();
