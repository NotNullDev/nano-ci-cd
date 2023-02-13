<script lang="ts">
	export let href = '/';
	const classes =
		'hover:bg-gray-700 p-1 cursor-pointer rounded px-2 active:scale-95 ' + $$restProps.class;
</script>

<a {href} {...$$restProps} class={classes} on:click>
	<slot />
</a>
