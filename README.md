# Chronicle

A local research workspace built with Go, Svelte 5, SvelteKit, and Svelte Flow.
Enter a subject, connect your model and MCP tools, then inspect an evidence-linked temporal graph.

## Run

Requires Node.js 22+ and Go 1.25+.

```sh
npm install
npm run dev
```

Open http://localhost:5173. The Go API listens on 127.0.0.1:8080.

For a single local process:

```sh
npm run build
npm start
```

Open http://localhost:8080. No account or hosted service is required.

## Assistant research

Chronicle exposes local MCP tools at `http://127.0.0.1:8080/api/mcp`. An assistant can capture sources,
validate a universal project, and save it directly; this path needs no separately configured model.
Saved projects appear live in the app. See [assistant tools](docs/assistant-tools.md).

## Research configuration

Open Settings and configure:

1. Ollama, or an OpenAI-compatible chat endpoint. For compatible endpoints, include `/v1` in the base URL when required.
2. An MCP server using Streamable HTTP.
3. Explicit search and document retrieval tools, selected after discovering their schemas.
4. A request budget and concurrency limit.

Enter a subject using New Research. The pipeline identifies it provisionally, plans chronological questions,
searches and retrieves documents, extracts schema-validated records, matches quotations, assesses evidence,
performs a targeted follow-up when the budget allows, checks chronology, and assembles evidence-linked chapters.

Tool descriptions and retrieved material are untrusted data. Only the two configured tools may be invoked.
The current adapter requires sources returned as HTTP(S) URLs and textual document content.
MCP stdio, PDF OCR, and archive-specific protocols need additional adapters.

## Workspace

The default backend stores atomic JSON checkpoints in `.chronicle/`, with owner-only file permissions.
Set `CHRONICLE_WORKSPACE` to choose another directory. The browser also caches projects for offline exploration.
Browser storage limits can be reached by large research audit logs; backend persistence and exports remain available.

Set `DATABASE_URL` to a PostgreSQL connection string to use PostgreSQL instead.
Migrations run transactionally on startup, with version tracking and an advisory lock.
The database preserves canonical project snapshots and normalized ontology records with foreign keys.

`CHRONICLE_ADDR` changes the listener address. The default listener is local to the workstation.

## Explore

- Select nodes, inspect claims and passages, and follow source links.
- Change temporal resolution independently of spatial zoom. Broad views aggregate events; detailed views expose claims and evidence.
- Filter date ranges, jump to a year, choose quarters of the subject's existence, and fit the visible period.
- Compare shadow hypotheses, parallel lanes, isolated conflicts, and narrative consequences.
- Select an active interpretation to change the projection while preserving canonical records.
- Play a deterministic sequence of events, select chapters, and change perspective or speed.
- Resize or collapse the explorer and inspector, switch themes, and annotate records.
- Save, reopen, import, or export projects. PDF uses the browser's print-to-PDF workflow.

Keyboard: Ctrl/Cmd+K opens research, Ctrl/Cmd+S saves, arrows move between milestones, Space toggles playback.

## Evidence and fixtures

Workspaces start empty. Synthetic fixtures for a person, family, institution, and place live only in
development tests; none are seeded into storage or bundled into the interface. Older demo projects are
hidden from the library without deleting their files.

A matching passage establishes that text was retrieved; it does not prove an assertion true.
Corroboration requires supporting passages from two explicitly reviewed independent source groups.
Contradicting evidence remains available. Missing passages leave claims unverified.
Narrative synthesis currently groups supported assertions into chronological chapters; scholarly interpretation requires review.

## Checks

```sh
npm run check
npm test
npm run build
npm run test:browser
```

Browser tests use an installed Chromium at `/usr/bin/chromium`, or Playwright's installed browser.
The research tests use a real MCP transport with synthetic archive responses and a deterministic model.
Live provider quality and PostgreSQL integration should be verified separately in the intended local environment.

Every authored text file contains fewer than 150 lines. `scripts/check-lines.mjs` enforces this;
generated dependency locks retain all their data in compact JSON. Build output, dependencies, and runtime data are excluded.

See [architecture](docs/architecture.md), [research behavior](docs/research.md), and [asset credits](docs/assets.md).
