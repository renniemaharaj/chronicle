# Setup and development

Run Chronicle locally, configure automatic research, and choose how projects are stored.

## Run

Requires Node.js 22+ and Go 1.26+.

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
Saved projects appear live in the app. See [assistant tools](assistant-tools.md).

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
generated dependency lockfiles are exempt from the authored-file line limit. Build output, dependencies, and runtime data are excluded.

See [architecture](architecture.md), [research behavior](research.md), and [asset credits](assets.md).

Story playback and optional Google Cloud voice setup are described in [Speech](speech.md).
