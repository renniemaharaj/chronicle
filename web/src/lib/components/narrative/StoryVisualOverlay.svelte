<script lang="ts">
  import { onDestroy, onMount } from 'svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { storyView } from '../../state/story-view.svelte';
  import { presentationTime } from '../../state/story-presentation.svelte';
  import { floatingImageSize } from '../../utils/story-image-fit';
  import { visualPreview, cancelVisualPreview } from '../../state/story-visual-preview.svelte';
  import { imageURL, passageImage, visualOpacity } from '../../utils/story-visual';
  import { intervalLabel } from '../../types/temporal';

  let failed = $state('');
  let anchor: HTMLDivElement;
  let bounds = $state({
    width: 800,
    height: 600,
    reader: 400,
    matchReaderWidth: false,
  });
  let aspect = $state(9 / 16);
  let caption = $state<HTMLElement>();
  let captionHeight = $state(70);
  const event = $derived(ws.event);
  const asset = $derived(passageImage(ws.project, event));
  const previewing = $derived(
    visualPreview.active && visualPreview.projectId === ws.project.id && visualPreview.eventId === event?.id,
  );
  const seconds = $derived(presentationTime(ws.project.id, event?.id).seconds);
  const size = $derived(
    floatingImageSize(
      bounds.width,
      bounds.height,
      aspect,
      bounds.reader,
      captionHeight,
      bounds.matchReaderWidth,
    ),
  );
  const opacity = $derived(
    storyView.showImages && event?.visual && asset && failed !== asset.url
      ? visualOpacity(event.visual, seconds)
      : 0,
  );

  $effect(() => {
    anchor?.parentElement?.style.setProperty(
      '--story-image-space',
      opacity > 0 ? `${size.width + 14}px` : '0px',
    );
  });

  $effect(() => {
    if (visualPreview.active && (!previewing || ws.playing)) {
      cancelVisualPreview();
    }
  });
  $effect(() => {
    if (!caption) {
      return;
    }

    const element = caption;
    const observer = new ResizeObserver(() => {
      captionHeight = element.getBoundingClientRect().height;
    });
    observer.observe(element);

    return () => observer.disconnect();
  });
  onDestroy(cancelVisualPreview);
  onMount(() => {
    const measure = () => {
      const reader = anchor.parentElement?.querySelector('.story-reading');
      bounds = {
        width: anchor.clientWidth,
        matchReaderWidth: window.innerWidth >= 760,
        height: window.innerWidth >= 760 ? anchor.clientHeight : anchor.clientHeight * 0.54 + 10,
        reader:
          window.innerWidth >= 760 ? (reader?.getBoundingClientRect().width ?? 0) : anchor.clientWidth / 2,
      };
    };

    const observer = new ResizeObserver(measure);
    observer.observe(anchor);

    return () => observer.disconnect();
  });
</script>

<div
  class="story-visual-anchor"
  bind:this={anchor}
>
  {#if asset && opacity > 0}
    <figure
      class="story-visual-overlay"
      style:opacity
      style:width={`${size.width}px`}
      style:--story-image-height={`${size.height}px`}
      aria-label="Timed passage image"
      data-media-id={asset.id}
    >
      <img
        src={imageURL(asset)}
        alt={asset.caption || asset.title}
        referrerpolicy="no-referrer"
        onload={(event) => {
          const image = event.currentTarget as HTMLImageElement;
          aspect = image.naturalWidth / image.naturalHeight;
        }}
        onerror={() => {
          failed = asset.url;
        }}
      />
      <figcaption bind:this={caption}>
        <span class="story-visual-kind">
          {asset.generated ? 'AI illustration' : 'Source image'}
        </span>
        <strong>{asset.title}</strong>
        <span>{asset.caption}</span>
        <small>
          {asset.depicts ? intervalLabel(asset.depicts) + ' · ' : ''}{asset.creator || 'Creator not recorded'}
        </small>
      </figcaption>
    </figure>
  {/if}
</div>
