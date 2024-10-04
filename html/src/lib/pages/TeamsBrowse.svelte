<script lang="ts">

import { Avatar, Button, Card, Modal } from "flowbite-svelte";
import toast from 'svelte-french-toast';
import { location } from "svelte-spa-router";
import DiscordIcon from '../components/DiscordIcon.svelte';
import Page from "../components/Page.svelte";
import CodeJamTeam from '../models/team';
import TeamMember from '../models/TeamMember';
import { getTeams, joinPublicTeam } from "../services/services";
import { activeUserStore, loggedInStore, userStore } from '../stores/stores';
import { onMount } from "svelte";

export const params: Record<string, never> = {};

let teamData: CodeJamTeam | null = null;
let loading: boolean = true;
let error: string | null = null;
let allTeams: CodeJamTeam[] = [];
let publicTeams: CodeJamTeam[] = [];
let clickOutsideModal = false;
let avatarUrls: Record<string, string> = {};

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

async function loadData() {
    try {
        const response = await getTeams();

        allTeams = await response.json();  // Array of teams...
        console.log("allTeams: ", allTeams)
    } catch (err) {
        error = `Failed to load team data: ${err}`;
    } finally {
        loading = false;
    }
}

async function getAvatarUrl(member: TeamMember): Promise<string> {
    let ext = member.AvatarUrl.startsWith("a_") ? ".gif" : ".png";
    return `https://cdn.discordapp.com/avatars/${member.ServiceUserId}/${member.AvatarUrl}${ext}`
}

async function loadAvatarUrls() {
    let members: TeamMember[] = [];

    for (let team of allTeams) {
        members.push(...team.TeamMembers)
    }

    const promises = members.map(async member => {
        const url = await getAvatarUrl(member);
        avatarUrls[member.UserId] = url;
    });

    await Promise.all(promises);
}

let currUserId: string | undefined= $userStore?.Id

function isUserInTeam(teamMembers: TeamMember[]): boolean {
    for (let teamMember of teamMembers) { 
        if (teamMember.UserId == currUserId) {
            return true
        }
    }
    return false
}

function getTeamOwner(teamMembers: TeamMember[]): string {
    let owner = teamMembers.find(member => member.TeamRole === "owner"); 
    if (owner) {
        return owner.DisplayName
    } else {
        return "No owner found."
    }
}

onMount(() => {
    loadData();
    loadAvatarUrls();
});

$: allTeams, loadAvatarUrls();
$: publicTeams = allTeams.filter(t => t.Visibility == "public" )

function isValidTeamId(resTeamId: string | ErrorResponse): resTeamId is string {
    // Check if resTeamId is an object with a 'Severity' property indicating an error
    if (typeof resTeamId === 'object' && 'Severity' in resTeamId && resTeamId.Severity === 'ERROR') {
        return false; // Not a valid Team ID
    }
    return true; // Valid Team ID
}
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
                    <b>Owner: </b>{getTeamOwner(Team.TeamMembers)}
                </span>
                <span>
                    <b>Members: </b>
                    <div class="flex mb-5 ml-3">
                    {#each Team.TeamMembers as Member}
                        <Avatar src="{avatarUrls[Member.UserId]}" title={Member.DisplayName} stacked />
                    {/each}
                    </div>
                </span>
                
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
                
                <!-- this loops for every team separately.  -->
                {#if !$loggedInStore}
                    NOT logged in, Log Into Discord button
                    <Button on:click={()=>(clickOutsideModal=true)}>Join di</Button>
                    <Modal classBackdrop={"bg-gray-900/15 space-y-9"} bind:open={clickOutsideModal} autoclose outsideclose>
                        <p class="text-base leading-relaxed text-gray-500 dark:text-gray-400"><a href="/oauth/redirect?redirect={$location}">Login with Discord <DiscordIcon /></a>  to join a team!</p>
                    </Modal>   
                {:else if $loggedInStore}
                    {#if !isUserInTeam(Team.TeamMembers)}
                        Show join button option
                        <Button on:click={()=>joinPublicTeam(Team.Id)
                            .then((resTeamId) => {
                                if (isValidTeamId(resTeamId)) {
                                    toast.success("Successfully joined team")
                                    window.location.href = '/#/team/' + resTeamId; // redirect to team page
                                } else {
                                    toast.error("Error:" + resTeamId.Message)
                                }
                            })}>Join
                        </Button>    
                    {:else if isUserInTeam(Team.TeamMembers)}
                        <Button disabled>Already joined</Button>
                    {/if}
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
        