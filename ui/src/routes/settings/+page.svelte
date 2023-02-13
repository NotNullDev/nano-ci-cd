<script lang="ts">
	import { onDestroy, onMount } from 'svelte';

	let changingSize = false;

	const handleMouseDown = (e: MouseEvent) => {
		if (changingSize) {
			changingSize = false;
		}
	};

	const handleMouseMove = (e: MouseEvent) => {
		if (changingSize && panelRef) {
			console.log(e);
			panelRef.style.width = `${panelRef.offsetWidth + e.movementX}px`;
		}
	};

	onMount(() => {
		document.addEventListener('mouseup', handleMouseDown);
		document.addEventListener('mousemove', handleMouseMove);
	});

	onDestroy(() => {
		document.removeEventListener('mouseup', handleMouseDown);
		document.removeEventListener('mousemove', handleMouseMove);
	});

	let panelRef: HTMLDivElement | null = null;
</script>

<svelte:head>
	<title>Settings</title>
</svelte:head>

<div class="flex flex-1 m-4 border border-b">
	<div class="flex">
		<div class="w-[200px] p-8 border m-2 mr-0 border-r-0 resize-x" bind:this={panelRef}>
			panel 1
		</div>
		<div class="bg-slate-300 w-[1px] my-2  relative">
			<div
				class="absolute w-[20px] right-0 top-0 h-full translate-x-[21px] cursor-e-resize"
				on:mousedown={() => (changingSize = true)}
			/>
		</div>
	</div>
	<div class="flex-1 p-8 border m-2">panel 2</div>
</div>
