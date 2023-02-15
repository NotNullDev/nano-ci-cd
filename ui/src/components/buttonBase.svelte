<svelte:options accessors={true} />

<script lang="ts">
	import clsx from 'clsx';
	import './buttonBase.d.ts';
	export let accent: 'primary' | 'danger' | 'ghost' = 'primary';

	const classesCombined = 'p-1 cursor-pointer rounded px-2  ' + $$restProps.class;
</script>

<button
	{...$$restProps}
	class={clsx(
		' active:scale-95  font-semibold ',
		{
			'bg-accent hover:bg-accent/80': accent === 'primary',
			'bg-danger hover:bg-danger/80': accent === 'danger',
			'bg-transparent hover:bg-gray-700 font-normal': accent === 'ghost'
		},
		classesCombined
	)}
	on:mouseenter
	on:click
>
	<slot />
</button>
