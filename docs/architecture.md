# Architecture

Chronicle runs one Go process and a static SvelteKit application. Vite proxies `/api` during development.
Production builds are served by Go from `web/build`.

## Canonical knowledge

`internal/domain` contains universal entities, events, relationships, claims, sources, evidence, conflicts,
coverage intervals, chapters, annotations, and research records. Extensible entity metadata lives in JSON attributes.
Domain validation checks references, precision, identifiers, provenance, and required fields before persistence.

The canonical model does not contain Svelte Flow nodes. `web/src/lib/graph` projects records into canvas nodes,
edges, phases, detail layers, and alternative interpretations. Date windows and a bounded event budget limit rendering.
Svelte Flow also renders only visible elements. Large project snapshots are currently loaded as whole documents;
server-side indexed viewport queries are a later performance step.

## Time

Historical dates use signed astronomical years: year 0 is 1 BCE. They never depend on JavaScript Date.
A value has explicit precision, calendar, optional approximation, and a numerical uncertainty range.
Intervals may have open ends. Gregorian and Julian calendars normalize deterministically to inclusive second bounds.
Unsupported historical calendars remain recorded but incomparable until an adapter is available.
Approximation alone never invents a numerical error margin.

Logical checks distinguish disjoint date intervals, physical presence outside a lifespan, suspicious parent chronology,
potentially incompatible locations, and unsupported or reversed causal ordering. Posthumous publication is permitted.
A check records a conflict only when at least two existing claims can be linked; it never manufactures supporting evidence.

## Persistence

A small storage interface supports owner-only atomic files and PostgreSQL.
PostgreSQL writes normalized entities, events, participants, sources, claims, evidence, relationships, and conflicts,
then commits the complete versioned workspace snapshot in the same transaction.
Canonical snapshots preserve narrative chapters, annotations, layouts, precision, alternatives, and audit records.
Normalized historical records are retained when omitted from a later project view.

## Research

The official MCP Go SDK handles discovery and Streamable HTTP. JSON Schema validates arguments before tool calls.
Provider adapters support Ollama and compatible chat-completions endpoints through a small model interface.
Research sessions checkpoint their plan, request count, model configuration, tool inputs, outputs, and partial project.
API keys are excluded from research audit records and returned settings.

The interface keeps presentation, state, projections, and API access in separate small modules.
A repository check enforces the user's limit of fewer than 150 lines per authored file.
