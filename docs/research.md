# Research behavior

## Stages

1. Provisionally identify the subject and adapt four chronological questions to its type.
2. Discover configured MCP capabilities and inspect their schemas.
3. Search independently for each research question, within bounded concurrency.
4. Choose an accessible URL present in the actual search output.
5. Retrieve it through the explicitly configured tool.
6. Extract typed entities, events, relationships, claims, and short quotations.
7. Reject quotations absent from the retrieved content. Unverified claims remain explicit.
8. Merge stable identities and preserve alternative claims.
9. Seek an independent source or missing period when the remaining budget permits.
10. Assess evidence, check temporal consistency, assemble chapters, and save the complete project.

## Limits and failures

The request budget includes model calls and research tool calls. Tool discovery has its own 150-tool ceiling.
Large documents are bounded before model processing and recorded as only partially reviewed.
Model requests time out; MCP calls have bounded timeouts and transport reconnection.
A failed question leaves the session incomplete and resumable. Completed questions are not repeated when resuming.
Pause and cancel stop in-flight calls through context cancellation. Restarted sessions are marked paused.

Structured model output is validated against inferred JSON Schema and decoded with unknown fields rejected.
Deterministic code validates the ontology and references. Invalid extraction remains visible in the activity log;
it cannot enter canonical knowledge silently.

The model never receives authority to choose arbitrary tools. Search and retrieval are the only configured capabilities.
A document containing instructions cannot authorize additional actions. Local file and write capabilities are not enabled.
Source access restrictions and retrieval limitations remain the responsibility of the configured research tools;
Chronicle does not bypass them.

## Assessment

A verified quotation proves passage matching, not semantic entailment or historical truth.
One verified supporting passage produces `supported_limited`.
Two explicitly reviewed independent groups with supporting passages can produce `corroborated`.
Supporting and contradicting passages produce `disputed`. Rejected or superseded assertions remain stored.
Host-based source grouping is a conservative initial hint and requires explicit independence review.

The default pipeline preserves incomplete coverage. It cannot establish historical completeness from a successful run.
Narrative chapters currently assemble supported assertions deterministically. They do not introduce unsupported causal stories.
