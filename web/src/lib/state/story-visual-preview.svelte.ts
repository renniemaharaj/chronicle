import type { StoryVisual } from '../types/media';

export const visualPreview = $state({
  projectId: '',
  eventId: '',
  elapsed: 0,
  active: false,
  running: false,
});

let timer: ReturnType<typeof setInterval> | undefined;

export function cancelVisualPreview() {
  clearInterval(timer);
  timer = undefined;
  visualPreview.active = false;
  visualPreview.running = false;
}

export function previewVisual(projectId: string, eventId: string, cue: StoryVisual) {
  cancelVisualPreview();

  const started = Date.now();
  const offset = Math.max(0, cue.startSeconds - 0.7);
  Object.assign(visualPreview, {
    projectId,
    eventId,
    elapsed: offset,
    active: true,
    running: true,
  });
  timer = setInterval(() => {
    visualPreview.elapsed = offset + (Date.now() - started) / 1000;
    if (visualPreview.elapsed >= cue.startSeconds + cue.durationSeconds + 0.3) {
      clearInterval(timer);
      timer = undefined;
      visualPreview.running = false;
    }
  }, 50);
}
