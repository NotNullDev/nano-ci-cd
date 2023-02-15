import { get } from 'svelte/store';
import { authStore } from '../../stores/authStore';

async function nanoFetch(path: string, options?: RequestInit) {
	const token = get(authStore()).token;
	const serverUrl = get(authStore()).serverUrl;
	const isLoggedIn = get(authStore()).isLoggedIn;

	if (options) {
		options.headers = {
			...options.headers,
			'nano-token': token
		};
	} else {
		options = {
			headers: {
				'nano-token': token
			}
		};
	}

	const resp = await fetch(serverUrl + path, options);

	if ((isLoggedIn && resp.status === 401) || resp.status === 403) {
		authStore().update((state) => {
			state.isLoggedIn = false;
			state.token = '';
			return state;
		});
	}

	if (!resp.ok) {
		const errMessage = await resp.json();
		return;
	}

	return resp;
}
