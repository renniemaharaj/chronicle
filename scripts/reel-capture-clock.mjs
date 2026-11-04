export async function startCaptureClock(page, reading) {
  await page.evaluate(
    ({ eventId, duration }) => {
      cancelAnimationFrame(window.chronicleCaptureFrame);

      const started = performance.now();

      const update = () => {
        const seconds = Math.min(duration, (performance.now() - started) / 1000);
        window.dispatchEvent(
          new CustomEvent('chronicle:story-capture-clock', {
            detail: { eventId, duration, seconds },
          }),
        );
        if (seconds < duration) {
          window.chronicleCaptureFrame = requestAnimationFrame(update);
        }
      };

      update();
    },
    {
      eventId: reading.eventId,
      duration: reading.audioEnd,
    },
  );
}

export async function stopCaptureClock(page) {
  await page.evaluate(() => cancelAnimationFrame(window.chronicleCaptureFrame));
}
