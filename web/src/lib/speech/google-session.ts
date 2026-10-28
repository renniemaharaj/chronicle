import type { SpeechInput } from './session';
import { cloudSession } from './cloud-session';

export function googleSession(input: SpeechInput) {
  return cloudSession(input, 'google');
}
