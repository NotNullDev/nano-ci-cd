import { writable } from 'svelte/store';
import { fetchNanoContext } from '../api/api';
import { NanoContext } from '../types/types';

type AuthStoreType = {
	isLoggedIn: boolean;
	token: string;
	serverUrl: string;
};

export const authStore = () => {
	const { set, subscribe, update } = writable<AuthStoreType>({
		isLoggedIn: false,
		token: '',
		serverUrl: 'http://localhost:8080'
	});

	return {
		subscribe,
		set,
		update
	};
};

export async function refetchNanoContext() {
	const resp = await fetchNanoContext();

	if (!resp || !resp) {
		throw new Error('Failed to fetch nano context');
	}
}

export function nanoContextStore() {
	const { set, subscribe, update } = writable<NanoContext>({
		apps: [],
		nanoConfig: {
			globalEnvironment: '',
			token: ''
		},
		buildingAppId: 0
	});

	return {
		subscribe,
		set,
		update
	};
}
