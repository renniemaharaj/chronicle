# Social media reach campaigns

Open a saved research project and choose **Campaigns** in the top bar. Research still uses the existing
assistant tools or configured autonomous research pipeline; campaign creation does not run a new LLM.
This first release provides series planning, vertical narrated takes, and a publication/performance ledger.
It does not autonomously sign into Facebook, upload files, collect insights, or promise monetization.

## Series workspace

Reel Series appears in the research explorer. The Series tab combines a production summary, searchable
episode table, coverage overview, and ordered distribution history per platform. Select an episode to inspect
its reading, linked canonical sources, latest observations and continuity boundaries. The inspector can be
closed or resized; on narrow screens it opens as a drawer.

Clips Library previews exported MP4 takes and lists their local output paths. Analytics uses the latest saved observation
for each metric and episode, preserves missing values, and groups revenue by currency. Growth compares
two distinct observation times for the same episode. These figures are manual observations, not live data.

Open an episode's studio to preview narration, record with OBS, review takes, or manage its publication.
Preview series reads across episode boundaries without changing take or publication checkpoints. Recording
controls stay outside the vertical capture frame. Artwork in the workspace comes from linked research media;
when absent, Chronicle displays a neutral graphic. Research thumbnails are not recorded video previews.

## Continuity

Each campaign references a canonical project and freezes its assessed event narration, historical context,
claim identifiers, source references/hashes, historical intervals, and editorial qualifications. Context is
presented first; events follow recorded historical dates, with unknown dates last. This is an editorial
ordering, not a conversion between historical calendars. Unverified narration becomes an explicit evidence
gap. Research coverage gaps remain visible; the planner cannot establish that an investigation is complete.

The planner targets 140 words per minute and divides narration into captions of at most 24 words, preserving
every UTF-8 byte in each qualified passage. Start/end offsets refer to that frozen passage, not the source
quotation or the original curated narration. Reel boundaries prefer sentence endings and may fall within an event. A long sentence can exceed the target. Joining a series'
passages reconstructs its entire qualified reading without overlap. Durations are estimates, not platform limits.

Take checkpoints follow `planned -> recording -> recorded -> published`. Each take has its own identifier.
An unfinished take retries its entire part; it never advances past a partly recorded passage. On restart,
inspect an open take, stop OBS if necessary, and either attach a complete file or discard the checkpoint.
The application does not verify a manually attached file's contents or duration.
Completed, unpublished parts can be retaken. Earlier file paths and measured durations remain in take
history; posted parts cannot be reset. Retaking does not delete any OBS file.

All actions require the current revision. One backend process serializes checkpoint changes and permits
only one open recording across campaigns. This assumes a single Chronicle backend; multiple backend
replicas require database transactions/leases before enabling concurrent campaign workers.

Research fingerprints prevent starting more recordings after relevant canonical research changes. Existing
recordings and posted history remain intact. Review what was already covered before creating a new series;
cross-series deduplication and incremental append/reconciliation are not yet implemented.

## Recording with OBS

1. Enable OBS WebSocket 5 on localhost port 4455 and set its password.
2. Create a scene named `Chronicle Reels` (or enter your chosen name in the studio).
3. Set output to 1080 x 1920 or another exact 9:16 resolution. Capture the vertical preview region,
   crop away application controls, and include the browser's narration audio. Verify both in a test take.
4. Choose an audio voice through the existing Voice & playback settings. Screen-only reading is rejected.
5. Choose a part and press **Record part with OBS**. The studio checks output geometry, active scene,
   and idle recording status, starts OBS, reads the part, then stops and saves its file/duration checkpoint.

The studio does not configure your capture source, crop, audio routing, encoding or scene transitions.
It uses a text composition without research images, whose capture/reuse rights need separate review.
Keep this browser tab active. Narration or OBS failures leave a recovery path; never treat a partial file
as a finished reel. A saved file does not establish publication. Capture uses the existing speech providers;
cloud voice setup, generation latency and usage charges follow their existing configuration.

