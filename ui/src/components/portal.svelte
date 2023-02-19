<script lang="ts">
	console.log('hi');

	function createPortal(node: HTMLElement) {
		const body = document.querySelector('body');

		body?.appendChild(node);

		return {
			destroy() {
				node.remove();
			}
		};
	}
</script>

<div use:createPortal {...$$restProps}>
	<slot />
</div>
