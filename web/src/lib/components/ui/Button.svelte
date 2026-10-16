<script lang="ts">
  import { Button as FlowbiteButton } from 'flowbite-svelte';
  import type { HTMLButtonAttributes } from 'svelte/elements';
  import type { Snippet } from 'svelte';

  type Variant = 'primary' | 'secondary' | 'ghost' | 'danger';
  type Props = Omit<HTMLButtonAttributes, 'class' | 'color' | 'disabled'> & {
    class?: string;
    children?: Snippet;
    size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl';
    tag?: 'button' | 'label';
    disabled?: boolean;
    variant?: Variant;
    active?: boolean;
    current?: boolean;
    pressed?: boolean;
    icon?: boolean;
  };

  let {
    children,
    class: className = '',
    variant = 'ghost',
    active,
    current,
    pressed,
    icon = false,
    size = 'sm',
    ...props
  }: Props = $props();

  const color = $derived(variant === 'danger' ? 'red' : variant === 'primary' ? 'primary' : 'alternative');
  const classes = $derived([
    className,
    active ? 'active' : '',
    current ? 'current' : '',
    pressed ? 'pressed' : '',
    variant === 'ghost' ? 'border-transparent bg-transparent dark:bg-transparent dark:border-transparent' : '',
    icon ? 'p-1.5 shrink-0 [&>svg]:shrink-0' : '',
  ].join(' '));
</script>

<FlowbiteButton {color} {size} class={classes} data-variant={variant} aria-pressed={active ?? pressed} {...props}>
  {@render children?.()}
</FlowbiteButton>