OBS integration follows the [OBS WebSocket protocol](https://github.com/obsproject/obs-websocket/blob/master/docs/generated/protocol.md)
through [obs-websocket-js](https://github.com/obs-websocket-community-projects/obs-websocket-js).
The OBS password stays in component memory and is not written to the backend or browser storage.

## Publication and performance

Saved-audio campaigns can target Facebook Reels and YouTube Shorts. Publication receipts and posting
order are independent for each platform. An episode becomes published when every selected target has
a receipt. A partially published episode cannot be retaken. Destination identities can be recorded with
the `destination` action after observing the signed-in account; accounts with posts cannot be replaced.
Existing campaigns and receipts retain their Facebook meaning. The analytics summary currently shows
Facebook observations; the publication ledger shows observations for its selected platform.

The Facebook composer link opens an external tab. Sign in there when necessary, select the intended Page
or profile, upload the recorded file and confirm the actual publication. Chronicle cannot inspect that tab's
login state. Paste the direct numeric `/reel/ID` or `/videos/ID` URL to register a manual publication receipt.
Only the next unpublished part accepts a receipt. Reusing a receipt within a series is rejected. This ledger
does not independently verify the URL or prevent a person from uploading a duplicate on Facebook.

Record observed views, reactions, comments, shares and revenue as timestamped manual observations.
Blank values mean unknown; zero means an observed zero. Revenue requires a currency code. These are
observations, not automatically fetched platform analytics or evidence of monetization eligibility.

The assistant's `replace` action records a confirmed replacement with `previousUrl`, the new `url`,
`file`, measured `duration`, and `metrics` refreshed within five minutes. It requires fewer than 10
observed views and protects any post previously observed at 10 or more. Old receipts, observations
and take files remain in replacement history; other platform receipts retain their own take files.
This action records completed external publishing; it does not remove or upload a post itself.

## Agent tools and next integration

`POST /api/campaigns/from-audio` and `chronicle_create_audio_campaign` inspect the selected voice's exact
chronological readings, stop at the first missing cache file, and measure existing audio with `ffprobe`.
They never generate paid speech. `scripts/render-reels.mjs CAMPAIGN_ID` records the actual Story map
with Follow location enabled through Playwright, holds each map view for its saved narration, and uses
FFmpeg to assemble 1080 x 1920 episodes. Caption timing starts from word-count estimates and moves
transitions to detected narration pauses when available.
Complete audio passages stay together; one long passage can exceed the target episode duration.
Map capture uses 60 fps, and exports use blended 120 fps frames with sharp captions.

Exports live under `WORKSPACE/exports/CAMPAIGN_ID`, with source descriptions and take checkpoints.
The video endpoint streams only registered MP4 files inside that export root, including byte ranges.
Local OBS files outside the export root remain local. `scripts/playwright-tool.mjs` connects to a local
Playwright MCP bridge for signed-in browser operations. Browser login and extension file permissions
remain user-controlled; credentials are not written into project records or repository configuration.

`scripts/publish-facebook-reels.mjs CAMPAIGN_ID [MAX_EPISODES]` and
`scripts/publish-youtube-shorts.mjs CAMPAIGN_ID [MAX_EPISODES] [PLAYLIST_TITLE]` drive an already connected
browser session, with Facebook in tab 0 and YouTube Studio in tab 1. Run them sequentially: they share the
active browser tab. Both verify the destination and save confirmed public links through the backend.
The local `social/posting-progress.json` records submissions before publishing. Facebook reconciles pending
submissions against the profile's reels; YouTube stops at an existing draft or submission for reconciliation.
These UI helpers depend on the observed platform layout and are not background workers or API adapters.
`scripts/publish-reel-series.mjs CAMPAIGN_ID [MAX_EPISODES]` resumes both targets sequentially, starting
with the target that has fewer confirmed parts. A paused target does not prevent trying the other target.
YouTube daily-limit rejections without an assigned video link keep a conservative 24-hour retry boundary.
A later run may retry that rejected upload after the boundary; drafts and submitted videos still require
reconciliation. The helper does not schedule itself or bypass channel-history and quota requirements.

The HTTP API exposes `GET/POST /api/campaigns`, `GET /api/campaigns/{id}`, and
`POST /api/campaigns/{id}/actions`. Equivalent MCP tools are `chronicle_list_campaigns`,
`chronicle_get_campaign`, `chronicle_create_campaign`, and `chronicle_update_campaign`.
Use `node scripts/chronicle-tool.mjs list` to discover them after restarting the backend.
Agents must inspect existing progress and preserve actual publication receipts, never infer success from
opening the composer. The current MCP tools manage checkpoints; they do not drive the studio or OBS.

Next steps are a durable runner for research, capture and publication jobs, word-level caption alignment,
append-only series reconciliation after research updates, and platform publishing adapters.
The adapter needs an explicit Page/profile destination, session/account checks, a login-required pause,
upload processing and public-receipt reconciliation before retry, plus attributable insights snapshots.
Scheduling, retries, target accounts, publication approvals and permission/token lifecycle must be durable.
Consult Meta's current [Reels publishing](https://developers.facebook.com/docs/video-api/guides/reels-publishing/)
and [video insights](https://developers.facebook.com/docs/graph-api/reference/video/video_insights/) contracts
when implementing that adapter; these external contracts and available metrics can change.
