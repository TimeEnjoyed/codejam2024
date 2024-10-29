<script lang="ts">
	import Page from '../components/Page.svelte';
	import { Avatar, Breadcrumb, BreadcrumbItem, Button, Card, Input } from 'flowbite-svelte';
	import CodeJamTeam from '../models/team';
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
	let teamInviteCode: string = '';

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

	async function removeMember(teamId: string, memberUserId: string): Promise<string> {
		if (confirm('Are you sure you want to remove this team member?')) {
			try {
				// Call API to remove member
				await removeMemberFromTeam(teamId, memberUserId);
				return teamId;
			} catch (err) {
				toast.error('Failed to remove member');
				console.error(err);
				return '';
			}
		}
		return '';
	}

	async function getAvatarUrl(member: TeamMember): Promise<string> {
		let ext = member.AvatarId.startsWith('a_') ? '.gif' : '.png';
		return `https://cdn.discordapp.com/avatars/${member.ServiceUserId}/${member.AvatarId}${ext}`;
	}

	async function loadAvatarUrls() {
		let members: TeamMember[] = [];

		members.push(...teamMembers);

		console.log(`MEMBERS: ${members}`);

		const promises = members.map(async (member) => {
			const url = await getAvatarUrl(member);

			// for this page, we use Id instead of UserId because it queries the User table, which uses 'Id'
			// may need to make a new class/model if we want to fix the error: Property 'Id' does not exist on type 'TeamMember'.
			avatarUrls[member.Id] = url;
		});

		await Promise.all(promises);
	}

	async function loadData(id: string) {
		try {
			getTeamById(params.id).then((response) => {
				response.json().then((data) => {
					formData = data as CodeJamTeam;
					teamData = data.Team;
					teamInviteCode = data.Team.InviteCode;
					teamMembers = data.TeamMembers;
					teamEvent = data.Event;
					teamName = data.Team.Name;
					teamVisibility = data.Team.Visibility;
					teamAvailability = data.Team.Availability;
					teamTechnologies = data.Team.Technologies;
					teamDescription = data.Team.Description;
					loadAvatarUrls();
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
	}
	let url: string = '';

	$: if (teamData?.Id) {
		url = `localhost:8080/#/team/invite/${teamData.InviteCode}`;
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
					{#each formData.TeamMembers as Member}
						<TableBodyRow>
							<TableBodyCell>
								<Avatar src={avatarUrls[Member.Id]} title={Member.DisplayName} />
							</TableBodyCell>
							<TableBodyCell>{Member.DisplayName}</TableBodyCell>
							<TableBodyCell>{Member.TeamRole}</TableBodyCell>
							<TableBodyCell>
								{#if Member.TeamRole !== 'owner'}
									<Button
										on:click={() =>
											teamData?.Id &&
											Member.Id &&
											removeMember(teamData.Id, Member.Id).then((resTeamId) => {
												if (resTeamId) {
													toast.success("You've successfully removed a member.");
													window.location.reload();
												} else {
													toast.error('Action was canceled or failed.');
												}
											})}
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
		{:else}
			<Spinner />
		{/if}
	</Card>
</Page>
