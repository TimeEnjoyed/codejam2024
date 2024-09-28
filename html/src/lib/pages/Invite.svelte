<script lang="ts">
	import Page from '../components/Page.svelte';
	import { Button, Card, type PweightType } from 'flowbite-svelte';
	import { onMount } from 'svelte';
	import CodeJamTeam from '../models/team';
    import { type User } from '../models/user'
	import { getTeamByInvite, joinTeam } from '../services/services';
	import { Label, Input } from 'flowbite-svelte';
	import type TeamMember from '../models/TeamMember';
	import CodeJamEvent from '../models/event';
	import { loggedInStore, userStore } from '../stores/stores';
	import DiscordIcon from '../components/DiscordIcon.svelte';
	import {location} from "svelte-spa-router";

	export let params: any; // set by svelte-spa-router
    //console.log(params) // returns Object { invitecode: "d1869a59b4fdf3" }
    console.log(params.invitecode)

	let teamData: CodeJamTeam | null = null;
	let teamMembers: TeamMember[] = [];
	let teamEvent: CodeJamEvent | null = null;
	let loading = true;
	let error: any = null;
    let teamId: string = '';

    interface TeamInfo {
        Team: CodeJamTeam
        Event: CodeJamEvent
        Members: TeamMember[]
    }

	async function loadData(invitecode: string) {
		try {
			const response = await getTeamByInvite(invitecode);
			const data: TeamInfo = await response.json();
			teamData = data.Team;

			if (teamData === null) {
				console.log("TeamData was unexpectedly null: Check the server for logs.")
				return
			}
            teamId = teamData.Id

			teamMembers = data.Members;
			teamEvent = data.Event;
		} catch (err) {
			error = 'Failed to load team data.';
			console.error(err);
		} finally {
			loading = false;
		}
	}

	$: if (params) {
		loadData(params.invitecode);
	}

	if (!params || teamId === null) {
		console.error("TeamID could not be assigned: Check the server for logs.")
	}
	

</script>

<Page>
	<Card>
		<h3>Join</h3>
		{#if $loggedInStore}
        <div class="py-4">
                <div>Hi {$userStore?.DisplayName},</div> 
                Click below to join {teamMembers[0]?.DisplayName}'s team: 
            </div>

            <Button on:click={()=>joinTeam(teamId, params.invitecode)} href="/#/teams">Join {teamData?.Name}</Button>
		{:else}
        <div class="py-4">
			Must be logged in to join a team.
        </div>
            <Button>
                <a href="/oauth/redirect?redirect={$location}">Login with Discord <DiscordIcon /></a> 
            </Button>
		{/if}
	</Card>
</Page>
