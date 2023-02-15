import { NanoContextSchema, type App, type NanoContext } from '../../types/types';
import { nanoFetch } from '../common/api';

export async function fetchNanoContext(): Promise<NanoContext> {
	const res = await nanoFetch('/');
	const dataRaw = await res?.json();

	console.log(dataRaw);

	const data = NanoContextSchema.parse(dataRaw);

	data.nanoConfig.globalEnvironment = btoa(data.nanoConfig.globalEnvironment);

	data.apps.map((app) => {
		app.envVal = btoa(app.envVal);
		app.buildVal = btoa(app.buildVal);
	});

	data.apps.sort((a, b) => b.ID - a.ID);
	return data;
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
	return btoa(data);
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
