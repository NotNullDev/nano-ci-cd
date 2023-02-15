<script>
	import { onMount } from 'svelte';
	import ButtonBase from '../../components/buttonBase.svelte';

	let appid = '';
	let logs = 'not logs available';

	let history = 'not history available';

	onMount(() => {
		appid = location?.href?.split('appId=')[1] ?? 'ERROR';
	});
</script>

<div class="p-4 px-8 shadow-xl justify-between flex gap-2">
	<div>App 1 [idle]</div>
	<div>
		<ButtonBase accent="danger">Delete app</ButtonBase>
	</div>
</div>

<div class="flex gap-2 pt-5 pr-3 ">
	<!-- logs -->
	<div class="w-full flex flex-col gap-2 flex-1">
		<!-- logs header -->
		<div class="flex items-end gap-2 h-[75px]">
			<input
				placeholder="search"
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

		<!-- logs content -->
		<div class="flex p-4 shadow-gray-900 shadow-xl h-[50vh] w-full">
			{logs}
		</div>

		<!-- logs footer -->
		<div class="mt-3">
			<ButtonBase>Download</ButtonBase>
			<ButtonBase>Auto refresh</ButtonBase>
		</div>
	</div>

	<!-- build history -->
	<div class="w-1/2 flex flex-col gap-2">
		<!-- history header -->
		<div class="h-[75px]">
			<h2>Build history</h2>
			<input
				placeholder="search"
				type="search"
				class="bg-gray-800 border-gray-700 border rounded p-2 h-min"
			/>
		</div>
		<!-- history list -->
		<div class="flex p-4 shadow-gray-900 shadow-xl w-full h-[50vh]">
			{history}
		</div>
	</div>
</div>
