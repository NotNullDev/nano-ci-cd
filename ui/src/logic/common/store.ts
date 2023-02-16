import { writable } from 'svelte/store';

export type AuthData = {
	isLoggedIn: boolean;
	token: string;
	serverUrl: string;
};

export type LoginFormData = {
	username: string;
	password: string;
	serverUrl: string;
};

export const authStore = writable<AuthData & LoginFormData>({
	isLoggedIn: false,
	token: '',
	password: '',
	username: '',
	serverUrl: 'http://localhost:8080'
});
