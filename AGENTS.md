# Chronicle conventions

- Every authored source, style, configuration, migration, and documentation file must contain fewer than 150 lines.
- Split by responsibility before a file reaches that limit.
- Keep dependency and build output out of version control, except the generated dependency lockfile.
- Use one universal ontology. Canvas nodes are projections of canonical records.
- Never invent citations or present synthetic fixture evidence as historical evidence.
- Preserve historical calendars, uncertainty, and competing claims.

## Research requests

When the user asks to research a subject, use Chronicle's local backend tools to save the result.
Read `docs/assistant-tools.md` and `docs/research-standards.md`. Prefer the registered `chronicle` MCP tools; if not available in the
current session, invoke the same tools with `node scripts/chronicle-tool.mjs` and JSON argument files.
Use available web search to investigate independently, then capture source passages through the backend.
Never populate research by editing storage files. Validate before saving, preserve competing claims,
and disclose coverage gaps. Assistant-led research does not need a separately configured LLM.
