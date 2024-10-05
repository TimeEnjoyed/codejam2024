<script lang="ts">
	import Page from '../components/Page.svelte';
	import { Breadcrumb, BreadcrumbItem, Button, Card, Input } from 'flowbite-svelte';
	import { FontAwesomeIcon } from '@fortawesome/svelte-fontawesome';
	import { faPenToSquare } from '@fortawesome/free-solid-svg-icons';
	import CodeJamTeam from '../models/team';
	import { onMount } from 'svelte';
	import { activeEventStore } from '../stores/stores';
	import { getTeamById, postTeam } from '../services/services';
	import Helper from 'flowbite-svelte/Helper.svelte';
	import Radio from 'flowbite-svelte/Radio.svelte';
	import Spinner from 'flowbite-svelte/Spinner.svelte';
	import Textarea from 'flowbite-svelte/Textarea.svelte';
	import Form from '../components/Form.svelte';
	import FormField from '../components/FormField.svelte';
	import toast from 'svelte-french-toast';
	import type CodeJamEvent from '../models/event';
	import type TeamMember from '../models/TeamMember';

	export let params: any; // set by svelte-spa-router

    // TODO: 
    // Show each team member with a delete button next to username
    // Successfully submit the form (update query to database.)
    
    let teamData: CodeJamTeam | null = null;
	let teamMembers: TeamMember[] = [];
	let teamEvent: CodeJamEvent | null = null;

    let teamName: string = '';
    let teamVisibility: string = '';
    let teamAvailability: string = '';
    let teamTechnologies: string = '';
    let teamDescription: string = '';
    
    let loading: boolean = true; 
	let formData: CodeJamTeam | null = null;
	let isSaving: boolean = false;
	let teamCreated: boolean = false;
    let error: string | null = null; 

	let clearErrors: () => {};
	let parseResponse: (response: object) => {};

	function saveForm() {
		if (formData !== null) {
			isSaving = true;
			clearErrors();

			formData.EventId = $activeEventStore?.Id || '';
			// Step 1: Post Team Data API
			postTeam(formData)
				.then((response) => {
					// parseResponse(response);
					// const url = new URL(response.url);
					// const pathSegments = url.pathname.split('/');
					// const teamId = pathSegments[pathSegments.length - 1];
					response
						.json()
						.then((data) => {
							// Team creation successful, letting svelte page know:
							teamCreated = true;
							// Stepp 1: GET team info
							// this uses routes.ts --> MyTeam.svelte page
							window.location.href = `/#/team/${data.id}`;
							toast.success("You've successfully edited a team");
							isSaving = false;
						})
						.catch(() => {
							isSaving = false;
						});
				})
				.catch((err) => {
					console.error('Error saving event', err);
					isSaving = false;
				});
		}
	}

	// METHOD 1: OnMount
	// this returns the team form with prefilled information
	// removed onMount because there's no HTML required to preload.
	// If the fetch-call (go-backedn) takes too long you will get race-issues sometimes. when use bind or interact with dom. TL;DR use onMOUNT to load data alaways.
	// onMount(() => {
	// 	if (params) {
	// 		getTeamById(params.id).then((response) => {
	// 			response.json().then((data) => {
	// 				formData = data as CodeJamTeam;
	// 			});
	// 		});
	// 	}
	// });

    async function loadData(id: string) {
        try {
            getTeamById(params.id).then((response) => {
				response.json().then((data) => {
					formData = data as CodeJamTeam;
                    teamData = data.Team;
			        teamMembers = data.Members;
			        teamEvent = data.Event;
                    teamName = data.Team.Name;
                    teamVisibility = data.Team.Visibility;
                    teamAvailability = data.Team.Availability;
                    teamTechnologies = data.Team.Technologies;
                    teamDescription = data.Team.Description;
				});
			});
        } catch (err) {
            error = `Failed to load team data: ${err}`; 
        } finally {
            loading = false;
        }
    }

    $: if(params) {
        loadData(params.id)
    }
</script>

<Page>
	<Breadcrumb solid class="mb-4 w-full max-w-screen-xl">
		<BreadcrumbItem href="/#/" home>Home</BreadcrumbItem>
		<BreadcrumbItem href="/#/team">Team Options</BreadcrumbItem>
		<BreadcrumbItem>Edit Team</BreadcrumbItem>
	</Breadcrumb>
	<Card size="xl" class="w-full">
		<div class="flex flex-row gap-8 my-8"></div>
		Edit team form here:
		{#if loading}
            {console.log(formData, "line 107")}
			<div class="p-4">Loading...</div>
		{:else if error}
			<div class="p-4 text-red-500">{error}</div>
        {:else if formData !== null}

            <div class="flex flex-col gap-8 my-8">
                <Form bind:clearErrors bind:parseResponse>
                    <FormField label="Team Name:" name="TeamName">
                        <Input bind:value={teamName}></Input>
                    </FormField>
                    <div>
                        <Radio name="team-type" bind:group={teamVisibility} value="public">Public Team</Radio>
                        <Helper class="ml-6 ">(If you want your team to be searchable.)</Helper>
                    </div>
                    <div>
                        <Radio name="team-type" bind:group={teamVisibility} value="private">Private Team</Radio>
                        <Helper class="ml-6">(Your team will be invite only)</Helper>
                    </div>

                    <FormField label="Your general availability:" name="TeamAvailability">
                        <Input bind:value={teamAvailability}>
                        </Input>
                    </FormField>

                    <!-- <MultiSelect id="multi-close" items={languages} bind:value={teamTechnologies} /> -->
                    <FormField label="Your technologies:" name="TeamTechnologies">
                        <Input
                            bind:value={teamTechnologies}
                        ></Input>
                    </FormField>

                    <FormField label="What do you want out of this team?" name="Description">
                        <!-- <Label for="aboutTextArea">What do you want out of this team?</Label> -->
                        <Textarea bind:value={teamDescription} />
                    </FormField>
                </Form>

                <Button on:click={saveForm} disabled={isSaving}>
                    {#if isSaving}
                        <Spinner />
                    {:else}
                        Save
                    {/if}
                </Button>
            </div>
        {:else}
            <Spinner />
        {/if}

	</Card>
</Page>
