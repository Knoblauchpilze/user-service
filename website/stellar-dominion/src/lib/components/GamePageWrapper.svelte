<script lang="ts">
	import '$styles/app.css';
	import { FlexContainer, Header, StyledText } from '$lib/components';

	import { type UiResource } from '$lib/game/resources';
	import { floorToInteger, toFlooredShortString } from '$lib/displayUtils';

	// https://svelte.dev/blog/zero-config-type-safety
	export let universeName: string;
	export let planetName: string;
	export let playerName: string;

	export let resources: UiResource[];

	function resourceTextColor(resource: UiResource): string {
		if (resource.amount < resource.storage) {
			return 'text-white';
		}
		return 'text-disabled';
	}
	function productionTextColor(resource: UiResource): string {
		if (resource.amount < resource.storage) {
			return 'text-enabled';
		}
		return 'text-disabled';
	}
</script>

<FlexContainer>
	<!-- https://stackoverflow.com/questions/67852559/pass-svelte-component-as-props -->
	<Header>
		<StyledText text={universeName} textColor="text-white" />
		<StyledText text={playerName} textColor="text-white" />
		<StyledText text={planetName} textColor="text-white" />
		<form method="POST" action="?/logout">
			<button class="hover:underline">Logout</button>
		</form>
	</Header>

	<FlexContainer align={'stretch'} styling={'w-full'}>
		<FlexContainer vertical={false} extensible={false} bgColor={'bg-black'} styling={'w-full'}>
			{#each resources as resource}
				<FlexContainer vertical={false} extensible={false} styling={'space-between'}>
					<StyledText text="{resource.name}:" textColor="text-white" />
					<StyledText
						text={floorToInteger(resource.amount).toString()}
						textColor={resourceTextColor(resource)}
						styling="px-1"
					/>
					<StyledText
						text="(+{floorToInteger(resource.production)}/h)"
						textColor={productionTextColor(resource)}
						styling="pr-1"
					/>
					<StyledText
						text="(storage: {toFlooredShortString(resource.storage)})"
						textColor="text-white"
					/>
				</FlexContainer>
			{/each}
		</FlexContainer>

		<slot />
	</FlexContainer>
</FlexContainer>
