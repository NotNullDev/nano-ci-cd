export class NotificationManager {
	notifications: number[] = [];

	showNotifiacation(message: string, type: string) {
		console.log(message);
	}

	close(id: number) {
		this.notifications = this.notifications.filter((n) => n !== id);
	}
}
