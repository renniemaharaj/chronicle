import { mkdir, readFile, writeFile } from 'node:fs/promises';
import { resolve } from 'node:path';

const [name = 'list', file] = process.argv.slice(2);
const endpoint = process.env.PLAYWRIGHT_ENDPOINT ?? 'http://127.0.0.1:8931/mcp';
const directory = resolve('.chronicle/social');
const sessionFile = resolve(directory, 'playwright-session.json');
await mkdir(directory, { recursive: true });

let session = '';
try {
  const saved = JSON.parse(await readFile(sessionFile, 'utf8'));
  if (saved.endpoint === endpoint) {
    session = saved.session;
  }
} catch {
  // The local MCP session is established on first use.
}

async function rpc(method, params) {
  const response = await fetch(endpoint, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Accept: 'application/json, text/event-stream',
      'MCP-Protocol-Version': '2025-11-25',
      ...(session ? { 'Mcp-Session-Id': session } : {}),
    },
    body: JSON.stringify({
      jsonrpc: '2.0',
      id: Date.now(),
      method,
      params,
    }),
    signal: AbortSignal.timeout(120000),
  });
  const text = await response.text();
  if (!response.ok) {
    throw new Error(`Playwright returned ${response.status}: ${text.slice(0, 400)}`);
  }

  if (response.headers.get('mcp-session-id')) {
    session = response.headers.get('mcp-session-id');
    await writeFile(sessionFile, JSON.stringify({ endpoint, session }), { mode: 0o600 });
  }

  const value = JSON.parse(
    text.startsWith('{')
      ? text
      : text
          .split('\n')
          .findLast((line) => line.startsWith('data:'))
          .slice(5),
  );
  if (value.error) {
    throw new Error(JSON.stringify(value.error));
  }

  return value.result;
}

if (!session) {
  await rpc('initialize', {
    protocolVersion: '2025-11-25',
    capabilities: {},
    clientInfo: {
      name: 'chronicle-series-driver',
      version: '1',
    },
  });
}

const args = file ? JSON.parse(await readFile(file, 'utf8')) : {};
const result = await rpc(
  name === 'list' ? 'tools/list' : 'tools/call',
  name === 'list'
    ? {}
    : {
        name,
        arguments: args,
      },
);
if (name === 'list') {
  await writeFile(resolve(directory, 'playwright-tools.json'), JSON.stringify(result, null, 2));
  process.stdout.write(
    JSON.stringify(
      result.tools.map((tool) => ({
        name: tool.name,
        arguments: tool.inputSchema.properties,
      })),
      null,
      2,
    ),
  );
} else {
  process.stdout.write(JSON.stringify(result, null, 2).replace(/token=[^\s"&\\]+/g, 'token=[redacted]'));
  if (result.isError) {
    process.exitCode = 1;
  }
}
