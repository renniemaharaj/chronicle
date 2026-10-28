# Configuration

The optional JSON file is selected by `CHRONICLE_CONFIG`. Paths are relative
to the process working directory unless absolute. Unknown fields, invalid
endpoints and invalid limits fail startup.

```sh
cp config/example.json config/local.json
CHRONICLE_CONFIG=config/local.json npm start
```

Native defaults are shown in `config/example.json`; container defaults are in
`config/docker.json`. Remove research fields you want managed through Settings.

## Precedence

Research uses defaults, saved Settings, explicit JSON research fields, then
environment variables. Other server settings use defaults, JSON, then environment.
Explicit empty values also override prior values. Files are read at startup.
Deployment credentials are redacted in API responses and are not copied into
saved provider settings when the Settings form is saved.

Docker Compose loads `.env`. `npm start` and `npm run dev` also load the optional
gitignored `.env.local` created by ElevenLabs onboarding; exported variables win.
Direct Go commands use exported variables. Vite reads its development
variables from the root `.env` and shell environment.

## Server and maps

| JSON field | Environment variable | Native default |
| --- | --- | --- |
| `address` | `CHRONICLE_ADDR` | `127.0.0.1:8080` |
| `workspace` | `CHRONICLE_WORKSPACE` | `.chronicle` |
| `staticDir` | `CHRONICLE_STATIC_DIR` | `web/build` |
| `databaseUrl` | `DATABASE_URL` | Empty selects files |
| `allowedOrigins` | `CHRONICLE_ALLOWED_ORIGINS` | Local Vite origins on 5173/5174 |
| `geocoderUrl` | `CHRONICLE_GEOCODER_URL` | `https://photon.komoot.io/api/` |
| `mapTileUrl` | `CHRONICLE_MAP_TILE_URL` | OpenStreetMap raster tiles |
| `mapAttribution` | `CHRONICLE_MAP_ATTRIBUTION` | OpenStreetMap contributor credits |

Origins use a comma-separated environment value or a JSON array. Each origin
includes scheme and port. Empty removes extra origins; same-host requests and
requests without Origin remain accepted. Container defaults have no extra origins.
Chronicle remains a private workspace without user authentication.

The browser fetches map settings from `/api/config`; no image rebuild is needed.
Reload the browser after changing them. Attribution is trusted administrator HTML;
retain provider credits. `VITE_MAP_TILE_URL` remains a fallback when the backend
is unavailable. The public endpoint never exposes service credentials or database URLs.

## Research services

These fields are inside the JSON `research` object.

| JSON field | Environment variable | Default |
| --- | --- | --- |
| `provider` | `CHRONICLE_LLM_PROVIDER` | `ollama` |
| `baseUrl` | `CHRONICLE_LLM_URL` | `http://localhost:11434` |
| `model` | `CHRONICLE_LLM_MODEL` | Empty |
| `apiKey` | `CHRONICLE_LLM_API_KEY` | Empty |
| `mcpUrl` | `CHRONICLE_MCP_URL` | Empty |
| `searchTool` | `CHRONICLE_SEARCH_TOOL` | Empty |
| `retrieveTool` | `CHRONICLE_RETRIEVE_TOOL` | Empty |
| `maxRequests` | `CHRONICLE_MAX_REQUESTS` | 40; range 4-100 |
| `concurrency` | `CHRONICLE_CONCURRENCY` | 2; range 1-4 |

`openai` selects a compatible chat-completions endpoint; include `/v1` when
required. MCP uses Streamable HTTP and explicitly selected tool names. These
values can configure autonomous research without any Settings interaction.
Assistant-led capture and save do not require these services.

## Speech

These fields are inside the JSON `speech` object.

| JSON field | Environment variable | Default |
| --- | --- | --- |
| `endpoint` | `CHRONICLE_SPEECH_URL` | `https://texttospeech.googleapis.com/v1` |
| `credentialsFile` | `GOOGLE_APPLICATION_CREDENTIALS` | Google default credentials discovery |
| `quotaProject` | `GOOGLE_CLOUD_QUOTA_PROJECT` | Credential quota project |

Credential paths must exist inside the container. The optional Compose speech
override mounts a host file; see [Docker](docker.md). Voice and language remain
available in the existing speech controls.

## Development and tool bridge

| Variable | Purpose | Default |
| --- | --- | --- |
| `CHRONICLE_DEV_HOST` | Vite host | `127.0.0.1` |
| `CHRONICLE_DEV_PORT` | Vite port | `5173` |
| `CHRONICLE_DEV_API_URL` | Vite API proxy | `http://127.0.0.1:8080` |
| `CHRONICLE_API` | CLI assistant-tool API | `http://127.0.0.1:8080` |

Vite uses a strict port. Add changed browser origins to the backend configuration.
Production Docker serves the built UI through Go and does not run Vite.
