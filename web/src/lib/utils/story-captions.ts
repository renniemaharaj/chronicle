export function storyCaptions(text: string) {
  const words = text.match(/\S+\s*/g) ?? [];
  const captions: { text: string; start: number; end: number }[] = [];
  let start = 0;

  for (let index = 0; index < words.length; index++) {
    const count = index - start + 1;
    if (count < 12 && index < words.length - 1 && !(count >= 5 && /[.!?]["'’]?\s*$/.test(words[index]))) {
      continue;
    }

    captions.push({
      text: words
        .slice(start, index + 1)
        .join('')
        .trim(),
      start,
      end: index + 1,
    });
    start = index + 1;
  }

  return {
    captions,
    words: words.length,
  };
}

export function activeStoryCaption(text: string, seconds: number, duration = 0) {
  const track = storyCaptions(text);
  const length = duration > 0 && Number.isFinite(duration) ? duration : track.words * 0.4;
  const position = Math.min(track.words - 1, (Math.max(0, seconds) / Math.max(0.1, length)) * track.words);

  return track.captions.find((caption) => caption.start <= position && position < caption.end)?.text ?? '';
}
