<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import { ModeWatcher } from 'mode-watcher';
	import { ThemeToggleBtn } from '$lib/components/custom';
	import { Toaster } from '$lib/components/ui/sonner';

	let { children } = $props();
</script>

<svelte:head>
	<title>Secret Vault</title>
	<link rel="icon" href={favicon} />
</svelte:head>

<div class="flex h-full flex-col overflow-hidden">
	<nav class="flex h-16 items-center px-4">
		<div class="ml-auto flex items-center space-x-4">
			<ThemeToggleBtn />
		</div>
	</nav>

	<main class="flex-auto overflow-hidden p-4">
		{@render children?.()}
	</main>
</div>

<ModeWatcher />
<Toaster />
