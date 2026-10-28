import { workspace as ws } from './workspace.svelte';
import { toggleStory } from './playback-controls';
import { closeStoryMode } from './story-mode';
import { saveProject } from './persistence';

export function handleKeyboard(event: KeyboardEvent) {
  if (event.key === 'Escape') {
    if (ws.dialog) {
      ws.dialog = '';
    } else if (ws.view === 'Story' && !document.fullscreenElement) {
      closeStoryMode(ws);
    }
  }

  if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === 'k') {
    event.preventDefault();
    ws.dialog = 'Research';
  }

  if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === 's') {
    event.preventDefault();
    saveProject();
  }

  const target = event.target as HTMLElement;
  if (!ws.hasProject || ws.dialog || target?.closest('input, textarea, select, [contenteditable="true"]')) {
    return;
  }

  if (event.key === 'ArrowRight') {
    ws.step(1);
  }

  if (event.key === 'ArrowLeft') {
    ws.step(-1);
  }

  if (event.code === 'Space' && !ws.dialog && !(event.target as HTMLElement)?.closest('button')) {
    event.preventDefault();
    toggleStory(ws);
  }
}
