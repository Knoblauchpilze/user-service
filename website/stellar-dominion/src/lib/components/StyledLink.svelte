<!-- https://svelte.dev/docs/svelte-components -->
<script lang="ts">
	export let text: string = 'Home';
	export let link: string = '/';
	export let showAsButton: boolean = false;
	export let disabled: boolean = false;
</script>

<!-- https://tailwindcss.com/docs/padding -->
<!-- https://tailwindcss.com/docs/hover-focus-and-other-states -->
{#if showAsButton}
	{#if disabled}
		<a href={link} class="px-12 py-4 rounded-[8px] bg-primary-selected text-secondary">
			{text}
		</a>
	{:else}
		<a
			href={link}
			class="px-12 py-4 rounded-[8px] bg-primary hover:bg-primary-hover text-secondary"
		>
			{text}
		</a>
	{/if}
{:else}
	<a class="text-secondary hover:text-secondary-hover" href={link}>{text}</a>
{/if}
