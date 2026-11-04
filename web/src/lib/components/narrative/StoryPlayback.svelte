<script lang="ts">
  import { onMount, untrack } from 'svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { focusChapter } from '../../state/chapter-focus';
  import { speech } from '../../speech/preferences.svelte';
  import { StoryPlayer } from '../../speech/player';
  import { storyPassages } from '../../speech/passages';
  import { storyPlayback } from '../../state/story-playback.svelte';

  const player = new StoryPlayer();
  let ready = $state(false);
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
    const model = speech.model;
    const speed = ws.speed;
    const passages = storyPassages(ws.project, ws.orderedEvents, ws.activeClaims);
    const narration = passages.find((passage) => passage.id === event?.id)?.text ?? '';
    const projectId = ws.project.id;
    const restartToken = storyPlayback.restartToken;
    untrack(() => {
      if (playing && chapter && ws.followChapter && ws.focusedChapterId !== chapter.id) {
        focusChapter(ws, chapter.id, false);
      }

      const key = JSON.stringify([
        projectId,
        event?.id,
        provider,
        voice,
        language,
        model,
        speed,
        narration,
        restartToken,
      ]);
      if (storyPlayback.passageId !== key) {
        storyPlayback.passageId = key;
        storyPlayback.eventId = event?.id ?? '';
        storyPlayback.projectId = projectId;
        storyPlayback.elapsed = 0;
        storyPlayback.duration = 0;
      }

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
        model,
        speed,
        progress: (seconds) => {
          storyPlayback.elapsed = seconds;
        },
        duration: (seconds) => {
          storyPlayback.duration = seconds;
        },
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
