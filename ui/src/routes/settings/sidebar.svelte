<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import ButtonBase from '../../components/buttonBase.svelte';
	let open = false;

	let sidebar: HTMLElement | null = null;
	let sidebarOverlay: HTMLElement | null = null;
	let sidebarContainer: HTMLElement | null = null;
	let portal: HTMLElement | null = null;

	const closeSidebar = () => {
		open = false;
	};

	$: {
		if (sidebar && open) {
			console.log('open!');
			sidebar.classList.remove('sidebar-closed');
			sidebar.classList.add('sidebar-open');

			sidebarOverlay?.classList.add('bg-gray-900/50');
			sidebarOverlay?.classList.add('pointer-events-auto');
			sidebarOverlay?.classList.remove('pointer-events-none');
			sidebarOverlay?.addEventListener('click', closeSidebar);
		}

		if (sidebar && !open) {
			console.log('closed!');
			sidebar.classList.remove('sidebar-open');
			sidebar.classList.add('sidebar-closed');

			sidebarOverlay?.classList.remove('bg-gray-900/50');
			sidebarOverlay?.classList.remove('pointer-events-auto');
			sidebarOverlay?.classList.add('pointer-events-none');
			sidebarOverlay?.removeEventListener('click', closeSidebar);
		}
	}
	onMount(() => {
		if (sidebarContainer) {
			portal = sidebarContainer.cloneNode(true) as HTMLElement;
			document.body.append(portal);
			sidebar?.remove();
			sidebar = portal.querySelector('#sidebar') as HTMLElement;
			sidebarOverlay = portal.querySelector('#sidebar-overlay') as HTMLElement;
		}
	});
	onDestroy(() => {
		if (portal) {
			portal.remove();
		}
	});
</script>

<ButtonBase
	on:click={() => {
		open = !open;
	}}>Toggle me</ButtonBase
>

<button class="p-1 px-2 rounded btn relative w-min whitespace-nowrap btn group">
	hover me!
	<div
		class="absolutea left-0 bottom-0 h-[2px] shadow-violet-800 shadow rounded-xl w-0  bg-violet-900 btn-helper group-hover:w-3/4 transition-all"
	/>
</button>

<div
	class=" overflow-hidden h-screen w-screens absolute inset-0 pointer-events-none"
	bind:this={sidebarContainer}
>
	<button class="relative z-10 h-full w-full pointer-events-none" id="sidebar-overlay">
		<div
			id="sidebar"
			class="z-20 bg-gray-800 border-l border-l-gray-900 rounded-tl-xl shadow-violet-900 shadow  absolute top-0 right-0  h-screen w-[300px]"
			bind:this={sidebar}
		/>
	</button>
</div>

<style>
	/* #sidebar {
		position: absolute;
		right: 0;
		top: 0;
		width: 400px;
		height: 100vh;
		border-top-left-radius: 10px;
		transform: translateX(100%);
	} */

	#sidebar:global(.sidebar-open) {
		animation-name: slide-in;
		animation-timing-function: ease-in;
		animation-duration: 0.4s;
		animation-fill-mode: forwards;
	}

	#sidebar:global(.sidebar-closed) {
		animation-name: slide-out;
		animation-duration: 0.4s;
		animation-timing-function: ease-out;
		animation-fill-mode: forwards;
	}

	@keyframes -global-slide-in {
		from {
			transform: translateX(100%);
		}
		to {
			transform: none;
		}
	}

	@keyframes -global-slide-out {
		from {
			transform: none;
		}
		to {
			transform: translateX(100%);
		}
	}
</style>
