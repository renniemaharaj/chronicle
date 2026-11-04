import { writeFile } from 'node:fs/promises';
import { resolve } from 'node:path';
import { backend } from './reel-backend.mjs';

function imageAddress(asset) {
  if (asset.generated && /^\/api\/media\/images\/[a-f0-9]{64}\.(png|jpg|gif)$/.test(asset.url)) {
    return `${backend}${asset.url}`;
  }

  const url = new URL(asset.url);
  if (!asset.retrieved || !['http:', 'https:'].includes(url.protocol)) {
    throw new Error('A reel image needs imported generation provenance or a retrieved source reference.');
  }

  return url.href;
}

async function download(asset, path) {
  const response = await fetch(imageAddress(asset), { signal: AbortSignal.timeout(30000) });
  if (!response.ok || !response.body) {
    throw new Error(`The passage image could not be loaded: ${asset.title}`);
  }

  const chunks = [];
  let size = 0;
  for await (const chunk of response.body) {
    size += chunk.length;
    if (size > 16 * 1024 * 1024) {
      throw new Error('The passage image exceeds 16 MB.');
    }

    chunks.push(chunk);
  }

  await writeFile(path, Buffer.concat(chunks));
}

function wrapped(value) {
  return (
    value
      .match(/.{1,40}(?:\s|$)|\S+/g)
      ?.map((line) => line.trim())
      .join('\n') ?? ''
  );
}

export async function reelImageOverlays(part, readings, directory, args, filters) {
  const labels = readings.map((_, index) => `[v${index}]`);
  let input = readings.length * 2;
  for (const [index, reading] of readings.entries()) {
    const visual = part.passages.find((item) => item.audioId === reading.audioId && item.visual)?.visual;
    if (!visual) {
      continue;
    }

    const { cue, image } = visual;
    const start = cue.startSeconds;
    const end = Math.min(start + cue.durationSeconds, reading.audioEnd);
    if (start >= end) {
      continue;
    }

    const path = resolve(directory, `visual-${part.number}-${index}.image`);
    const label = resolve(directory, `visual-${part.number}-${index}.txt`);
    await download(image, path);
    await writeFile(
      label,
      `${image.generated ? 'AI ILLUSTRATION | RECONSTRUCTION' : 'SOURCE IMAGE'}\n${wrapped(image.title)}`,
    );
    if (image.url.toLowerCase().split('?')[0].endsWith('.gif')) {
      args.push('-ignore_loop', '0');
    }

    args.push('-i', path);

    const fade = Math.min(cue.transitionSeconds, (end - start) / 2);
    const fades =
      fade > 0
        ? `,fade=t=in:st=${start}:d=${fade}:alpha=1,fade=t=out:st=${end - fade}:d=${fade}:alpha=1`
        : '';
    const motion = `scale=1080:1920:force_original_aspect_ratio=increase,crop=1080:1920,setsar=1,fps=60,format=rgba,tpad=stop_mode=clone:stop_duration=${reading.audioEnd}`;
    const window = `between(t,${start},${end})`;
    filters.push(
      `[${input++}:v]${motion},trim=duration=${reading.audioEnd},drawbox=x=35:y=95:w=885:h=230:color=0x071523@0.9:t=fill,drawtext=textfile='${label}':expansion=none:fontsize=32:fontcolor=white:x=60:y=120:line_spacing=10${fades}[picture${index}]`,
    );
    filters.push(
      `[v${index}][picture${index}]overlay=0:0:enable='${window}':eof_action=pass[illustrated${index}]`,
    );
    labels[index] = `[illustrated${index}]`;
  }

  return labels;
}
