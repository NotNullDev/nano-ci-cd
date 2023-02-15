<script type="ts">
	import ButtonBase from '../../components/buttonBase.svelte';
	import InputBase from '../../components/inputBase.svelte';
	import { authStore } from '../../logic/common/store';
</script>

<div class="flex flex-1 items-center justify-center -translate-y-28">
	<div class="p-12 pb-9 flex flex-col gap-2 shadow-xl shadow-orange-900">
		<h2 class="text-2xl text-center  mb-2">Welcome to Nano CI CD</h2>
		<InputBase placeholder="login" />
		<InputBase placeholder="password" />
		<InputBase placeholder="serverUrl" />
		<ButtonBase
			class="mt-2"
			on:click={() => {
				authStore().update((store) => {
					store.isLoggedIn = true;
					return store;
				});
			}}>LOGIN</ButtonBase
		>
		<div class="text-end mt-3 text-slate-400 cursor-pointer">Forgot password?</div>
	</div>
</div>
