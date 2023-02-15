import { get } from 'svelte/store';
import { authStore, nanoContextStore } from '../stores/authStore';
import { App, AppLogsType, AppLogsTypeSchema } from '../types/types';

export async function resetToken(): Promise<string> {
	const res = await nanoFetch('/reset-token', {
		method: 'POST'
	});
	let data = (await res?.text()) as string;

	// data comes in format "<token>", so we need to remove the quotes
	data = data.slice(1);
	data = data.slice(0, data.length - 2);

	return data;
}

export async function clearBuilds(): Promise<void> {
	const res = await nanoFetch('/clear-builds', {
		method: 'GET'
	});
}

export async function dockerSystemPrune(): Promise<void> {
	const res = await nanoFetch('/docker-system-prune', {
		method: 'GET'
	});
}

export async function updateGlobalEnv(updatedEnv: string): Promise<string> {
	const res = await nanoFetch('/update-global-env', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/text'
		},
		body: updatedEnv
	});
	const data = (await res?.text()) as string;
	return base64Decode(data);
}

export async function createApp(name: string) {
	const res = await nanoFetch('/create-app', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({
			appName: name
		})
	});
	const data = (await res?.json()) as App;
	return data;
}

export async function updateApp(app: App) {
	const res = await nanoFetch('/update-app', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify(app)
	});
	const data = (await res?.json()) as App;
	return data;
}

export async function deleteApp(appId: number) {
	const res = await nanoFetch('/delete-app?id=' + appId, {
		method: 'DELETE'
	});
	const data = (await res?.text()) as string;

	return Number(data);
}

// https://github.com/Azure/fetch-event-source
export async function runBuild(appName: string) {
	const res = await nanoFetch('/build?appName=' + appName, {
		method: 'POST',
		headers: {
			Authorization: get(nanoContextStore()).nanoConfig.token ?? ''
		}
	});

	const data = (await res?.text()) as string;
	return data;
}

export async function updateUser(username: string, password: string, repeatPassword: string) {
	if (password !== repeatPassword) {
		console.error('Passwords do not match');
		return;
	}

	const res = await nanoFetch('/update-user', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({
			username: username,
			password: password
		})
	});
	const data = (await res?.text()) as string;
	return data;
}

export async function login(username: string, password: string) {
	const res = await nanoFetch('/login', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({
			username,
			password
		})
	});

	if (!res?.ok) {
		throw new Error('Login failed');
	}

	let data = (await res?.text()) as string;
	// data comes in format "<token>", so we need to remove the quotes
	data = data.slice(1);
	data = data.slice(0, data.length - 2);
	return data;
}

export async function logout() {
	authStore().update((store) => {
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

export async function fetchLogs(appId: number, limit = 1): Promise<AppLogsType> {
	const res = await nanoFetch('/logs?appId=' + appId + '&limit=' + limit);

	const data = AppLogsTypeSchema.parse(await res?.json());

	return data;
}

function base64Decode(str: string) {
	// return Buffer.from(str, 'base64').toString('ascii');
	return btoa(str);
	// return Buffer.from(str, 'base64').toString('ascii');
}
s;
function nanoFetch(arg0: string, arg1: { method: string }) {
	throw new Error('Function not implemented.');
}
