<script lang="ts">
	import clsx from 'clsx';
	import { onDestroy, onMount } from 'svelte';
	export let defaultOpen = false;
	export let className = '';

	$: open = defaultOpen ?? false;

	export let sidebar: HTMLElement | null = null;
	let sidebarOverlay: HTMLElement | null = null;
	let sidebarContainer: HTMLElement | null = null;

	export const closeSidebar = () => {
		console.log('close');

		_closeSidebar({
			target: sidebarOverlay
		});
	};

	const _closeSidebar = (e: any) => {
		if (sidebar) {
			if (e.target !== sidebarOverlay) return;
			console.log('closing...');

			sidebar.classList.remove('sidebar-open');
			sidebar.classList.add('sidebar-closed');

			sidebarOverlay?.classList.remove('bg-gray-900/50');
			sidebarOverlay?.classList.remove('pointer-events-auto');
			sidebarOverlay?.classList.add('pointer-events-none');
			sidebarOverlay?.removeEventListener('click', _closeSidebar);
			open = false;
		}
	};

	export const openSidebar = () => {
		if (sidebar) {
			sidebar.classList.remove('sidebar-closed');
			sidebar.classList.add('sidebar-open');

			sidebarOverlay?.classList.add('bg-gray-900/50');
			sidebarOverlay?.classList.add('pointer-events-auto');
			sidebarOverlay?.classList.remove('pointer-events-none');
			sidebarOverlay?.addEventListener('click', _closeSidebar);
			open = true;
		}
	};

	onMount(() => {
		if (sidebarContainer) {
			document.body.append(sidebarContainer);
			sidebar = sidebarContainer.querySelector('#sidebar') as HTMLElement;
			sidebarOverlay = sidebarContainer.querySelector('#sidebar-overlay') as HTMLElement;
		}
	});
	onDestroy(() => {
		if (sidebarContainer) {
			sidebarContainer.remove();
		}
	});
</script>

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
	<div class="relative z-10 h-full w-full pointer-events-none" id="sidebar-overlay">
		<div
			id="sidebar"
			style="transform: translateX(100%)"
			class={clsx(
				'z-20 bg-gray-800 border-l border-l-gray-900 rounded-tl-xl shadow-violet-900 shadow  absolute top-0 right-0  h-screen w-[300px]',
				className,
				{
					'sidebar-open': defaultOpen
				}
			)}
			bind:this={sidebar}
		>
			<slot />
		</div>
	</div>
</div>

<style>
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
