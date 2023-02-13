<script lang="ts">
	import '../app.css';
	import LinkBase from '../components/linkBase.svelte';
</script>

<div class="min-h-screen flex flex-col bg-gray-800 text-slate-200">
	<header class="p-4 flex">
		<a href="/" class="hover:bg-gray-700 p-1 cursor-pointer rounded px-2">Nano UI</a>
	</header>
	<div class="flex flex-1">
		<nav class="p-4 shadow-xl rounded-xl justify-between items-center flex flex-col">
			<div class="flex flex-col gal-2">
				<LinkBase
					href="/"
					on:click={() => {
						console.log('hello!');
					}}>Dashboard</LinkBase
				>
				<LinkBase
					href="/system"
					on:click={() => {
						console.log('hello!');
					}}>System</LinkBase
				>
			</div>

			<a href="/settings">
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
			</a>
		</nav>
		<main class="flex flex-1 flex-col">
			<slot />
		</main>
	</div>
</div>
