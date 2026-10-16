# Stories, historical context and speech

Every saved event has a portable `narration` with readable `text`, `claimIds`, and an
`interpretation` flag. Existing projects receive evidence-qualified readings when opened by the backend.
Assistant research may supply richer stories, but must reference the claims supporting each assertion.
Rejected or unverified claims cannot support factual narration. Disputes remain visible and audible.
Changing the active interpretation changes the reading without changing the saved historical record.

`contextPeriods` stores scoped historical context with an uncertain temporal interval, summary and claim IDs.
Its windows describe the researcher's selected context, not authoritative boundaries for universal eras.
No era is invented when context is missing. These optional fields preserve portable format version 1.
JSON snapshots retain context; event JSONB attributes also retain narration in PostgreSQL.

## Playback

Play reads the story on screen and advances after speech finishes. Pause preserves the current position.
Previous/next or changing a voice cancels the old reading. Speed changes restart the current reading.
Follow chapter isolates exact event membership, including chapters with overlapping dates.
All history clears that filter and disables following. Chapter names come from each research project.
The silent reading mode advances at an estimated 150 words per minute, adjusted by playback speed.

Browser voices use the browser's Web Speech API. Voice availability and whether processing is local
are controlled by the operating system and browser. Google voices appear only if that browser supplies them.
If speech is unavailable, playback pauses with an error; it does not pretend that audio was generated.

## Optional Google Cloud Text-to-Speech

Choose **Voice & playback > Google Cloud Text-to-Speech**. This sends narration to Google.
The backend uses Application Default Credentials (ADC); no credential is stored in the frontend.
Configure an existing Google Cloud project with the Text-to-Speech API enabled and appropriate access.
With the Google Cloud CLI installed, establish local credentials and a quota project:

```sh
gcloud auth application-default login
gcloud auth application-default set-quota-project YOUR_PROJECT_ID
```

Alternatively set `GOOGLE_APPLICATION_CREDENTIALS` to an authorized credentials file before starting Chronicle.
Google Cloud service charges and quotas apply. Chronicle does not create projects or enable billing.
Click **Load Google voices**, select a voice, then Play. Credentials are resolved only when Cloud speech is requested.
The backend limits text to 5,000 UTF-8 bytes, bounds responses, times out requests and propagates cancellation.

Local endpoints: `GET /api/speech/voices?language=en-US` and `POST /api/speech/synthesize`.
The latter accepts `{ "text": "...", "voiceName": "...", "languageCode": "en-US" }` and returns MP3 audio.

References: [Google authentication](https://docs.cloud.google.com/text-to-speech/docs/authentication),
[voice discovery](https://docs.cloud.google.com/text-to-speech/docs/reference/rest/v1/voices/list),
[synthesis](https://docs.cloud.google.com/text-to-speech/docs/reference/rest/v1/text/synthesize),
[browser speech](https://developer.mozilla.org/en-US/docs/Web/API/SpeechSynthesis).
