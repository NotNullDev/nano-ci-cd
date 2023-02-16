<script lang="ts">
	import { goto } from '$app/navigation';
	import ButtonBase from '../components/buttonBase.svelte';
	import IconPlus from '../components/icons/IconPlus.svelte';
	import { authStore } from '../logic/common/store';
	import { createApp } from '../logic/index/api';
	import { refetchNanoContext } from '../logic/index/functions';
	import { indexPageStore } from '../logic/index/store';
	import App from './app.svelte';

	let isLoggedIn = $authStore.isLoggedIn;
	let isBrowser = typeof window !== 'undefined';

	authStore.subscribe((val) => {
		isLoggedIn = val.isLoggedIn;
	});

	$: {
		if (isBrowser && !isLoggedIn) {
			console.log('redirect to login');
			goto('/login');
		}
	}

	$: {
		if (isBrowser && isLoggedIn) {
			console.log('fetching data');
			(async () => {
				await refetchNanoContext();
			})();
		}
	}

	let newAppName = '';
	let appsFilter = '';
	let filteredApps = $indexPageStore.apps;

	$: {
		filteredApps = $indexPageStore.apps.filter((app) => {
			return app.appName.includes(appsFilter);
		});
	}
</script>

<svelte:head>
	<title>Dashboard</title>
</svelte:head>

{#if isLoggedIn}
	<div class="p-4 px-8 shadow-xl justify-between flex gap-2">
		<div>
			<input
				placeholder="find by name"
				class="bg-gray-800 border-gray-700 border rounded p-2"
				bind:value={appsFilter}
			/>
		</div>
		<div class="gap-2 flex">
			<input
				placeholder="app name"
				class="bg-gray-800 border-gray-700 border rounded p-2"
				bind:value={newAppName}
			/>
			<ButtonBase
				class="w-16 h-10 flex items-center"
				on:click={async () => {
					const newApp = await createApp(newAppName);
					await refetchNanoContext();
				}}
			>
				<IconPlus className="" />
				<div>Add</div>
			</ButtonBase>
		</div>
	</div>

	<div class="mt-5 flex flex-wrap gap-2 mb-4">
		{#each filteredApps as app}
			{#key app.ID}
				<App {app} />
			{/key}
		{/each}
	</div>
{/if}
