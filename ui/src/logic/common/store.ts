import { writable, type Writable } from 'svelte/store';

export type AuthData = {
	isLoggedIn: boolean;
	token: string;
	serverUrl: string;
};

export type AuthStoreType = {
	dummy?: string;
} & Writable<AuthData>;

export function authStore(): AuthStoreType {
	const s = writable({
		isLoggedIn: false,
		token: '',
		serverUrl: 'http://localhost:8080'
	});

	return {
		...s
	};
}
