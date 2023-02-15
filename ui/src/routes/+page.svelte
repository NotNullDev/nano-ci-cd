<script>
	import ButtonBase from '../components/buttonBase.svelte';
	import App from './app.svelte';
</script>

<svelte:head>
	<title>Dashboard</title>
</svelte:head>

<div class="p-4 px-8 shadow-xl justify-between flex gap-2">
	<div>
		<input placeholder="search by name" class="bg-gray-800 border-gray-700 border rounded p-2" />
	</div>
	<div>
		<input placeholder="app name" class="bg-gray-800 border-gray-700 border rounded p-2" />
		<ButtonBase
			on:click={() => {
				console.log('add');
			}}>Add</ButtonBase
		>
	</div>
</div>

<div class="mt-5 flex flex-wrap gap-2 mb-4">
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
	<App />
</div>
