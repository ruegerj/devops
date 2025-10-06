<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import { ModeWatcher } from 'mode-watcher';
	import { ThemeToggleBtn } from '$lib/components/custom';

	let { children } = $props();
</script>

<svelte:head>
	<title>Universe Vault</title>
	<link rel="icon" href={favicon} />
</svelte:head>

<main>
	<ThemeToggleBtn />
	{@render children?.()}
</main>

<ModeWatcher />
