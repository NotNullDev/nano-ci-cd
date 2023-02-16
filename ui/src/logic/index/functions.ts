import { fetchNanoContext } from './api';
import { indexPageStore } from './store';

/**
 * Refetches the NanoContext and updates the store
 */
export async function refetchNanoContext() {
	const resp = await fetchNanoContext();
	indexPageStore.set(resp);
}
