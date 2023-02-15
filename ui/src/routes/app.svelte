<script lang="ts">
	import dayjs from 'dayjs';
	import LinkBase from '../components/linkBase.svelte';

	// input

	export let appName = 'App 1';
	export let appStatus = 'idle';
	export let lastUpdated = dayjs().format('DD.MM.YYYY HH:mm').toString();
	export let appId = 'fsdafadsf';
</script>

<div
	class="flex flex-col p-4 w-48 h-52 border-gray-700 border rounded hover:shadow-indigo-900 hover:shadow hover:scale-105
transition-all duration-300"
>
	<div class="mb-2">{appName}</div>
	<div class="flex-1 flex flex-col gap-2">
		<div>Status: {appStatus}</div>
		<div>Last updated: {lastUpdated}</div>
	</div>
	<div class="flex justify-end mt-2">
		<LinkBase href={'/details?appId=' + appId} class="text-sm">Details</LinkBase>
	</div>
</div>
