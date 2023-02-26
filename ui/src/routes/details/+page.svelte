<script lang="ts">
	import { onMount } from 'svelte';
	import { get } from 'svelte/store';
	import ButtonBase from '../../components/buttonBase.svelte';
	import Sidebar from '../../components/sidebar.svelte';
	import { NanoUtils } from '../../logic/common/pure';
	import { runBuild, useFetchBuilds } from '../../logic/details/api';
	import { detailsPageStore } from '../../logic/details/store';
	import { indexPageStore } from '../../logic/index/store';
	import type { App } from '../../types/types';
	import SidebarContent from './sidebarContent.svelte';

	let appid = '';
	let logs = 'no logs available';

	let history = 'not history available';

	let app: App | null = null;

	onMount(() => {
		appid = location?.href?.split('appId=')[1] ?? 'ERROR';

		app = get(indexPageStore).apps.find((app) => app.ID === Number(appid)) ?? null;
		useFetchBuilds();
	});

	let openSidebar: () => void | null;
	let closeSidebar: () => void | null;
</script>

<div class="p-4 px-8 shadow-xl justify-between flex gap-2">
	<div class="flex gap-4 items-center">
		<div>App {app?.appName} [idle]</div>
		<ButtonBase
			on:click={() => {
				if (app) {
					runBuild(app.appName);
				}
			}}>Build</ButtonBase
		>
		<ButtonBase
			on:click={() => {
				openSidebar();
			}}>History</ButtonBase
		>
	</div>
	<div>
		<ButtonBase accent="danger">Delete app</ButtonBase>
	</div>
</div>
<div class="flex flex-1">
	<div class="w-[150px] h-full flex flex-col gap-2 mt-10 mr-3">
		<ButtonBase>build</ButtonBase>
		<ButtonBase>force stop</ButtonBase>
		<ButtonBase>build history</ButtonBase>
		<ButtonBase>logs</ButtonBase>
		<ButtonBase>details</ButtonBase>
	</div>
	<div class="flex-1">
		<div class="flex gap-2 pt-5 pr-3 ">
			<!-- logs -->
			<div class="w-full flex flex-col gap-2 flex-1">
				<!-- logs header -->
				<div class="flex items-end gap-2 h-[75px]">
					<input
						placeholder="filter"
						type="search"
						class="bg-gray-800 border-gray-700 border rounded p-2 h-min"
					/>
					<label class="ml-4">
						<div>from</div>
						<input type="datetime-local" class="bg-gray-800 border-gray-700 border rounded p-2" />
					</label>
					<label>
						<div>to</div>
						<input type="datetime-local" class="bg-gray-800 border-gray-700 border rounded p-2" />
					</label>
				</div>

				<div class="flex gap-4">
					<div>Build ID: <span>{$detailsPageStore.currentBuild?.ID ?? '-'}</span></div>
					<div>
						Started at: <span
							>{NanoUtils.formatDate($detailsPageStore.currentBuild?.startedAt ?? '') ?? '-'}</span
						>
					</div>
					<div>
						Finished at: <span
							>{NanoUtils.formatDate($detailsPageStore.currentBuild?.finishedAt ?? '') ?? '-'}</span
						>
					</div>
					<div>Status: <span>{$detailsPageStore.currentBuild?.buildStatus ?? '-'}</span></div>
				</div>

				<!-- logs content -->
				<div class="flex p-4 shadow-gray-900 shadow-xl h-[50vh] w-full">
					{$detailsPageStore.currentBuild?.logs ?? 'no logs available'}
				</div>

				<!-- logs footer -->
				<div class="mt-3">
					<ButtonBase>Download</ButtonBase>
					<ButtonBase>Auto refresh</ButtonBase>
				</div>
			</div>
		</div>

		{#key $detailsPageStore.availableBuilds?.length}
			<Sidebar
				bind:openSidebar
				bind:closeSidebar
				className="text-slate-300 flex items-center flex-col"
			>
				<SidebarContent {closeSidebar} />
			</Sidebar>
		{/key}
	</div>
</div>
<!-- <Portal class="absolute inset-0 pointer-events-none">
	<div>wtf!!</div>
</Portal> -->
