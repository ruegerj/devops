<script lang="ts">
	/* eslint-disable svelte/no-at-html-tags */
	import Prism from 'prismjs';
	import 'prismjs/components/prism-json';
	import { codejar } from './code-jar.action';

	interface Props {
		class: string;
		code: string;
	}

	let { class: clazz, code = $bindable() }: Props = $props();

	let container = $state<HTMLPreElement>();

	function highlight(text: string, syntax: string): string {
		return Prism.highlight(text, Prism.languages[syntax], syntax);
	}
</script>

<pre
	class={clazz ?? ''}
	bind:this={container}
	use:codejar={{
		value: code,
		syntax: 'json',
		highlight: highlight
	}}>{@html highlight(code, 'json')}</pre>
