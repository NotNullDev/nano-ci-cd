import { get } from 'svelte/store';
import { AppLogsTypeSchema, type App, type AppLogsType } from '../../types/types';
import { nanoFetch } from '../common/api';
import { indexPageStore } from '../index/store';

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
