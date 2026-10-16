export type StoryScene = 'Map' | 'Timeline' | 'Graph';
export type StoryPanel = 'Chapters' | 'Places' | 'Sources';

export const storyView = $state({
  scene: 'Map' as StoryScene,
  panel: 'Chapters' as StoryPanel,
  sidebarOpen: true,
  followLocation: true,
});
