<script lang="ts">
  let { size = 32 }: { size?: number } = $props();
</script>

<span
  class="chronicle-mark"
  style:width={`${size}px`}
  style:height={`${size}px`}
  aria-hidden="true"
></span>

<style>
  .chronicle-mark {
    display: block;
    flex-shrink: 0;
    color: var(--blue);
    background-color: currentColor;
    mask: url('/favicon.svg') center / contain no-repeat;
  }
</style>
