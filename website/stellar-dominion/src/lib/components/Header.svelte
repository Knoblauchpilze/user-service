<script lang="ts">
	import { FlexContainer } from '$lib/components';
</script>

<FlexContainer
	vertical={false}
	bgColor={'bg-black'}
	styling={'w-full h-1/10 fixed top-0 left-0 text-white'}
>
	<slot />
</FlexContainer>
