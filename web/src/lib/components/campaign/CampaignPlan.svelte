<script lang="ts">
  import { Clapperboard, Clock3, X } from '@lucide/svelte';
  import { Input, Label } from 'flowbite-svelte';
  import Button from '../ui/Button.svelte';
  import type { CampaignWorkspace } from '../../campaign/workspace.svelte';

  let { model }: { model: CampaignWorkspace } = $props();
  let seconds = $state(60);
</script>

<section
  class="reel-create-card"
  aria-label="Create a reel series"
>
  <div class="reel-section-heading">
    <h2><Clapperboard size={18} />A new series, from your research</h2>
    <Button
      icon
      aria-label="Close series setup"
      disabled={model.locked}
      onclick={() => (model.createOpen = false)}
    >
      <X size={16} />
    </Button>
  </div>
  <p>Keep the full account in sequence, with manageable episodes and a saved place to resume.</p>
  <form
    class="reel-create-form"
    onsubmit={(event) => {
      event.preventDefault();
      void model.create(seconds);
    }}
  >
    <div>
      <Label for="reel-target">Target seconds per reel</Label>
      <Input
        id="reel-target"
        type="number"
        min="15"
        max="180"
        required
        bind:value={seconds}
      />
    </div>
    <div class="reel-create-note">
      <Clock3 size={18} />
      <span>
        Sentence boundaries come first.
        <br />
        Final duration depends on narration.
      </span>
    </div>
    <Button
      variant="primary"
      type="submit"
      disabled={model.locked}
    >
      Build series
    </Button>
  </form>
</section>
