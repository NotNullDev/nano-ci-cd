<script type="ts">
	import { goto } from '$app/navigation';
	import ButtonBase from '../../components/buttonBase.svelte';
	import InputBase from '../../components/inputBase.svelte';
	import { authStore } from '../../logic/common/store';

	authStore.subscribe((val) => {
		if (val.isLoggedIn) {
			goto('/');
		}
	});
</script>

<div class="flex flex-1 items-center justify-center">
	<div class="p-12 pb-9 flex flex-col gap-2 shadow-xl shadow-orange-900 -translate-y-12" id="login">
		<h2 class="text-2xl text-center  mb-2 ">Welcome to Nano CI CD</h2>
		<InputBase placeholder="login" />
		<InputBase placeholder="password" />
		<InputBase placeholder="serverUrl" />
		<ButtonBase
			class="mt-2"
			on:click={() => {
				authStore.update((store) => {
					store.isLoggedIn = true;
					return { ...store };
				});
			}}>LOGIN</ButtonBase
		>
		<div class="text-end mt-3 text-slate-400 cursor-pointer">Forgot password?</div>
	</div>
</div>

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
