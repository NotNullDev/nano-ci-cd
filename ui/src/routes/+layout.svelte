<script lang="ts">
	import '../app.css';
	import LinkBase from '../components/linkBase.svelte';
</script>

<div class="min-h-screen flex flex-col bg-gray-800 text-slate-200">
	<header class="p-4 flex">
		<a href="/" class="hover:bg-gray-700 p-1 cursor-pointer rounded px-2">Nano UI</a>
	</header>
	<div class="flex flex-1">
		<nav
			class="p-4 shadow-xl rounded-xl justify-between items-center flex flex-col border-r border-gray-700 mr-4 "
		>
			<div class="flex flex-col gap-1w">
				<LinkBase href="/" class="flex gap-1 items-center">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="icon icon-tabler icon-tabler-dashboard"
						width="24"
						height="24"
						viewBox="0 0 24 24"
						stroke-width="2"
						stroke="currentColor"
						fill="none"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<path stroke="none" d="M0 0h24v24H0z" fill="none" />
						<path d="M12 13m-2 0a2 2 0 1 0 4 0a2 2 0 1 0 -4 0" />
						<path d="M13.45 11.55l2.05 -2.05" />
						<path d="M6.4 20a9 9 0 1 1 11.2 0z" />
					</svg>
					<div>Dashboard</div>
				</LinkBase>
				<LinkBase href="/system" class="flex gap-1 items-center">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="icon icon-tabler icon-tabler-chart-area-line-filled"
						width="24"
						height="24"
						viewBox="0 0 24 24"
						stroke-width="2"
						stroke="currentColor"
						fill="none"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<path stroke="none" d="M0 0h24v24H0z" fill="none" />
						<path d="M4 19l4 -6l4 2l4 -5l4 4l0 5l-16 0" fill="currentColor" />
						<path d="M4 12l3 -4l4 2l5 -6l4 4" />
					</svg>
					<div>System</div>
				</LinkBase>

				<LinkBase class="flex gap-1 items-center" href="/settings">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="hover:stroke-gray-400 cursor-pointer"
						width="24"
						height="24"
						viewBox="0 0 24 24"
						stroke-width="2"
						stroke="currentColor"
						fill="none"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<path stroke="none" d="M0 0h24v24H0z" fill="none" />
						<path
							d="M10.325 4.317c.426 -1.756 2.924 -1.756 3.35 0a1.724 1.724 0 0 0 2.573 1.066c1.543 -.94 3.31 .826 2.37 2.37a1.724 1.724 0 0 0 1.065 2.572c1.756 .426 1.756 2.924 0 3.35a1.724 1.724 0 0 0 -1.066 2.573c.94 1.543 -.826 3.31 -2.37 2.37a1.724 1.724 0 0 0 -2.572 1.065c-.426 1.756 -2.924 1.756 -3.35 0a1.724 1.724 0 0 0 -2.573 -1.066c-1.543 .94 -3.31 -.826 -2.37 -2.37a1.724 1.724 0 0 0 -1.065 -2.572c-1.756 -.426 -1.756 -2.924 0 -3.35a1.724 1.724 0 0 0 1.066 -2.573c-.94 -1.543 .826 -3.31 2.37 -2.37c1 .608 2.296 .07 2.572 -1.065z"
						/>
						<path d="M12 12m-3 0a3 3 0 1 0 6 0a3 3 0 1 0 -6 0" />
					</svg>
					<div>Settings</div>
				</LinkBase>
			</div>
		</nav>
		<main class="flex flex-1 flex-col">
			<slot />
		</main>
	</div>
</div>
