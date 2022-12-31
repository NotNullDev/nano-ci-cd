<script lang="ts">
	import { beforeNavigate } from '$app/navigation';
	import { assets, base } from '$app/paths';
	import '../app.css';

	beforeNavigate(() => {
		console.log('bye!');
	});
	console.log(base);
	console.log(assets);
</script>

<svelte:head>
	<title>Nano CI CD</title>
	<meta name="viewport" content="width=device-width, initial-scale=1" />
	<link rel="icon" type="image/png" href="/favicon.png" />
</svelte:head>

<div class="min-w-screen min-h-screen  text-gray-300 flex flex-col">
	<a class="p-6 shadow" href="/">
		<button class="btn btn-ghost ">Nano CI CD</button>
	</a>
	<slot />
</div>
