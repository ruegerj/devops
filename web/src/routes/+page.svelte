<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { toast } from 'svelte-sonner';
	import IconUnlockKeyhole from 'virtual:icons/lucide/unlock-keyhole';
	import IconLockKeyhole from 'virtual:icons/lucide/lock-keyhole';

	let secretContent = $state('');

	async function fetchSecret() {
		const response = await fetch('/api/secret');
		if (!response.ok) {
			console.error('Secret retrieval failed', response.status, response.statusText);
			toast.error('Vault is temporarily sealed - try again later...');
			return;
		}

		const content = await response.json();
		secretContent = JSON.stringify(content, null, 3);
	}

	function reloadPage() {
		window.location.reload();
	}
</script>

<div class="px-16 md:px-20 lg:px-24 xl:px-28">
	<Card.Root>
		<Card.Header class="flex justify-center">
			<Card.Title class=" text-4xl font-extrabold tracking-tight text-balance">
				Secret Vault
			</Card.Title>
		</Card.Header>
		<Card.Content class="flex justify-center">
			{#if !secretContent}
				<p class="text-gray-500">Open it up to see its contents...</p>
			{:else}
				<div>
					<pre>{secretContent}</pre>
				</div>
			{/if}
		</Card.Content>
		<Card.Footer class="flex justify-center">
			{#if !secretContent}
				<Button class="cursor-pointer" onclick={fetchSecret}>
					<IconUnlockKeyhole />
					Open
				</Button>
			{:else}
				<Button class="cursor-pointer" onclick={reloadPage}>
					<IconLockKeyhole />
					Lock
				</Button>
			{/if}
		</Card.Footer>
	</Card.Root>
</div>
