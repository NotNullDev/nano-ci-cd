/// <reference types="svelte" />

interface ButtonBasAttributes<T> extends HTMLButtonAttributes<T> {
	// Add your custom properties and methods here
	variant: 'primary' | 'danger';
}
