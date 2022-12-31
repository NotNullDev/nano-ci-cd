<script lang="ts">
	let area = `
    gaga
    f
    fas
    fds
    fas
    df
    dsaf
    asddf
    asda
    sdsa
    dasd
    a
    d
    asdas
    d
    as
    dsad
    asd
    asd
    asd
    as
    das
    das`;
</script>

<div class="flex flex-1 flex-col items-center p-4">
	<h1 class="text-2xl font-bold mb-6">Configuration</h1>
	<div class="text-xl mb-2">Environment variables</div>
	<textarea
		class="textarea textarea-bordered shadow-xl w-1/2 h-[300px] overflow-y-auto"
		bind:value={area}
	/>
</div>
