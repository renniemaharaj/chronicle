# Places and travel

Choose **Map** in Explorer or the map icon above the canvas. The map shares the
project's chapters, date window, selected record and Story pane. **Places** shows
canonical place entities; selecting one exposes its timeline events and an
OpenStreetMap link. **Travel** shows evidence-backed events that explicitly
require the subject's physical presence. **Through [year]** follows the timeline.

## Geographic references

Each place may have `attributes.mapLocation` containing `latitude`, `longitude`,
`label`, `kind`, `sourceURL`, `provider` and `retrievedAt`. Coordinates must be
finite and within geographic bounds. These modern geographic references do not
assert historical boundaries or precise historical addresses. Country, region
and city markers represent the named area. They are portable with the project;
canvas state does not create a second collection of place entities.

**Locate** runs an explicit place search. Choose a matching result before its
coordinates are saved. There is no search autocomplete or automatic geocoding.
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
