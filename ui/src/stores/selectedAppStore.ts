import { writable } from 'svelte/store';

export type NanoApplication = {
	name: string;
	env: string;
};

export const selectedNanoApplication = writable<NanoApplication>();

export const nanoApplications = writable<NanoApplication[]>([]);
