<script lang="ts">
  import { onMount, untrack } from 'svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { focusChapter } from '../../state/chapter-focus';
  import { speech } from '../../speech/preferences.svelte';
  import { StoryPlayer } from '../../speech/player';
  import { eventStory, contextFor, contextText } from '../../speech/story';

  const player = new StoryPlayer();
  let ready = $state(false);
  let previousChapter = '';
  let previousProject = '';
  let contextIntro = '';
  let chapterKey = '';
  onMount(() => {
    speech.load();
    ready = true;
    return () => player.stop();
  });
  $effect(() => {
    if (!ready) {
      return;
    }
    const event = ws.event;
    const chapter = ws.chapter;
    const playing = ws.playing;
    const provider = speech.provider;
    const voice = speech.voice;
    const language = speech.language;
    const speed = ws.speed;
    const text = eventStory(ws.project, event, ws.activeClaims);
    const contexts = contextFor(ws.project, event);
    const projectId = ws.project.id;
    untrack(() => {
      if (previousProject !== projectId) {
        player.stop();
        previousProject = projectId;
        previousChapter = '';
        chapterKey = '';
        contextIntro = '';
      }
      if (playing && chapter && chapter.id !== previousChapter) {
        contextIntro = contexts.map((item) => `${item.title}. ${contextText(ws.project, item)}`).join(' ');
        chapterKey = event?.id ?? '';
        previousChapter = chapter.id;
      }
      if (playing && chapter && ws.followChapter && ws.focusedChapterId !== chapter.id) {
        focusChapter(ws, chapter.id, false);
      }
      const narration = `${event?.title ?? ''}. ${event?.id === chapterKey ? contextIntro : ''} ${text}`;
      const key = JSON.stringify([projectId, event?.id, provider, voice, language, speed, narration]);
      if (playing) {
        speech.status = provider === 'screen' ? 'Reading on screen' : 'Reading aloud';
        speech.error = '';
      } else if (speech.status !== 'Speech unavailable') {
        speech.status = ws.storyFinished ? 'Story complete' : 'Paused';
      }
      player.update(key, playing && !!event, provider, {
        text: narration,
        voice,
        language,
        speed,
        done: () => {
          speech.status = 'Ready to read';
          ws.step(1);
        },
        fail: (message) => {
          ws.playing = false;
          speech.error = message;
          speech.status = 'Speech unavailable';
          ws.notify(message);
        },
      });
    });
  });
</script>
