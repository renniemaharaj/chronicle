# Docker

The image contains the Go API and built Svelte UI. Configuration is applied at
runtime. Docker Compose 2.24 or newer is required.

```sh
cp .env.example .env
mkdir -p -m 700 .chronicle
docker compose up --build -d
```

Open http://localhost:8080. The default host binding is local to this machine.
Compose mounts the existing host `.chronicle/` folder read/write at `/data`.
Captures, projects, and saved settings are shared with the native app and survive
container recreation. Back up this host folder to preserve your research.

Compose runs as UID/GID 1000 by default. Set `CHRONICLE_UID` and `CHRONICLE_GID`
in `.env` to the output of `id -u` and `id -g` if your host IDs differ. The folder
keeps its existing ownership and permissions. Missing folders cause startup to
fail instead of silently creating an empty workspace owned by root.

Set `CHRONICLE_WORKSPACE_HOST` to mount another existing host workspace at `/data`.
Keep the container's `CHRONICLE_WORKSPACE` set to `/data`. Previous deployments
used the `chronicle_workspace` named volume; that volume is retained when switching
to this bind mount. Inspect it before removing it if it contains earlier research.

## Configuration

Edit `.env`, or copy `config/docker.json` to `config/local.json` and set
`CHRONICLE_CONFIG_FILE=./config/local.json` in `.env`. Apply edits with:

```sh
docker compose up -d --force-recreate chronicle
```

To use another environment file, select it for interpolation and the container:

```sh
CHRONICLE_ENV_FILE=deployment.env docker compose --env-file deployment.env up -d
```

`CHRONICLE_PORT` changes the host port. To change the container listener, set
`CHRONICLE_ADDR=0.0.0.0:9000` and `CHRONICLE_CONTAINER_PORT=9000` together.
The health check reads that listener. Keep container listeners reachable on all
container interfaces. See [Configuration](configuration.md) for every service.

## PostgreSQL

Add these values to `.env`, choosing a password for your installation:

```dotenv
COMPOSE_PROFILES=postgres
POSTGRES_DB=chronicle
POSTGRES_USER=chronicle
POSTGRES_PASSWORD=replace-me
DATABASE_URL=postgres://chronicle:replace-me@postgres:5432/chronicle?sslmode=disable
```

Run `docker compose up --build -d`. The app waits for database health and runs
migrations at startup. PostgreSQL has a separate persistent volume and no
published port. Percent-encode special characters in URL credentials. External
PostgreSQL only needs `DATABASE_URL`; switching storage does not migrate data.

## Models and research tools

Enable `COMPOSE_PROFILES=ollama` and use
`CHRONICLE_LLM_URL=http://ollama:11434`. Combine profiles with
`COMPOSE_PROFILES=postgres,ollama`. Pull the model you want:

```sh
docker compose up -d ollama
docker compose exec ollama ollama pull YOUR_MODEL
```

Set `CHRONICLE_LLM_MODEL` accordingly. Compatible providers use
`CHRONICLE_LLM_PROVIDER=openai`, their URL, model and optional key. Set
`CHRONICLE_MCP_URL`, `CHRONICLE_SEARCH_TOOL`, and `CHRONICLE_RETRIEVE_TOOL`,
or configure them through Settings. No model weights or search service are
downloaded automatically with the app. Assistant research at `/api/mcp` needs
no separately configured model.

Host services are reachable as `host.docker.internal` when they listen on a
Docker-reachable interface. A host listener restricted to `127.0.0.1` is not
reachable that way. Container `localhost` refers to the container itself.

## Google speech

Set `GOOGLE_ADC_FILE` to a host credentials JSON file and optionally set
`GOOGLE_CLOUD_QUOTA_PROJECT`, then mount it using:

```sh
docker compose -f compose.yaml -f compose.speech.yaml up --build -d
```

The read-only credential is mounted at `/run/secrets/google_adc` and must be
readable by the configured `CHRONICLE_UID`. Browser speech remains available without Google setup.
See [Speech](speech.md) for credentials and voices.

## Image and service options

Compose accepts `CHRONICLE_IMAGE`, `NODE_VERSION`, `GO_VERSION`,
`ALPINE_VERSION`, `POSTGRES_IMAGE`, and `OLLAMA_IMAGE`. Choose tags or digests
according to your update policy. Ollama also accepts `OLLAMA_HOST`,
`OLLAMA_KEEP_ALIVE` and `OLLAMA_NUM_PARALLEL`. `CHRONICLE_BIND` controls the host binding.

A standalone image accepts `docker run --env-file` and JSON mounted at the path
selected by `CHRONICLE_CONFIG`. Image defaults use `/etc/chronicle/config.json`,
`/data`, port 8080, and UID 10001. To share the host workspace directly:

```sh
docker run --rm --user "$(id -u):$(id -g)" -p 127.0.0.1:8080:8080 \
  --mount "type=bind,source=$PWD/.chronicle,target=/data" chronicle:local
```

The host folder must exist. See [Docker's environment documentation](https://docs.docker.com/compose/how-tos/environment-variables/set-environment-variables/).
