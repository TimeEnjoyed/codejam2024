<script lang="ts">

import Page from "../components/Page.svelte";
import {Button, Card, Modal} from "flowbite-svelte";
import {getTeams, joinPublicTeam} from "../services/services";
import { loggedInStore, userStore } from '../stores/stores';
import CodeJamTeam from '../models/team';
import type TeamMember from '../models/TeamMember';
import DiscordIcon from '../components/DiscordIcon.svelte';
import {location} from "svelte-spa-router";
export const params: Record<string, never> = {};


let teamData: CodeJamTeam | null = null;
let teamMembers: TeamMember[] = [];
let loading: boolean = true;
let error: string | null = null;
let allTeams: CodeJamTeam[] = [];
let publicTeams: CodeJamTeam[] = [];
let clickOutsideModal = false;

async function loadData() {
    try {
        const response = await getTeams();
        allTeams = await response.json();  // Array of teams...
    } catch (err) {
        error = `Failed to load team data: ${err}`;
    } finally {
        loading = false;
    }
}

loadData();

$: publicTeams = allTeams.filter(t => t.Visibility == "public" )

</script>
    
<Page>
    <Card size="md" class="w-full flex">
        <h3>Browse All Teams</h3>


		{#if loading}
			<div class="p-4">Loading...</div>
		{:else if error}
			<div class="p-4 text-red-500">{error}</div>
        {:else}
            {#each publicTeams as Team}
            <Card size="xl" class="flex w-full p-8 px-4 py-z gap-3">

				<center class="p-2">
					<h4>Team {Team.Name}</h4>
                </center>
                <span>
					<b>Visibility: </b>{Team.Visibility}
				</span>
				<span>
					<b>Technologies: </b>{Team.Technologies}
				</span>
				<span>
					<b>Availability: </b>{Team.Availability}
				</span>
				<span>
					<b>Description: </b>{Team.Description}
				</span>

                // only shows if user isnt a member
                {#if $loggedInStore}
                <Button on:click={()=>joinPublicTeam(Team.Id)}>Join</Button>
                {:else}
                <Button on:click={()=>(clickOutsideModal=true)}>Join</Button>
                <Modal classBackdrop={"bg-gray-900/15 space-y-9"} bind:open={clickOutsideModal} autoclose outsideclose>
                    <p class="text-base leading-relaxed text-gray-500 dark:text-gray-400"><a href="/oauth/redirect?redirect={$location}">Login with Discord <DiscordIcon /></a>  to join a team!</p>
                </Modal>
                {/if}
			</Card>
            {:else}
                <div>
                    Looks like you don't have any teams! Go to <a href="/#/">browse teams</a> to join one!
                </div>
        {/each}
			
		{/if}
	</Card>

</Page>
    