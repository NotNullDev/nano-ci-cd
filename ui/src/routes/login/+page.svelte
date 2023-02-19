<script type="ts">
	import { goto } from '$app/navigation';
	import { tick } from 'svelte';
	import ButtonBase from '../../components/buttonBase.svelte';
	import InputBase from '../../components/inputBase.svelte';
	import { authStore } from '../../logic/common/store';
	import { login } from '../../logic/logic/api';

	authStore.subscribe((val) => {
		if (val.isLoggedIn) {
			goto('/');
		}
	});

	let username = $authStore.username;
	let password = '';
	let serverUrl = $authStore.serverUrl;
</script>

<div class="flex flex-1 items-center justify-center">
	<form
		on:submit={(e) => e.preventDefault()}
		class="p-12 pb-9 flex flex-col gap-2 shadow-xl shadow-orange-900 -translate-y-12"
		id="login"
	>
		<h2 class="text-2xl text-center  mb-2 ">Welcome to Nano CI CD</h2>
		<InputBase placeholder="login" bind:val={username} />
		<InputBase placeholder="password" bind:val={password} />
		<InputBase placeholder="serverUrl" bind:val={serverUrl} />
		<ButtonBase
			class="mt-2"
			on:click={async () => {
				authStore.setKey('serverUrl', serverUrl);

				await tick();

				const token = await login(username, password);

				authStore.set({
					...authStore.get(),
					username: username,
					serverUrl: serverUrl,
					token: token,
					isLoggedIn: true
				});
			}}>LOGIN</ButtonBase
		>
		<div class="text-end mt-3 text-slate-400 cursor-pointer">Forgot password?</div>
	</form>
</div>

<svelte:head>
	<title>Login</title>
</svelte:head>

<style>
	#login {
		animation-name: login;
		animation-duration: 0.5s;
		animation-fill-mode: forwards;
	}

	@keyframes login {
		0% {
			opacity: 0;
			transform: translateY(100px);
		}
		100% {
			transform: translateY(-6rem);
		}
	}
</style>
