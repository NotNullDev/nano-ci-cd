import { nanoFetch } from '../common/api';

export async function resetToken(): Promise<string> {
	const res = await nanoFetch('/reset-token', {
		method: 'POST'
	});
	let data = (await res?.text()) as string;

	// data comes in format "<token>", so we need to remove the quotes
	data = data.slice(1);
	data = data.slice(0, data.length - 2);

	return data;
}

export async function clearBuilds(): Promise<void> {
	const res = await nanoFetch('/clear-builds', {
		method: 'GET'
	});
}

export async function dockerSystemPrune(): Promise<void> {
	const res = await nanoFetch('/docker-system-prune', {
		method: 'GET'
	});
}
