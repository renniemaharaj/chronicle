# Places and travel

Choose **Map** in Explorer or the map icon above the canvas. The map shares the
project's chapters, date window, selected record and Story pane. **Places** shows
canonical place entities; selecting one exposes its timeline events and an
OpenStreetMap link. **Travel** shows evidence-backed events that explicitly
require the subject's physical presence. **Through [year]** follows the timeline.

## Story mode

Choose **Story** in Explorer or **Story mode** above the Story pane. Narration
plays over a map that follows the current event's canonical `locationIds`.
Multiple recorded locations are framed together. Repeated locations keep the
camera steady, and reduced motion preferences disable animated map movement.
**Follow location** can be switched off to explore the map manually.

Passages without mapped locations show the chapter's wider mapped context
(or the story's context when no chapter applies). The location caption explains
missing references. Research supplies source-backed coordinates; playback does not
infer them from narration. The base map does not depict period borders.

Use Play, Pause, previous/next, the passage slider or the chapter selector to
navigate. Existing voice, speed and evidence controls remain available.
**Exit story mode** or Escape returns to the previous view at the current event.
The fullscreen button expands the player; Escape first leaves fullscreen.

## Geographic references

Each place may have `attributes.mapLocation` containing `latitude`, `longitude`,
`label`, `kind`, `sourceURL`, `sourceId`, `claimIds`, `provider` and `retrievedAt`.
Optional `note` and `precisionMeters` retain source limitations. Coordinates must be
finite and within geographic bounds. These modern geographic references do not
assert historical boundaries or precise historical addresses. Country, region
and city markers represent the named area. They are portable with the project;
canvas state does not create a second collection of place entities.

Assistant research resolves these references from captured geographic sources,
with entity claims and verified coordinate passages. Existing projects can be
enriched through the same validate/save tools. A full history rerun is unnecessary.
Regions and natural features use representative points and broad camera framing.
Coordinate sources are accessible in Story mode and beside each map place.

**Locate** is an optional manual correction tool. It runs an explicit place search;
the selected result supplies an OpenStreetMap reference. There is no search autocomplete.
The local backend serializes requests with at least one second between starts
and caches up to 256 distinct queries. No API key is required.

The default search provider is [Photon](https://photon.komoot.io/), which permits
reasonable project use and may throttle or become unavailable. Set
`CHRONICLE_GEOCODER_URL` to another Photon-compatible `/api/` endpoint to change
providers without changing code. Map records retain OpenStreetMap attribution.

The map uses Leaflet and browser-cached OpenStreetMap tiles with visible
attribution. Only the current viewport loads tiles; no offline download or
prefetch is provided. Follow the [tile usage policy](https://operations.osmfoundation.org/policies/tiles/).
`VITE_MAP_TILE_URL` can select a compatible OSM-derived raster tile service at
build time; retain attribution and comply with that provider's terms.

## Travel limits

Travel requires `requiresPresence: true`, subject participation, a located event,
and a supported claim with verified supporting evidence. Posthumous events are
excluded using the subject's recorded ending. Unknown or unsupported calendars
are omitted from travel ordering; historical date labels stay unchanged.

Dashed connections indicate a sequence of recorded visits, never a road, flight
or inferred itinerary. Connections are omitted between overlapping date ranges,
ambiguous multiple locations and unmapped stops. Broad dates remain broad;
overlap alone does not establish simultaneous presence in different places.
The event list retains dated visits even when a place has not been located.
