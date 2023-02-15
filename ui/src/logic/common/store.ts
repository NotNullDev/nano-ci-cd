import { writable } from 'svelte/store';

export type AuthData = {
	isLoggedIn: boolean;
	token: string;
	serverUrl: string;
};

export const authStore = writable<AuthData>({
	isLoggedIn: false,
	token: '',
	serverUrl: 'http://localhost:8080'
});
