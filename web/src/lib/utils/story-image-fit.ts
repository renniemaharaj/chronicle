export function floatingImageSize(
  width: number,
  height: number,
  ratio: number,
  readerWidth = 0,
  captionHeight = 70,
  matchReaderWidth = false,
) {
  const aspect = Number.isFinite(ratio) && ratio > 0 ? ratio : 9 / 16;
  const availableWidth = Math.max(1, width - readerWidth - 42);
  const originalWidth = Math.min(320, availableWidth, Math.max(1, height - 160) * aspect);
  const maximumHeight = Math.max(1, height - 78 - captionHeight);

  if (matchReaderWidth) {
    const imageWidth = Math.min(readerWidth, availableWidth);

    return {
      width: imageWidth,
      height: Math.min(Math.max(1, imageWidth - 2) / aspect, Math.max(1, height - 88 - captionHeight)),
    };
  }

  const imageWidth = Math.min(originalWidth * 1.1, availableWidth, maximumHeight * aspect);

  return {
    width: imageWidth,
    height: imageWidth / aspect,
  };
}
