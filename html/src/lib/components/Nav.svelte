<script lang="ts">
	import { activeContent } from '../stores/stores';
	import { Button, Modal } from 'flowbite-svelte';
	let goalsModal = false;
	let timelineModal = false;
	let rulesModal = false;
    import {userStore, loggedInStore} from "../stores/stores";
    

	function setActiveContent(content: string) {
		activeContent.set(content);
		console.log(activeContent);
		if ($loggedInStore) {
			document.getElementById(content)?.classList.add('card-module');
		}
	}
</script>

<nav>
	<ul class="flex row sp-between nav-ul">
		{#if loggedInStore}
			<!-- LOGGED IN -->
			<!-- # have a codejam badge -->
			<li>
				<a href={'#'} on:click|preventDefault={() => setActiveContent('profile')}>
                    <Button class="bg-slate-500 hover:bg-rose-200 text-white border-2 border-teal-300">Profile</Button>
                </a>
			</li>
			<li>
				<Button
					class="bg-rose-500 hover:bg-rose-200 text-white"
					on:click={() => (goalsModal = true)}>Goals</Button
				>
			</li>
			<li>
				<Button
					class="bg-violet-500 hover:bg-violet-200 text-white"
					on:click={() => (timelineModal = true)}>Timeline</Button
				>
			</li>
			<li>
				<Button class="bg-sky-500 hover:bg-sky-200 text-white" on:click={() => (rulesModal = true)}
					>Rules</Button
				>
			</li>
		{:else if !loggedInStore}
			<!-- NOT LOGGED IN -->
			<!-- Currently work as anchor tags -->
			<li>
				<a href="#goals"><Button class="bg-rose-500 hover:bg-rose-200 text-white">Goals</Button></a>
			</li>
			<li>
				<a href="#timeline"
					><Button class="bg-violet-500 hover:bg-violet-200 text-white">Timeline</Button></a
				>
			</li>
			<li>
				<a href="#rules"><Button class="bg-sky-500 hover:bg-sky-200 text-white">Rules</Button></a>
			</li>
		{/if}
	</ul>
</nav>

<Modal title="Codejam Goals" bind:open={goalsModal} autoclose outsideclose>
	<p class="text-base leading-relaxed text-gray-500 dark:text-gray-400 bg-rose-100">
		replace this text with goals content
	</p>
	<p class="text-base leading-relaxed text-gray-500 dark:text-gray-400">
		The European Union’s General Data Protection Regulation (G.D.P.R.) goes into effect on May 25
		and is meant to ensure a common set of data rights in the European Union. It requires
		organizations to notify users as soon as possible of high-risk data breaches that could
		personally affect them.
	</p>
	<svelte:fragment slot="footer">
		<Button on:click={() => alert('Handle "success"')}>Close</Button>
	</svelte:fragment>
</Modal>
<Modal title="Doot" bind:open={timelineModal} autoclose outsideclose>
	<p class="text-base leading-relaxed text-gray-500 dark:text-gray-400">
		With less than a month to go before the European Union enacts new consumer privacy laws for its
		citizens, companies around the world are updating their terms of service agreements to comply.
	</p>
	<p class="text-base leading-relaxed text-gray-500 dark:text-gray-400">
		The European Union’s General Data Protection Regulation (G.D.P.R.) goes into effect on May 25
		and is meant to ensure a common set of data rights in the European Union. It requires
		organizations to notify users as soon as possible of high-risk data breaches that could
		personally affect them.
	</p>
	<svelte:fragment slot="footer">
		<Button on:click={() => alert('Handle "success"')}>Close</Button>
	</svelte:fragment>
</Modal>

<Modal title="Rules" bind:open={rulesModal} autoclose outsideclose>
	<p class="text-base leading-relaxed text-gray-500 dark:text-gray-400">
		With less than a month to go before the European Union enacts new consumer privacy laws for its
		citizens, companies around the world are updating their terms of service agreements to comply.
	</p>
	<p class="text-base leading-relaxed text-gray-500 dark:text-gray-400">
		The European Union’s General Data Protection Regulation (G.D.P.R.) goes into effect on May 25
		and is meant to ensure a common set of data rights in the European Union. It requires
		organizations to notify users as soon as possible of high-risk data breaches that could
		personally affect them.
	</p>
	<svelte:fragment slot="footer">
		<Button on:click={() => alert('Handle "success"')}>Close</Button>
	</svelte:fragment>
</Modal>

<style>
	.nav-ul {
		padding: 1rem;
		justify-content: space-around;
	}
	li {
		list-style-type: None;
		gap: 20px;
		font-size: 1rem;
	}
</style>
