import { spawn } from 'node:child_process';
import { writeFile } from 'node:fs/promises';
import { resolve } from 'node:path';
import { backend } from './reel-backend.mjs';
import { subtitleTrack } from './reel-subtitles.mjs';
import { reelImageOverlays } from './reel-image-overlays.mjs';

export function mediaCommand(command, args) {
  return new Promise((accept, reject) => {
    const child = spawn(command, args, { stdio: ['ignore', 'pipe', 'pipe'] });
    let stdout = '';
    let stderr = '';
    child.stdout.on('data', (data) => (stdout += data));
    child.stderr.on('data', (data) => (stderr = (stderr + data).slice(-6000)));
    child.on('error', reject);
    child.on('exit', (code) =>
      code === 0 ? accept(stdout) : reject(new Error(`${command} failed (${code}): ${stderr}`)),
    );
  });
}

export async function videoMetadata(path) {
  const result = await mediaCommand('ffprobe', [
    '-v',
    'error',
    '-show_entries',
    'format=duration:stream=codec_type,width,height,avg_frame_rate',
    '-of',
    'json',
    path,
  ]);
  const metadata = JSON.parse(result);
  const video = metadata.streams.find((stream) => stream.codec_type === 'video');
  if (
    video?.width !== 1080 ||
    video.height !== 1920 ||
    video.avg_frame_rate !== '120/1' ||
    !metadata.streams.some((stream) => stream.codec_type === 'audio')
  ) {
    throw new Error('The exported clip needs a 1080 x 1920 video stream and narration audio.');
  }

  return Number(metadata.format.duration);
}

export async function encodeEpisode(part, directory, output, mapDirectory = resolve(directory, 'map')) {
  const readings = [...new Map(part.passages.map((passage) => [passage.audioId, passage])).values()];
  const args = ['-hide_banner', '-loglevel', 'error', '-y'];
  const filters = [];
  const total = readings.reduce((sum, reading) => sum + reading.audioEnd, 0);
  for (const [index, reading] of readings.entries()) {
    const audio = resolve(directory, `${reading.audioId}.mp3`);
    const response = await fetch(`${backend}/api/speech/audio/${reading.audioId}`);
    if (!response.ok) {
      throw new Error('Saved narration disappeared before export.');
    }

    await writeFile(audio, Buffer.from(await response.arrayBuffer()));
    args.push('-i', resolve(mapDirectory, `${reading.audioId}.webm`), '-i', audio);
    filters.push(
      `[${index * 2}:v]fps=60,scale=1080:1920,setsar=1,tpad=stop_mode=clone:stop_duration=${reading.audioEnd},trim=duration=${reading.audioEnd},setpts=PTS-STARTPTS[v${index}]`,
    );
    filters.push(`[${index * 2 + 1}:a]aresample=44100,asetpts=PTS-STARTPTS[a${index}]`);
  }

  const subtitles = resolve(directory, `episode-${part.number}.ass`);
  await writeFile(subtitles, await subtitleTrack(part, directory));

  const videoInputs = (await reelImageOverlays(part, readings, directory, args, filters)).join('');
  const audioInputs = readings.map((_, index) => `[a${index}]`).join('');
  filters.push(`${videoInputs}concat=n=${readings.length}:v=1:a=0[base]`);
  filters.push(`${audioInputs}concat=n=${readings.length}:v=0:a=1[audio]`);
  filters.push(
    `[base]minterpolate=fps=120:mi_mode=blend,drawbox=x=48:y=1320:w=890:h=340:color=0x071523@0.88:t=fill,ass='${subtitles.replaceAll("'", "'\\''")}'[video]`,
  );
  args.push(
    '-filter_complex_threads',
    '4',
    '-filter_complex',
    filters.join(';'),
    '-map',
    '[video]',
    '-map',
    '[audio]',
    '-c:v',
    'libx264',
    '-preset',
    'veryfast',
    '-crf',
    '22',
    '-threads',
    '4',
    '-c:a',
    'aac',
    '-b:a',
    '192k',
    '-movflags',
    '+faststart',
    '-shortest',
    output,
  );
  await mediaCommand('ffmpeg', args);

  const duration = await videoMetadata(output);
  if (Math.abs(duration - total) > 0.6 || duration > 180) {
    throw new Error('The export duration does not match its saved narration.');
  }

  return duration;
}
