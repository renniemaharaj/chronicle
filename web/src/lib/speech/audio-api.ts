import type { SpeechInput } from './session';

export type CloudProvider = 'google' | 'elevenlabs';

export interface SavedAudio {
  available: boolean;
  cached: boolean;
  url: string;
}

export class SpeechRequestError extends Error {
  constructor(
    message: string,
    readonly code = 'request_failed',
  ) {
    super(message);
  }
}

export async function speechRequest<T>(path: string, options?: RequestInit): Promise<T> {
  const response = await fetch(`/api/speech/${path}`, options);
  const body = await response.json().catch(() => ({ error: 'Speech service unavailable. Try again later.' }));
  if (!response.ok) {
    throw new SpeechRequestError(body.error || 'Speech request failed.', body.code);
  }

  return body as T;
}

export function prepareAudio(
  text: string,
  provider: CloudProvider,
  input: Pick<SpeechInput, 'voice' | 'language' | 'model'>,
  generate: boolean,
  signal?: AbortSignal,
) {
  return speechRequest<SavedAudio>('prepare', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    signal,
    body: JSON.stringify({
      text,
      provider,
      voiceName: input.voice,
      languageCode: input.language,
      modelId: input.model,
      generate,
    }),
  });
}
