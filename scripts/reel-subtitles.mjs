import { spawn } from 'node:child_process';
import { resolve } from 'node:path';

function timestamp(seconds) {
  const hundredths = Math.max(0, Math.round(seconds * 100));
  const hours = Math.floor(hundredths / 360000);
  const minutes = Math.floor(hundredths / 6000) % 60;
  const remaining = Math.floor(hundredths / 100) % 60;

  return `${hours}:${String(minutes).padStart(2, '0')}:${String(remaining).padStart(2, '0')}.${String(hundredths % 100).padStart(2, '0')}`;
}

function caption(text) {
  return text
    .replaceAll('\\', '\\\\')
    .replaceAll('{', '｛')
    .replaceAll('}', '｝')
    .replace(/\s+/g, ' ')
    .trim();
}

function silenceGaps(audio) {
  return new Promise((accept, reject) => {
    const child = spawn(
      'ffmpeg',
      ['-hide_banner', '-i', audio, '-af', 'silencedetect=noise=-38dB:d=0.18', '-f', 'null', '-'],
      { stdio: ['ignore', 'ignore', 'pipe'] },
    );
    let stderr = '';
    child.stderr.on('data', (data) => (stderr += data));
    child.on('error', reject);
    child.on('exit', (code) => {
      if (code !== 0) {
        reject(new Error(`Unable to analyze narration pauses: ${stderr.slice(-1000)}`));

        return;
      }

      const gaps = [];
      let start = 0;
      for (const match of stderr.matchAll(
        /silence_start:\s*([\d.]+)|silence_end:\s*([\d.]+)\s*\|\s*silence_duration:\s*([\d.]+)/g,
      )) {
        if (match[1]) {
          start = Number(match[1]);
        } else {
          gaps.push({
            start,
            end: Number(match[2]),
            duration: Number(match[3]),
          });
        }
      }

      accept(gaps);
    });
  });
}

function pauseBoundary(expected, after, gaps) {
  const candidates = gaps.filter(
    (gap) => gap.duration >= 0.2 && gap.end > after && Math.abs(gap.end - expected) <= 1.35,
  );
  candidates.sort((left, right) => Math.abs(left.end - expected) - Math.abs(right.end - expected));

  return candidates[0]?.end ?? expected;
}

export async function subtitleTrack(part, directory) {
  const audios = new Map();
  for (const passage of part.passages) {
    if (!audios.has(passage.audioId)) {
      const file = resolve(directory, `${passage.audioId}.mp3`);
      audios.set(passage.audioId, await silenceGaps(file));
    }
  }

  const lines = [
    '[Script Info]',
    'ScriptType: v4.00+',
    'PlayResX: 1080',
    'PlayResY: 1920',
    'WrapStyle: 0',
    '[V4+ Styles]',
    'Format: Name, Fontname, Fontsize, PrimaryColour, SecondaryColour, OutlineColour, BackColour, Bold, Italic, Underline, StrikeOut, ScaleX, ScaleY, Spacing, Angle, BorderStyle, Outline, Shadow, Alignment, MarginL, MarginR, MarginV, Encoding',
    'Style: Default,Liberation Sans,52,&H00F4F0E9,&H00F4F0E9,&H00191008,&H90000000,0,0,0,0,100,100,0,0,1,2,0,5,105,185,100,1',
    '[Events]',
    'Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text',
  ];
  let offset = 0;
  let audioID = '';
  let previousEnd = 0;

  for (const [index, passage] of part.passages.entries()) {
    if (passage.audioId !== audioID) {
      offset += previousEnd;
      previousEnd = 0;
      audioID = passage.audioId;
    }

    const next = part.passages[index + 1];
    const end =
      next?.audioId === audioID
        ? pauseBoundary(passage.audioEnd, previousEnd, audios.get(audioID))
        : passage.audioEnd;
    const start = timestamp(offset + previousEnd);
    lines.push(
      `Dialogue: 0,${start},${timestamp(offset + end)},Default,,0,0,0,,{\\pos(500,1490)}${caption(passage.text)}`,
    );
    previousEnd = end;
  }

  return lines.join('\n') + '\n';
}
