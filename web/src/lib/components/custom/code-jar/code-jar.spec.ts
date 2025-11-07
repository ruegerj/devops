// @vitest-environment jsdom

import { mount, unmount } from 'svelte';
import { describe, expect, it } from 'vitest';
import CodeJar from '$lib/components/custom/code-jar/code-jar.svelte';

describe('CodeJar Component', () => {
	it('should mount and render the given code', () => {
		// GIVEN
		const exptedCode = `{ "message": "Hello world" }`;

		// WHEN
		const codeJar = mount(CodeJar, {
			target: document.body,
			props: {
				code: exptedCode,
				class: ''
			}
		});
		const node = document.querySelector('pre');

		// THEN
		expect(node).toBeDefined();
		expect(node!.textContent).toBe(exptedCode);

		unmount(codeJar);
	});

	it('should apply the given classes to the dom node', () => {
		// GIVEN
		const expectedClasses = 'py-3 font-bold';

		// WHEN
		const codeJar = mount(CodeJar, {
			target: document.body,
			props: {
				code: `"foobar"`,
				class: expectedClasses
			}
		});
		const node = document.querySelector('pre');

		// THEN
		expect(node).toBeDefined();
		expect(node!.className).toBe(expectedClasses);

		unmount(codeJar);
	});
});
