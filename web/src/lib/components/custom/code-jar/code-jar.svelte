<script lang="ts">
	/* eslint-disable svelte/no-at-html-tags */
	import Prism from 'prismjs';
	import 'prismjs/components/prism-json';
	import { codejar } from './code-jar.action';

	interface Props {
		class: string;
		code: string;
		enabled?: boolean;
	}

	let { class: clazz, code = $bindable(), enabled }: Props = $props();

	let container = $state<HTMLPreElement>();

	function highlight(text: string, syntax: string): string {
		return Prism.highlight(text, Prism.languages[syntax], syntax);
	}

	function onCodeUpdate(newValue: string) {
		code = newValue;
	}
</script>

<pre
	class={clazz ?? ''}
	bind:this={container}
	use:codejar={{
		value: code,
		syntax: 'json',
		editorEnabled: enabled ?? true,
		highlight: highlight,
		onUpdate: onCodeUpdate
	}}>{@html highlight(code, 'json')}</pre>
