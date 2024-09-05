<script lang="ts">
	import { Breadcrumb, BreadcrumbItem, Card } from 'flowbite-svelte';
	import Page from '../components/Page.svelte';
	import type TeamMember from '../models/TeamMember';
	import CodeJamEvent from '../models/event';
	import CodeJamTeam from '../models/team';
	import { getTeamById } from '../services/services';

	interface Params {  // This is what the params is because you pass an id with type of string.
		id: string
	}
	export let params: Params;

	let teamData: CodeJamTeam | null = null;
	let teamMembers: TeamMember[] = [];
	let teamEvent: CodeJamEvent | null = null;
	let loading: boolean = true; 
	let error: string | null = null; 
    let teamCreated: boolean = false; // Reactivate variable to track if team was just created

	async function loadData(id: string) {
		try {
			const response = await getTeamById(id);
			const data = await response.json();
			teamData = data.Team;
			teamMembers = data.Members;
			teamEvent = data.Event;
		} catch (err) {
			error = `Failed to load team data: ${err}`; 
		} finally {
			loading = false;
		}
	}

	$: if (params) {
		loadData(params.id);  // See this doesn't error cause it expects an id: You can run to see I guess
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
			<div class="p-4">{teamMembers[0]?.DisplayName}, your team has been successfully created!</div>
            {/if}
			<h3 class="p-4">{teamEvent?.Title}</h3>
			<Card size="xl" class="flex w-full p-8 px-4 py-6 space-y-3">
				<center class="p-2">
					<h4>Team {teamData.Name}</h4>
					<span><small>(edit)</small></span>
				</center>
				<span>
					<b>Team Members: </b>
					<ul>
						{#each teamMembers as member}
							<li>{member.DisplayName}</li>
						{/each}
					</ul>
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
				<p>Invite Link: <a href="/#/team/invite/{teamData.InviteCode}">here </a></p>
			</Card>
		{/if}
	</Card>
</Page>
