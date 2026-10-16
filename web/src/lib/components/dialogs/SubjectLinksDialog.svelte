<script lang="ts">
  import { ExternalLink, Globe2, ShieldCheck } from '@lucide/svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { safeSourceURL } from '../../api/links';
  const entity = $derived(ws.entity ?? ws.subject);
  const links = $derived((ws.project.subjectLinks ?? []).filter((item) => item.entityId === entity.id));
  const verified = $derived(links.filter((item) => item.verified && safeSourceURL(item.url)));
</script>

<p class="dialog-intro">
  Official websites and social accounts for {entity.name}, with the source used to confirm each link.
</p>
<div class="official-links">
  {#each verified as link}
    {@const source = ws.project.sources.find((item) => item.id === link.sourceId)}
    <article class="official-link">
      <a class="official-link-title" href={safeSourceURL(link.url)} target="_blank" rel="noopener noreferrer">
        <Globe2 size={22} /><span><strong>{link.label}</strong><small>{link.url}</small></span><ExternalLink
          size={16}
        />
      </a>
      <p>
        <ShieldCheck size={14} />{link.verificationMethod === 'official_site_link'
          ? 'Linked by the official website'
          : 'Official website reference'}
      </p>
      <small>Checked {link.verifiedAt?.slice(0, 10) ?? 'date not recorded'}</small>
      {#if safeSourceURL(source?.url)}<a href={safeSourceURL(source?.url)} target="_blank" rel="noopener noreferrer"
          >Verification source: {source?.title}</a
        >{/if}
    </article>
  {/each}
</div>
{#if !verified.length}<p class="empty-state">No official links have been verified for this subject yet.</p>{/if}
{#if links.length > verified.length}<p class="subject-search-note">
    {links.length - verified.length} candidates could not be verified. They are preserved in the research record.
  </p>{/if}
<p class="subject-search-note">
  Verification here means an evidenced official reference, rather than a platform badge. Public accounts may be
  maintained by a subject’s representatives.
</p>
