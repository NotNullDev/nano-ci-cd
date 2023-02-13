<svelte:options accessors={true} />

<script lang="ts">
	const classes = 'hover:bg-gray-700 p-1 cursor-pointer rounded px-2  ' + $$restProps.class;
</script>

<button {...$$restProps} class={classes + ' active:scale-95'} on:click>
	<slot />
</button>
