<script lang="ts">
	import Page from '../components/Page.svelte';
	import { Avatar, Breadcrumb, BreadcrumbItem, Button, Card, Input } from 'flowbite-svelte';
	import { FontAwesomeIcon } from '@fortawesome/svelte-fontawesome';
	import { faPenToSquare } from '@fortawesome/free-solid-svg-icons';
	import CodeJamTeam from '../models/team';
	import { onMount } from 'svelte';
	import { activeEventStore } from '../stores/stores';
	import { getTeamById, postTeam, removeMemberFromTeam } from '../services/services';
	import Helper from 'flowbite-svelte/Helper.svelte';
	import Radio from 'flowbite-svelte/Radio.svelte';
	import Spinner from 'flowbite-svelte/Spinner.svelte';
	import Textarea from 'flowbite-svelte/Textarea.svelte';
	import Form from '../components/Form.svelte';
	import FormField from '../components/FormField.svelte';
	import toast from 'svelte-french-toast';
	import type CodeJamEvent from '../models/event';
	import type TeamMember from '../models/TeamMember';
	import {
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow,
		TableHead,
		TableHeadCell
	} from 'flowbite-svelte';

	export let params: any; // set by svelte-spa-router

	// TODO:
	// Show each team member with a delete button next to username
	// Successfully submit the form (update query to database.)

	let teamData: CodeJamTeam | null = null;
	let teamMembers: TeamMember[] = [];
	let teamEvent: CodeJamEvent | null = null;
	let avatarUrls: Record<string, string> = {};

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

			postTeam(formData)
				.then((response) => {
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

	function removeMember(teamId: string, memberId: string) {
		if (confirm('Are you sure you want to remove this team member?')) {
			// Call API to remove member
			removeMemberFromTeam(teamId, memberId)
				.then(() => {
					// Update formData to reflect changes
					teamMembers = teamMembers.filter((member) => member.UserId !== memberId);
					toast.success('Member removed successfully');
				})
				.catch((err) => {
					toast.error('Failed to remove member');
					console.error(err);
				});
		}
	}

	async function getAvatarUrl(member: TeamMember): Promise<string> {
		let ext = member.AvatarUrl.startsWith('a_') ? '.gif' : '.png';
		return `https://cdn.discordapp.com/avatars/${member.ServiceUserId}/${member.AvatarUrl}${ext}`;
	}

	//avatarUrls[member.UserId] = url;
	async function loadAvatarUrls() {
		let members: TeamMember[] = [];
		const promises = members.map(async (member) => {
			const url = await getAvatarUrl(member);
		});
		await Promise.all(promises);
	}

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

	$: if (params) {
		loadData(params.id);
		loadAvatarUrls();
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
		{#if loading}
			{console.log(formData, 'line 107')}
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
						<Input bind:value={teamAvailability}></Input>
					</FormField>

					<!-- <MultiSelect id="multi-close" items={languages} bind:value={teamTechnologies} /> -->
					<FormField label="Your technologies:" name="TeamTechnologies">
						<Input bind:value={teamTechnologies}></Input>
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
		{#if loading}
			{console.log(formData, 'line 107')}
			<div class="p-4">Loading...</div>
		{:else if error}
			<div class="p-4 text-red-500">{error}</div>
		{:else if formData !== null}
			<h2>Team Members</h2>

			<Table>
				<TableHead>
					<TableHeadCell>Avatar</TableHeadCell>
					<TableHeadCell>Username</TableHeadCell>
					<TableHeadCell>Role</TableHeadCell>
					<TableHeadCell>
						<span class="sr-only">Delete</span>
					</TableHeadCell>
				</TableHead>
				<TableBody tableBodyClass="divide-y">
					{#each teamMembers as member}
						<TableBodyRow>
							<TableBodyCell
								><Avatar src={avatarUrls[member.Id]} title={member.DisplayName} /></TableBodyCell
							>
							<TableBodyCell>{member.DisplayName}</TableBodyCell>
							<TableBodyCell>{member.TeamRole}</TableBodyCell>
							<TableBodyCell>
								{#if member.TeamRole !== 'owner'}
									{console.log(member)}
									<Button
										on:click={() =>
											teamData?.Id && member.Id && removeMember(teamData.Id, member.Id)}
										class="btn-remove text-red-500 hover:text-red-700"
										color="light"
									>
										Remove
									</Button>
								{/if}
							</TableBodyCell>
						</TableBodyRow>
					{/each}
				</TableBody>
			</Table>
		{/if}
	</Card>
</Page>
