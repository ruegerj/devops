import type { CodeJar } from 'codejar';
import type { ActionReturn } from 'svelte/action';

type HighlightFunc = (text: string, syntax: string) => string;
type HighlightElementFunc = (element: HTMLElement) => void;

export interface CodeJarOptions {
	value: string;
	syntax: string;
	editorEnabled: boolean;
	highlight: HighlightFunc;
	onUpdate: (value: string) => void;
}

export function codejar(node: HTMLElement, options: CodeJarOptions): ActionReturn<CodeJarOptions> {
	let jar: CodeJar;

	// import component on mount since it requires a window obj to be present upon initialization (would fail for ssr)
	import('codejar').then(({ CodeJar }) => {
		jar = CodeJar(node, wrapHighlight(options.highlight, options.syntax));
		jar.onUpdate(options.onUpdate);
		updateVisibility(node, options.editorEnabled);
	});

	return {
		destroy() {
			jar?.destroy();
		},
		update(newOptions: CodeJarOptions) {
			if (jar && newOptions.value !== jar.toString()) {
				const pos = jar.save();
				jar.updateCode(newOptions.value);
				jar.restore(pos);
			}
			if (newOptions.editorEnabled !== options.editorEnabled) {
				updateVisibility(node, newOptions.editorEnabled);
			}
		}
	};
}

function updateVisibility(node: HTMLElement, enabled: boolean) {
	node.contentEditable = enabled ? 'plaintext-only' : 'false';
}

function wrapHighlight(highlight: HighlightFunc, syntax: string): HighlightElementFunc {
	return (element: HTMLElement) => {
		element.innerHTML = highlight(element.textContent ?? '', syntax);
	};
}
