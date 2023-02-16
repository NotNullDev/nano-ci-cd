import { writable } from 'svelte/store';
import type { NanoContext } from '../../types/types';

export const indexPageStore = writable<NanoContext>({
	apps: [],
	nanoConfig: {
		globalEnvironment: '',
		token: ''
	},
	buildingAppId: 0
});
