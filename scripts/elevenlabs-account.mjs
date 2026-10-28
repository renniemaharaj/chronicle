import { ElevenLabsClient } from '@elevenlabs/elevenlabs-js';

// Credentials stay in the server environment; this command never generates audio.
const client = new ElevenLabsClient();

try {
  const usage = await client.user.subscription.get();
  process.stdout.write(
    `${JSON.stringify(
      {
        tier: usage.tier,
        used: usage.characterCount,
        limit: usage.characterLimit,
        resetsAt: usage.nextCharacterCountResetUnix,
      },
      null,
      2,
    )}\n`,
  );
} catch {
  process.stderr.write('Unable to read ElevenLabs usage. Check the key permissions and connection.\n');
  process.exitCode = 1;
}
