// Split on words where possible, respecting Google's UTF-8 byte limit.
export function speechChunks(text: string, maxBytes = 4500): string[] {
  const encoder = new TextEncoder();
  const chunks: string[] = [];
  let current = '';
  let size = 0;
  for (const character of text) {
    const bytes = encoder.encode(character).length;
    if (size + bytes > maxBytes) {
      const boundary = current.lastIndexOf(' ');
      if (boundary > current.length / 2) {
        chunks.push(current.slice(0, boundary));
        current = current.slice(boundary + 1);
        size = encoder.encode(current).length;
      } else {
        chunks.push(current);
        current = '';
        size = 0;
      }
    }
    current += character;
    size += bytes;
  }
  if (current.trim()) chunks.push(current);
  return chunks;
}
