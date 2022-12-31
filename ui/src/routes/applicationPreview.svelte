<script>
	import { goto } from '$app/navigation';
</script>

<div class="flex items-center justify-center">
	<div class="w-[300px] h-[250px] rounded-xl p-4 shadow-xl    cursor-pointer flex flex-col">
		<h2 class="text-center text-xl font-bold">App 1</h2>
		<div class="flex flex-1 flex-col justify-end">
			<button
				class="btn btn-primary btn-sm w-1/3"
				on:click={() => {
					goto('/app');
				}}>Details</button
			>
		</div>
	</div>
</div>
