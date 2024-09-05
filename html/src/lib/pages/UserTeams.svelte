<script lang="ts">
import Page from '../components/Page.svelte';
import { Breadcrumb, BreadcrumbItem, Card} from 'flowbite-svelte';

import CodeJamTeam from '../models/team';
import CodeJamEvent from '../models/event';
import { getUserTeams } from '../services/services';
import type TeamMember from '../models/TeamMember';

export const params: Record<string, never> = {};

let teamData: CodeJamTeam | null = null;
let teamEvent: CodeJamEvent | null = null;
let teamMembers: TeamMember[] = [];
let loading: boolean = true;
let error: string | null = null;
let userTeams: CodeJamTeam[]

async function loadData() {
    try {
        const response = await getUserTeams();
        userTeams = await response.json();  // Array of teams...
    } catch (err) {
        error = `Failed to load team data: ${err}`;
        console.error(err);
    } finally {
        loading = false;
    }
}


loadData();
</script>
<!-- TODO: Make this page only accessible by team owner 
If not team owner, cannot view-->

<Page>
    <Breadcrumb solid class="mb-4 w-full max-w-screen-xl">
		<BreadcrumbItem href="/#/" home>Home</BreadcrumbItem>
		<BreadcrumbItem href="/#/teams">My Teams</BreadcrumbItem>
	</Breadcrumb>

	<Card size="md" class="w-full flex">
        <h3>Your Teams</h3>
        
		{#if loading}
			<div class="p-4">Loading...</div>
		{:else if error}
			<div class="p-4 text-red-500">{error}</div>
        {:else if userTeams === null}
            <div>
                Error, please contact admin.
            </div>
		{:else if userTeams.length === 0}
            <div>
                Looks like you don't have any teams. Go to <a href="/#/teams/browse">browse</a> teams to join one or <a href="/#/teams/create">create</a> your own team!
            </div>
        {:else}
            {#each userTeams as userTeam}
            <Card size="xl" class="flex w-full p-8 px-4 py-6 space-y-3">
				<center class="p-2">
					<h4>Team {userTeam.Name}</h4>
				</center>
                <span>
					<b>Visibility: </b>{userTeam.Visibility}
				</span>
				<span>
					<b>Technologies: </b>{userTeam.Technologies}
				</span>
				<span>
					<b>Availability: </b>{userTeam.Availability}
				</span>
				<span>
					<b>Description: </b>{userTeam.Description}
				</span>
			</Card>
            {/each}		
		{/if}
	</Card>
</Page>
