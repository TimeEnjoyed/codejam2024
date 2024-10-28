<script lang="ts">
	import { Breadcrumb, BreadcrumbItem, Card } from 'flowbite-svelte';
	import { onMount } from 'svelte';
	import Page from '../components/Page.svelte';
	import TeamMember from '../models/TeamMember';
	import CodeJamEvent from '../models/event';
	import CodeJamTeam from '../models/team';
	import { getTeamById } from '../services/services';
	import toast from 'svelte-french-toast';
    import {userStore} from "./../stores/stores";
	interface Params {
		// This is what the params is because you pass an id with type of string.
		id: string;
	}

	export let params: Params;
	const urlParams = new URLSearchParams(document.location.search);
	console.log('urlParams: ', urlParams);

	console.log(params);
	let teamData: CodeJamTeam | null = null;
	let teamMembers: TeamMember[] = [];
	let teamEvent: CodeJamEvent | null = null;
	let loading: boolean = true;
	let error: string | null = null;
	let teamCreated: boolean = false; // Reactivate variable to track if team was just created

    // to display the owner the team
	function getTeamOwner(teamMembers: TeamMember[]): string {
		let owner = teamMembers.find((member) => member.TeamRole === 'owner');
		if (owner) {
			return owner.DisplayName;
		} else {
			return 'No owner found.';
		}
	}

	async function loadData(id: string) {
		try {
			const response = await getTeamById(id);
			const data = await response.json();
			teamData = data.Team;
			teamMembers = data.TeamMembers;
			console.log('teamMembers in loadData: ', teamMembers);
			teamEvent = data.Event;
		} catch (err) {
			error = `Failed to load team data: ${err}`;
		} finally {
			loading = false;
		}
	}

	$: if (params) {
		loadData(params.id); // See this doesn't error cause it expects an id: You can run to see I guess
	}

	onMount(() => {
		const hash = window.location.hash;
		const queryString = hash.split('?')[1];
		if (queryString) {
			const urlParams = new URLSearchParams(queryString);
			teamCreated = urlParams.get('teamCreated') === 'true';
		}
		console.log('onMount teamCreated: ', teamCreated);
	});

	let url: string = '';

	$: if (teamData?.Id) {
		url = `localhost:8080/#/team/edit/${teamData.Id}`;
	}

	function copyToClipboard(): void {
		navigator.clipboard
			.writeText(url)
			.then(() => {
				toast.success('Copied to clipboard');
			})
			.catch((err: Error) => {
				toast.error(`Failed to copy: ${err}`);
			});
	}
</script>

<!-- TODO: Create Owner, Member, and Public View -->

<Page>
	<Breadcrumb solid class="mb-4 w-full max-w-screen-xl">
		<BreadcrumbItem href="/#/" home>Home</BreadcrumbItem>
		<BreadcrumbItem>Team {teamData?.Name}</BreadcrumbItem>
	</Breadcrumb>

	<Card size="md" class="w-full flex">
		{#if loading}
			<div class="p-4">Loading...</div>
		{:else if error}
			<div class="p-4 text-red-500">{error}</div>
		{:else if teamData !== null}
			{#if teamCreated}
				<div class="p-4">
					{teamMembers[0]?.DisplayName}, your team has been successfully created!
				</div>
			{/if}
			<h3 class="p-4">{teamEvent?.Title}</h3>
			<Card size="xl" class="flex w-full p-8 px-4 py-6 space-y-3">
				<center class="p-2">
					<h4>Team {teamData.Name}</h4>
                    {#if $userStore?.DisplayName == getTeamOwner(teamMembers)}

					<span><small>(<a href="/#/team/edit/{teamData.Id}">edit</a>)</small></span>
                    {/if}
				</center>
                <span>
                    <b>Owner: </b>{getTeamOwner(teamMembers)}
                </span>
				<span>
					<b>Team Members: </b>
					{#each teamMembers as member}
						<ul>
							<li>{member.DisplayName}</li>
						</ul>
					{/each}
				</span>
				<span>
					<b>Visibility: </b>{teamData.Visibility}
				</span>
				<span>
					<b>Technologies: </b>{teamData.Technologies}
				</span>
				<span>
					<b>Availability: </b>{teamData.Availability}
				</span>
				<span>
					<b>Description: </b>{teamData.Description}
				</span>
				<span>
					<b>Invite Link: </b>
				</span>
				<div>
					<textarea
						class="border border-slate-300 bg-white text-gray-400 rounded-md w-full resize-none"
						bind:value={url}
						readonly
					></textarea>
					<button
						class="my-2 border border-slate-300 bg-white text-gray-400 p-2"
						on:click={copyToClipboard}>Copy Text</button
					>
				</div>
			</Card>
		{/if}
	</Card>
</Page>
