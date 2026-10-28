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
Click **Load Google voices**, select a voice, then generate missing audio before Play. Credentials are resolved only when Cloud speech is requested.
The backend limits text to 5,000 UTF-8 bytes, bounds responses, times out requests and propagates cancellation.

Local endpoints: `GET /api/speech/voices?language=en-US` and `POST /api/speech/synthesize`.
The latter accepts `{ "text": "...", "voiceName": "...", "languageCode": "en-US" }` and returns MP3 audio.

References: [Google authentication](https://docs.cloud.google.com/text-to-speech/docs/authentication),
[voice discovery](https://docs.cloud.google.com/text-to-speech/docs/reference/rest/v1/voices/list),
[synthesis](https://docs.cloud.google.com/text-to-speech/docs/reference/rest/v1/text/synthesize),
[browser speech](https://developer.mozilla.org/en-US/docs/Web/API/SpeechSynthesis).

## ElevenLabs

Choose **Voice & playback > ElevenLabs**. Set `ELEVENLABS_API_KEY` in the server environment
or enter a key in the dialog. Environment/server configuration takes precedence over workspace keys.
The UI detects the key without receiving its value, and automatically loads the account's paginated
voice catalog, text-to-speech models and subscription usage. Restricted keys may allow speech while
blocking usage or voice discovery; those failures are shown independently. Refresh usage after generation.
Workspace keys are stored on the backend, outside portable projects and browser preferences.
An empty key removes a workspace credential. `CHRONICLE_ELEVENLABS_URL` overrides the API endpoint.

ElevenLabs audio uses 44.1 kHz, 128 kbps MP3, stability 0.5 and similarity 0.75.
The selected model is sent to the provider; model availability comes from its API, not screenshots.
Usage shows account-wide provider counters and reset dates. Missing text characters are a generation
estimate, not a credit quote: pricing and credit multipliers vary by model and subscription.

## Saved narration

For Google and ElevenLabs, select a voice and choose **Generate missing audio** before Play.
The dialog checks each story passage and shows saved parts, missing parts and missing text characters.
Long passages are split on Unicode-safe boundaries below 4,500 UTF-8 bytes. Generation saves each part
before continuing. Stop lets the current part finish; closing the dialog also stops after that part.
Completed parts survive interruption and provider failures. Retry generates only parts still missing.

Files live in `WORKSPACE/audio/` as private MP3 files, persisted across server restarts, including
PostgreSQL deployments (retain the workspace volume). Audio is not embedded in project exports.
The cache identity includes exact narration, provider endpoint, voice, language, model and the fixed
encoding/synthesis settings version. Playback speed is applied locally and does not regenerate audio.
Changing a voice keeps previous versions; switching back reuses them. Changed text needs new audio.
Concurrent generation within one server rechecks the cache to avoid duplicate provider charges.
Keep one Chronicle server per workspace when generating audio. Cache files remain until removed from
that directory; do not remove files while generation is active. Retain them to avoid regeneration charges.

Playback checks local saved audio and streams files with byte-range support. Missing files pause with
an instruction to generate them; Play never silently starts paid synthesis. Browser Web Speech cannot
export reusable audio through this API, so browser voices continue speaking on demand.

Local endpoints: `GET/PUT /api/speech/settings`, `GET /api/speech/voices?provider=elevenlabs`,
`GET /api/speech/models`, `GET /api/speech/usage`, `POST /api/speech/prepare` and
`GET /api/speech/audio/{id}`. Prepare accepts the synthesis fields plus `provider`, `modelId` and
`generate`; false only inspects the cache, true generates missing files and returns their local URL.
The legacy synthesis endpoint now also persists and reuses audio.

References: [ElevenLabs voices](https://elevenlabs.io/docs/api-reference/voices/search),
[models](https://elevenlabs.io/docs/api-reference/models/list),
[usage](https://elevenlabs.io/docs/api-reference/user/subscription/get),
[speech](https://elevenlabs.io/docs/api-reference/text-to-speech/convert).

## Preparation progress and updates

Library and **Story preparation status** show evidence-linked events and supported
narration text as percentages of saved events. They do not imply that all history
has been researched. Research runs show completed planned questions; final
validation must succeed before the run can reach 100 percent.

**Voice & playback** shows audio completion for the current story perspective,
voice and model, plus the names of passages still missing audio. Credit exhaustion,
key spending limits and temporary rate limits stop generation with actionable
messages. No automatic paid retry occurs. After renewal or a top-up, choose
**Generate missing audio** again. Reopening the dialog checks the disk cache,
including after a restart; previously completed parts are reused.

Choose **Update story**, then **Update current story** in Research to investigate
recent developments, new scholarship and remaining gaps in the existing project.
This requires the configured research model and MCP tools. Assistant-led updates
can instead use capture, validate and save. Refresh is on demand; there is no
background news polling. New events retain their dates, captured quotations and
competing claims. Announced future events remain plans. Unchanged narration reuses
saved audio; changed text needs regeneration. Concurrent story edits detected
before saving leave findings in the research session instead of replacing edits.

`npm run speech:account` uses the ElevenLabs SDK to inspect the reported allowance
without generating audio. Provider usage and reset timestamps are account data,
not a guarantee that the selected model or key limit can fund a complete story.

Reference: [ElevenLabs error codes](https://elevenlabs.io/docs/eleven-api/resources/errors).
