import dayjs from 'dayjs';

export function base64Decode(str: string) {
	return btoa(str);
}

export class NanoUtils {
	static formatDate(dateToFormat: string) {
		return dayjs(dateToFormat).format('DD.MM.YYYY HH:mm:ss');
	}
}
