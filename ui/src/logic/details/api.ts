import { get } from 'svelte/store';
import {
	AppLogsTypeSchema,
	BuildMetadata,
	NanoBuildSchema,
	type App,
	type AppLogsType
} from '../../types/types';
import { nanoFetch } from '../common/api';
import { indexPageStore } from '../index/store';
import { detailsPageStore } from './store';

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
			Authorization: get(indexPageStore).nanoConfig.token ?? ''
		}
	});

	const data = (await res?.text()) as string;
	return data;
}

export async function fetchLogs(appId: number, limit = 1): Promise<AppLogsType> {
	const res = await nanoFetch('/logs?appId=' + appId + '&limit=' + limit);

	const data = AppLogsTypeSchema.parse(await res?.json());

	return data;
}

export async function getBuild(buildId: number) {
	const res = await nanoFetch('/build?buildId=' + buildId);

	const data = await res?.json();

	const validatedData = NanoBuildSchema.parse(data);

	return validatedData;
}

export async function useFetchBuild(buildId: number) {
	const data = await getBuild(buildId);
	detailsPageStore.setKey('currentBuild', data);
}

export async function fetchBuilds(): Promise<BuildMetadata[]> {
	const res = await nanoFetch('/available-builds-metadata');

	const data = await res?.json();

	const validatedData = data.map((d: any) => BuildMetadata.parse(d)) as BuildMetadata[];

	return validatedData;
}

export async function useFetchBuilds() {
	const data = await fetchBuilds();
	console.log(data);
	detailsPageStore.setKey('availableBuilds', data);
}
