import { persistentMap } from '@nanostores/persistent';
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

export const authStore = persistentMap<AuthData & LoginFormData>(
	'auth-store:',
	{
		isLoggedIn: false,
		token: '',
		password: '',
		username: '',
		serverUrl: 'http://localhost:8080'
	},
	{
		encode(value) {
			return JSON.stringify(value);
		},
		decode(value) {
			try {
				return JSON.parse(value);
			} catch (e) {
				return value;
			}
		}
	}
);
