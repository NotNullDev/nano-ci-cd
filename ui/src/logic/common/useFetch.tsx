import { ZodError } from 'zod';

export function useFetch<T>(
	fetchFunction: () => Promise<() => T>,
	onSuccess: (data: T) => void,
	onError: (error: ZodError | any) => void,
	onSettled: () => void
) {
	return {};
}
