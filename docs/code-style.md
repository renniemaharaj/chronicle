# Readable code

Separate functions and distinct logical steps with blank lines. Keep setup,
its immediate guard and cleanup together, then separate the next operation.
Keep related assignments together. Do not compress logic to fit the file limit.

Expand control-flow bodies and put statements on separate lines. Use clear names
and direct guards. Multi-field configuration literals use one field per line.
Short comments explain constraints or surprising behavior without banners.

Svelte control blocks and substantial child elements occupy separate lines.
Multi-attribute elements place attributes on separate lines. CSS declarations
and rules remain expanded. Preserve content, ordering, bindings and events.

```sh
npm run format
npm run check
```

The formatter combines AST-based control-flow and spacing edits with Prettier
and gofmt. `format:check` checks without writing. Authors still review logical
grouping and responsibility boundaries.

Every authored file contains fewer than 150 lines. Split by responsibility
before reaching the limit. Generated dependency lockfiles are exempt. Build
output, dependencies, secrets and runtime storage stay out of version control.
