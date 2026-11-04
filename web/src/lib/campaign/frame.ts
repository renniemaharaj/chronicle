import { tick } from 'svelte';

export async function prepareCaptureFrame() {
  await tick();
  await new Promise((resolve) => requestAnimationFrame(() => requestAnimationFrame(resolve)));

  const stage = document.querySelector<HTMLElement>('.campaign-stage');
  if (!stage) {
    throw new Error('The reel preview is unavailable.');
  }

  const bounds = stage.getBoundingClientRect();
  const children = [...stage.querySelectorAll<HTMLElement>('h2, h3, p, footer')];
  const clipped = children.some((element) => {
    const child = element.getBoundingClientRect();

    return child.bottom > bounds.bottom || child.right > bounds.right;
  });
  if (clipped || stage.scrollHeight > stage.clientHeight + 1) {
    throw new Error('This passage does not fit the capture frame. Adjust the composition before recording.');
  }
}
