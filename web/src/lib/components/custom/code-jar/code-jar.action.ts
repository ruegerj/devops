import type { CodeJar } from 'codejar';
import type { ActionReturn } from 'svelte/action';

type HighlightFunc = (text: string, syntax: string) => string;
type HighlightElementFunc = (element: HTMLElement) => void;

export interface CodeJarOptions {
	value: string;
	syntax: string;
	highlight: HighlightFunc;
}

export function codejar(node: HTMLElement, options: CodeJarOptions): ActionReturn<CodeJarOptions> {
	let jar: CodeJar;

	// import component on mount since it requires a window obj to be present upon initialization (would fail for ssr)
	import('codejar').then(({ CodeJar }) => {
		jar = CodeJar(node, wrapHighlight(options.highlight, options.syntax));
	});

	return {
		destroy() {
			jar?.destroy();
		}
	};
}

function wrapHighlight(highlight: HighlightFunc, syntax: string): HighlightElementFunc {
	return (element: HTMLElement) => {
		element.innerHTML = highlight(element.textContent ?? '', syntax);
	};
}
