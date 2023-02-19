<script lang="ts">
	import dayjs from 'dayjs';
	import ButtonBase from '../../components/buttonBase.svelte';
	import InputBase from '../../components/inputBase.svelte';

	import { detailsPageStore } from '../../logic/details/store';

	export let closeSidebar: () => void | null;
</script>

<div class="p-4 shadow text-2xl">Build history</div>
<InputBase placeholder="Search" />
<div class="flex flex-col gap-2 mt-4">
	{#each $detailsPageStore.availableBuilds ?? [] as build}
		<div class="flex gap-2 p-2 justify-between w-full items-center">
			<div class="flex justify-between">
				<div class="mr-2 border-r-gray-900 pr-2 border-r-2">{build.id}</div>
				<div>{dayjs(build.date).format('DD-MM-YYYY mm:HH')}</div>
			</div>

			<ButtonBase
				on:click={() => {
					detailsPageStore.setKey('selectedBuildId', build.id);
					closeSidebar();
				}}>Select a</ButtonBase
			>
		</div>
	{/each}
</div>
