export interface MotionSettings {
  pulse: boolean;
  camera: boolean;
}

// Interpolate between media-time reports without advancing through loading stalls.
export function motionLoop(draw: (seconds: number, settings: MotionSettings) => void) {
  let frame: number | undefined;
  let seconds = 0;
  let speed = 1;
  let sampledAt = 0;
  let settings: MotionSettings = {
    pulse: false,
    camera: false,
  };

  const render = (now: number) => {
    const betweenSamples = Math.min(0.1, Math.max(0, (now - sampledAt) / 1000));
    draw(seconds + betweenSamples * speed, settings);
    frame = now - sampledAt < 100 ? requestAnimationFrame(render) : undefined;
  };

  const stop = () => {
    if (frame !== undefined) {
      cancelAnimationFrame(frame);
      frame = undefined;
    }
  };

  return {
    update(time: number, rate: number, running: boolean, effects: MotionSettings) {
      seconds = Number.isFinite(time) ? Math.max(0, time) : 0;
      speed = Number.isFinite(rate) && rate > 0 ? rate : 1;
      sampledAt = performance.now();
      settings = effects;
      if (running && seconds > 0 && (settings.pulse || settings.camera)) {
        if (frame === undefined) {
          frame = requestAnimationFrame(render);
        }
      } else {
        stop();
        draw(seconds, settings);
      }
    },
    stop,
  };
}
