import { NanoContext } from '../../types/types';

export async function fetchNanoContext(): Promise<NanoContext> {
	const res = await nanoFetch('/');
	const dataRaw = await res?.json();

	const data = NanoContextSchema.parse(dataRaw);

	data.nanoConfig.globalEnvironment = base64Decode(data.nanoConfig.globalEnvironment);

	data.apps.map((app) => {
		app.envVal = base64Decode(app.envVal);
		app.buildVal = base64Decode(app.buildVal);
	});

	data.apps.sort((a, b) => b.ID - a.ID);
	return data;
}
