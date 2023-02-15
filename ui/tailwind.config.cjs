/** @type {import('tailwindcss').Config} */
module.exports = {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	theme: {
		extend: {
			colors: ({ theme, colors }) => ({
				...colors,
				primary: colors.gray[800],
				accent: colors.orange[800],
				accentHover: colors.orange[700],
				danger: colors.red[900]
			})
		}
	},
	plugins: []
};
