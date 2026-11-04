# Timed passage images

In Story mode, expand **Passage image & timing** under the passage. Choose one image from the project's
media, set its start and on-screen duration in narration seconds, and set each fade duration. Save before
using **Preview saved image cue**. Choose **Map only** to remove the cue. The toolbar's image button can
hide all timed images without deleting their saved settings.

The duration includes both fades. The image returns to the map when its window ends or the passage
changes. Narration continues underneath. Saved audio uses its actual playback position; screen reading
and browser speech use a playback clock, with browser speech timing approximate. Pausing freezes the
cue. Loading saved audio does not advance it. Preview plays the saved visual window without generating
or playing speech. Image loading failures leave the map visible.

## Canonical records

Each event has an optional `visual` with `mediaId`, `startSeconds`, `durationSeconds` and
`transitionSeconds`. It references one canonical `media` item. Timings are finite, nonnegative seconds;
duration must be at least one second, each fade at most five seconds, and both fades must fit the window.
Images with definite temporal bounds outside the event's definite bounds are rejected. Unresolved or
historical calendars retain their existing uncertainty; the validator does not invent date conversions.

Existing source images retain captured URLs, attribution, source IDs and reuse terms. Generated images
carry `generated.tool`, the complete `prompt`, `createdAt`, historical `referenceSourceIds`, and a
`depicts` interval. Guidance sources contextualize a reconstruction; the image is never supporting evidence.
GIF animations play in Story mode and loop for the duration of a saved reel cue. Generated images cannot
become retrieved source images or verified avatars. Rights not recorded remain
unknown. Their caption and visible AI label distinguish reconstruction from an archival image.

## Importing an illustration

Generate the image with an image tool, then POST its PNG or JPEG bytes to `/api/media/images`.
The local backend accepts PNG, JPEG and GIF files up to 16 MB, 8192 pixels per side and 40 megapixels.
It returns a content-addressed
`url`, `width` and `height`. Images live in `WORKSPACE/media/images` and are served through a confined
read-only endpoint. Keep that directory with workspace backups; a project JSON does not embed the bytes.

Add the generated `media` record and the event cue through `chronicle_validate_project` followed by
`chronicle_save_project`. Existing claims, evidence and source captures must be preserved. The local
assistant tool bridge supports the same workflow after the backend is restarted with the new schema.

## Reels

New saved-audio campaigns freeze each selected image and cue with the passage. The FFmpeg renderer
inserts it over the recorded map, burns its AI/source label, fades back out and preserves narration and
captions. Image windows are clipped to the saved reading's duration. Changing later media or timings
does not change an existing take. The narration fingerprint excludes presentation-only cue changes.

The Rome project starts with three built-in image-tool illustrations: early Palatine huts, aqueduct
maintenance in the first to second centuries CE, and the Colosseum around 80 CE. Their original PNGs and
complete prompt set are in `.chronicle/illustrations/rome`; generation provenance is also saved on each
canonical media record. The map/image/map preview is in `.chronicle/exports/rome-image-preview`.
