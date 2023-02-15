import { nanoFetch } from '../common/api';

export async function login(username: string, password: string) {
	const res = await nanoFetch('/login', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({
			username,
			password
		})
	});

	if (!res?.ok) {
		throw new Error('Login failed');
	}

	let data = (await res?.text()) as string;
	// data comes in format "<token>", so we need to remove the quotes
	data = data.slice(1);
	data = data.slice(0, data.length - 2);
	return data;
}
