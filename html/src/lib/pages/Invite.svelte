<script lang="ts">
	import { Button, Card } from 'flowbite-svelte';
	import toast from 'svelte-french-toast';
	import { location } from 'svelte-spa-router';
	import DiscordIcon from '../components/DiscordIcon.svelte';
	import Page from '../components/Page.svelte';
	import CodeJamEvent from '../models/event';
	import CodeJamTeam from '../models/team';
	import TeamMember from '../models/TeamMember';
	import { getTeamByInvite, joinPublicTeam, joinTeam } from '../services/services';
	import { activeUserStore, loggedInStore, userStore } from '../stores/stores';

	export let params: any; // set by svelte-spa-router
	//console.log(params) // returns Object { invitecode: "d1869a59b4fdf3" }

    let teamData: CodeJamTeam | null = null;
    let teamMembers: TeamMember[] = [];
    let teamEvent: CodeJamEvent | null = null;
    let loading = true;
    let error: any = null;
    let teamId: string = '';
    $: currUserId = $userStore?.Id;

    interface ErrorResponse {
		Severity: string;
		Detail?: string;
		Code: string;
		Message: string;
		Hint?: string;
		Position?: number;
		InternalPosition?: number;
		InternalQuery?: string;
		Where?: string;
		SchemaName?: string;
	}

    interface TeamInfo {
		Team: CodeJamTeam;
		Event: CodeJamEvent;
		Members: TeamMember[];
	}

	async function loadData(invitecode: string) {
		try {
			const response = await getTeamByInvite(invitecode);
			const data = await response.json();
			teamData = data.Team;
			teamId = data.Team.Id;
			teamMembers = data.TeamMembers;
			teamEvent = data.Event;
		} catch (err) {
			error = 'Failed to load team data.';
			console.error(err);
		} finally {
			loading = false;
		}
	}

     // to tell if join button works or not:
    function isUserInTeam(teamMembers: TeamMember[]): boolean {     
        for (let teamMember of teamMembers) {
            if (teamMember.Id == currUserId) {
                console.log("teamMember and curr userid: ", teamMember.UserId, "==", currUserId);
                console.log("user ALREADY in team");
                return true;
            }
        }
        return false
    }


    function isValidTeamId(resTeamId: string | ErrorResponse): resTeamId is string {
		// Check if resTeamId is an object with a 'Severity' property indicating an error
		if (
			typeof resTeamId === 'object' &&
			'Severity' in resTeamId &&
			resTeamId.Severity === 'ERROR'
		) {
			return false; // Not a valid Team ID
		}
		return true; // Valid Team ID
	}

	$: if (params) {
		loadData(params.invitecode);
	}

	if (!params || teamId === null) {
		console.error('TeamID could not be assigned: Check the server for logs.');
	}
</script>

<Page>
	<Card>
		{#if loading}
			<p>Loading...</p>
		{:else if error}
			<p>{error}</p>
		{:else}
			<h3>Join {teamData?.Name}</h3>
			{#if $loggedInStore}
				<div class="py-4"> 
                    {#if !isUserInTeam(teamMembers)}
                        <div>Hi {$userStore?.DisplayName},</div>
                        <p>Click below to join {teamMembers[0]?.DisplayName}'s team:</p>

                        <Button class="my-5 text-white hover:text-pink-100"
                            on:click={() => {
                                joinTeam(teamId, params.invitecode);
                                window.location.href = `/#/team/${teamId}`;
                                toast.success(`Successfully joined ${teamData?.Name}`);
                            }}
                        >
                            Join {teamData?.Name}
                        </Button>
					{:else if isUserInTeam(teamMembers)}
                        <p>Hi {$userStore?.DisplayName}, Looks like you're already in this team! Go <a href="/#/team/${teamId}">here</a> to view it.</p>
						<Button class="my-5" disabled>Already joined</Button>
					{/if}
				</div>
			{:else}
				<div class="py-4">
					<p>Must be logged in to join a team.</p>
					<Button class="my-5">
						<a class="text-white hover:text-pink-100" href="/oauth/redirect?redirect={$location}">Login with Discord <DiscordIcon /></a>
					</Button>
				</div>
			{/if}
		{/if}
	</Card>
</Page>
