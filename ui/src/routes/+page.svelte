<script lang="ts">
	import { goto } from '$app/navigation';
	import ButtonBase from '../components/buttonBase.svelte';
	import { authStore } from '../logic/common/store';
	import { fetchNanoContext } from '../logic/index/api';
	import type { NanoContext } from '../types/types';
	import App from './app.svelte';

	let nanoContext: NanoContext | null = null;

	let isLoggedIn = false;

	const auth = authStore().subscribe((value) => {
		isLoggedIn = value.isLoggedIn;
	});

	let isBrowser = typeof window !== 'undefined';

	$: {
		if (isBrowser && !isLoggedIn) {
			goto('/login');
		}

		if (isBrowser && isLoggedIn) {
			(async () => {
				nanoContext = await fetchNanoContext();
			})();
		}
	}
</script>

<svelte:head>
	<title>Dashboard</title>
</svelte:head>

{#if isLoggedIn}
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

	{JSON.stringify(nanoContext, null, 2)}

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
{:else}
	<div>You must be logged in to proceed</div>
{/if}
