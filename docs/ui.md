# Interface components

Chronicle uses Flowbite Svelte with Tailwind CSS 4. The Vite plugin compiles
utilities; `web/src/lib/styles/flowbite.css` registers the Flowbite plugin,
component sources, primary palette, and class-based dark mode.

`ThemeProvider` in the root layout applies the shared configuration from
`web/src/lib/components/ui/theme.ts`. The workspace controls the `dark` class
and Chronicle surface variables for both themes.

## Reuse

- Use Flowbite `Input`, `Select`, `Label`, `Checkbox`, and `Textarea` for forms.
- Use the shared `ui/Button.svelte` for actions. Its `variant` is `primary`,
  `secondary`, `ghost`, or `danger`; it composes Flowbite's Button.
- Add dialog contents through `DialogHost`, which provides one Flowbite Modal.
- Use Flowbite Tabs, Toggle, Card, and Popover for their respective interactions.
- Use `CollectionSearch` for Library, subjects, and media searches.
- Use `MediaPreview` for image and video previews in galleries and Story.
- Use `ResizeHandle` for all pane separators; `resize-handles.css` owns their styling.

Keep Chronicle's graph, historical date controls, and pane sizing in focused
components. These represent application behavior rather than generic controls.
Keep every authored file below 150 lines; generated dependency locks are exempt.

## Media and chronology

Media belongs to canonical records. `utils/media.ts` matches an asset's claim
IDs to an event's claims or to the referenced claims' `eventId`. A shared
subject or year alone does not establish that a clip records a performance.
YouTube previews use a validated video ID with an image-error fallback. Links
open the source in a new tab; thumbnails do not load an embedded video player.

See the [Flowbite Svelte quickstart](https://flowbite-svelte.com/docs/pages/quickstart)
and [theme provider](https://flowbite-svelte.com/docs/pages/theme-provider).
See [Places and travel](maps.md) for the Leaflet canvas, geographic references,
provider configuration and limits of travel reconstruction.
