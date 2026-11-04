<script lang="ts">
  let { size = 16 }: { size?: number } = $props();
</script>

<span
  class="facebook-mark"
  style={`--mark-size: ${size}px`}
  aria-hidden="true"
>
  f
</span>

<style>
  .facebook-mark {
    display: inline-flex;
    align-items: end;
    justify-content: center;
    width: var(--mark-size);
    height: var(--mark-size);
    flex-shrink: 0;
    border-radius: 4px;
    background: #1877f2;
    color: white;
    font:
      700 var(--mark-size) Arial,
      sans-serif;
    line-height: 0.95;
  }
</style>
