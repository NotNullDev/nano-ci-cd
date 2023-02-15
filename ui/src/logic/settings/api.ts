import { nanoFetch } from '../common/api';

export async function updateUser(username: string, password: string, repeatPassword: string) {
	if (password !== repeatPassword) {
		console.error('Passwords do not match');
		return;
	}

	const res = await nanoFetch('/update-user', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({
			username: username,
			password: password
		})
	});
	const data = (await res?.text()) as string;
	return data;
}
