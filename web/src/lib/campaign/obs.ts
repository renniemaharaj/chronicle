import OBSWebSocket from 'obs-websocket-js';

export class Recorder {
  private client = new OBSWebSocket();
  private ownsRecording = false;
  private uncertainStart = false;

  async connect(password: string, sceneName: string) {
    await this.client.connect('ws://127.0.0.1:4455', password, { rpcVersion: 1 });

    const video = await this.client.call('GetVideoSettings');
    if (video.outputWidth * 16 !== video.outputHeight * 9) {
      throw new Error('Set OBS output to 9:16, for example 1080 x 1920, before recording.');
    }

    const scene = await this.client.call('GetCurrentProgramScene');
    if (!sceneName.trim() || scene.currentProgramSceneName !== sceneName.trim()) {
      throw new Error('Select the configured Chronicle capture scene in OBS before recording.');
    }

    const status = await this.client.call('GetRecordStatus');
    if (status.outputActive) {
      throw new Error('OBS is already recording. Finish or recover that take first.');
    }
  }

  async start() {
    this.uncertainStart = true;
    await this.client.call('StartRecord');
    this.ownsRecording = true;
    this.uncertainStart = false;
  }

  async stop() {
    if (!this.ownsRecording) {
      throw new Error('This studio does not own an active OBS recording.');
    }

    const status = await this.client.call('GetRecordStatus');
    if (!status.outputActive) {
      this.ownsRecording = false;
      throw new Error('OBS stopped unexpectedly. Recover the take before retrying.');
    }

    const output = await this.client.call('StopRecord');
    this.ownsRecording = false;

    return {
      file: output.outputPath,
      duration: status.outputDuration / 1000,
    };
  }

  async discard() {
    if (this.uncertainStart) {
      throw new Error('OBS did not acknowledge the start request. Inspect OBS and recover the saved take.');
    }

    if (this.ownsRecording) {
      await this.stop();
    }
  }

  async disconnect() {
    await this.client.disconnect();
  }
}
