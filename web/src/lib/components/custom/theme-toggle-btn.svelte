<script lang="ts">
	import IconMoon from 'virtual:icons/lucide/moon';
	import IconSun from 'virtual:icons/lucide/sun';

	import { toggleMode } from 'mode-watcher';
	import { Button } from '$lib/components/ui/button';
</script>

<Button onclick={toggleMode} variant="outline" size="icon" class="cursor-pointer">
	<IconSun
		class="h-[1.2rem] w-[1.2rem] scale-100 rotate-0 transition-all dark:scale-0 dark:-rotate-90"
	/>
	<IconMoon
		class="absolute h-[1.2rem] w-[1.2rem] scale-0 rotate-90 transition-all dark:scale-100 dark:rotate-0"
	/>
	<span class="sr-only">Toggle theme</span>
</Button>
