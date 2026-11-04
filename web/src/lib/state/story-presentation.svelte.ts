import { storyPlayback } from './story-playback.svelte';
import { visualPreview } from './story-visual-preview.svelte';

export const captureClock = $state({
  eventId: '',
  projectId: '',
  seconds: 0,
  duration: 0,
  active: false,
});

export function presentationTime(projectId: string, eventId = '') {
  if (captureClock.active && captureClock.projectId === projectId && captureClock.eventId === eventId) {
    return {
      seconds: captureClock.seconds,
      duration: captureClock.duration,
    };
  }

  if (visualPreview.active && visualPreview.projectId === projectId && visualPreview.eventId === eventId) {
    return {
      seconds: visualPreview.elapsed,
      duration: 0,
    };
  }

  if (storyPlayback.projectId === projectId && storyPlayback.eventId === eventId) {
    return {
      seconds: storyPlayback.elapsed,
      duration: storyPlayback.duration,
    };
  }

  return {
    seconds: 0,
    duration: 0,
  };
}
