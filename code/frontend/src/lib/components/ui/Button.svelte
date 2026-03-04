<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		type?: 'button' | 'submit' | 'reset';
		variant?: 'primary' | 'secondary' | 'ghost' | 'outline';
		size?: 'sm' | 'md' | 'lg';
		disabled?: boolean;
		onclick?: () => void;
		children: Snippet;
		class?: string;
	}

	let {
		type = 'button',
		variant = 'primary',
		size = 'md',
		disabled = false,
		onclick,
		children,
		class: className = ''
	}: Props = $props();

	const variants = {
		primary: 'bg-brand-primary hover:bg-brand-primary/90 text-white shadow-lg glow-primary',
		secondary: 'bg-brand-secondary hover:bg-brand-secondary/90 text-white shadow-md',
		outline: 'border border-white/20 hover:bg-white/5 text-white',
		ghost: 'hover:bg-white/5 text-white/80 hover:text-white'
	};

	const sizes = {
		sm: 'px-3 py-1.5 text-sm',
		md: 'px-5 py-2.5 text-base',
		lg: 'px-8 py-3.5 text-lg font-semibold'
	};
</script>

<button
	{type}
	{disabled}
	{onclick}
	class="inline-flex items-center justify-center rounded-xl transition-all duration-200 active:scale-95 disabled:opacity-50 disabled:pointer-events-none {variants[
		variant
	]} {sizes[size]} {className}"
>
	{@render children()}
</button>
