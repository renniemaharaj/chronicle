# Historical research standards

Chronicle investigates a subject's documented existence and its wider historical setting.
A finished run is a bounded investigation, never a declaration that all available history has been found.

## Six research questions

The backend combines model-proposed leads with six deterministic areas of investigation:

1. Origins and setting: earliest attestation, precursors, birth or formation, geography, material conditions,
   communication, transport, technologies, social institutions and competing accounts.
2. Development: education or adaptation, work, structural changes, maturity, setbacks and less documented intervals.
3. Relationships: family or constituent entities, collaborators, rivals, institutions, ownership, places and influence.
4. Turning points: accomplishments, works, failures, decisions, conflicts, causes and consequences.
5. Endings and aftermath: latest documented state, death or dissolution if established, and effects on surviving subjects.
6. Public record: official references, public accounts, attributable images, external video links and further source leads.

For an ongoing subject, the ending question investigates its latest state. A person's biography must not
supply a fictitious death. An institution or object must not be forced into a human childhood or lifespan.

## Historical context

Research the conditions that preceded the subject, including technologies that were actually available in
that place and period. Date their introduction or adoption separately. Availability does not prove use,
inspiration or causal influence. An era name provides orientation, not an explanation.

Scoped `contextPeriods` carry claims, evidence and uncertain intervals. Editorial context windows are labeled
as interpretation. Do not assert that a subject experienced an event simply because it happened nearby.

## Geographic references

Research geographic references for canonical event places as part of the investigation. Capture coordinate
sources, attach exact quotations to entity claims, and save their IDs with `attributes.mapLocation`.
Do not ask the user to locate ordinary researched places. Prefer scholarly gazetteers and authoritative
geographic records; accessible encyclopedias can supply attributed coordinate references.
Use broad regional framing for regions, provinces, rivers and seas. A representative point does not establish
an event site, a journey, historical borders or continuous presence. Preserve ambiguous identifications as gaps.

## Endings and aftermath

Separate the date, place, reported immediate cause, underlying cause, source attribution and uncertainty.
Public reporting is not a new medical diagnosis. Respect living people's privacy and avoid unsupported
claims about symptoms, private grief, motivations or recovery.

Record mourning, commemorations, succession, institutional changes, public reactions and affected places
as separate events. Link demonstrated causes with `causeIds`; use ordinary relationships for association.
Posthumous events may involve a deceased person without implying physical presence.
Keep those events distinct from the subject's documented lifetime.

## Sources

Seek original records, scholarly publications, books, archives, reputable reporting, wikis, public forums
and public social posts where accessible. Classify each source and record its dependencies.
A forum post is an attributed account or research lead; repetition does not make it independent corroboration.
A book catalogue or search snippet does not establish access to the book's full text.
Use configured MCP document/archive tools for formats the built-in HTML/text capture cannot read.
Record blocked, missing, truncated or unavailable material and coverage gaps explicitly.

## Images and video

Capture preserves observed HTML/Markdown image, video and link resources with the source hash.
The resource budget prioritizes official references and media even on long pages.
Recognized YouTube player IDs may be normalized to external watch URLs.
Media is linked, not downloaded or hosted. Images retain caption, creator, license and source URL.
Unknown rights remain unknown. A logo, nearby image or filename alone must not become someone's portrait.
An avatar needs identifying entity claims with supporting captured evidence and an observed image URL.
Source-caption matching does not authenticate the underlying photograph.

## Official links

A website needs an explicit official reference in another captured source, supported by an entity claim.
A social account then needs an observed link from that established website, for the same entity.
The backend recomputes verification, rejecting caller-supplied flags and invented resource lists.
This confirms the reference chain, not a paid platform badge, current account activity or endorsement.
Record the method, source and retrieval date. Unverified candidates remain in portable records;
only confirmed references appear as official links in the interface.

## Reproducibility and limits

Research keeps tool activity, model settings, request budgets, source hashes, quotations, assessments,
competing interpretations, scoped context and coverage gaps. Assistant-led work uses the same capture,
validation and save tools; never edit storage files to populate research.
The default 40-request budget supports six questions and targeted review. Source access and model
configuration still constrain depth. Continue research when those limits leave substantive gaps.
