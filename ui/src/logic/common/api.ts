import { get } from 'svelte/store';
import { authStore } from './store';

export async function nanoFetch(path: string, options?: RequestInit) {
	const token = get(authStore).token;
	const isLoggedIn = get(authStore).isLoggedIn;
	const serverUrl = get(authStore).serverUrl;

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

	const fetchUrl = serverUrl + path;
	console.log('fetchUrl', fetchUrl);

	options = {
		...options,
		headers: {
			...options.headers,
			'Content-Type': 'application/json'
		}
	};

	const resp = await fetch(fetchUrl, options);
	console.log('resp', resp);

	if ((isLoggedIn && resp.status === 401) || resp.status === 403) {
		authStore.update((state) => {
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

export async function logout() {
	authStore.update((store) => {
		store.isLoggedIn = false;
		store.token = '';
		return store;
	});
}

export async function resetGlobalBuildStatus() {
	const res = await nanoFetch('/reset-global-build-status', {});
	const data = (await res?.text()) as string;
	return data;
}
