<script lang="ts">
	import clsx from 'clsx';

	export function onChange() {}

	export let val = '';
	export let placeholder = '';
	export let className = '';
</script>

<input
	class={clsx('bg-gray-800 border-gray-700 border rounded p-2 h-min', className)}
	{placeholder}
	on:keydown
	on:blur
	bind:value={val}
/>
