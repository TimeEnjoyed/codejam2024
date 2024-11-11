<script lang="ts">
	import {
		Avatar,
		Breadcrumb,
		BreadcrumbItem,
		Button,
		Card,
		Input,
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow,
		TableHead,
		TableHeadCell
	} from 'flowbite-svelte';
	import Helper from 'flowbite-svelte/Helper.svelte';
	import Radio from 'flowbite-svelte/Radio.svelte';
	import Spinner from 'flowbite-svelte/Spinner.svelte';
	import Textarea from 'flowbite-svelte/Textarea.svelte';
	import toast from 'svelte-french-toast';
	import Form from '../components/Form.svelte';
	import FormField from '../components/FormField.svelte';
	import Page from '../components/Page.svelte';
	import type CodeJamEvent from '../models/event';
	import CodeJamTeam from '../models/team';
	import CodeJamTeamExtended from '../models/teamExtended';
	import type TeamMember from '../models/TeamMember';
	import { getTeamById, putTeam, removeMemberFromTeam } from '../services/services';

	export let params: any; // set by svelte-spa-router

	// TODO:
	// Show each team member with a delete button next to username
	// Successfully submit the form (update query to database.)

	let formData: CodeJamTeamExtended;
	let teamMembers: TeamMember[] = [];
	let teamEvent: CodeJamEvent | null = null;
	let avatarUrls: Record<string, string> = {};

	let teamName: string = '';
    let teamId: string = '';
	let teamVisibility: string = '';
	let teamAvailability: string = '';
	let teamTechnologies: string = '';
	let teamDescription: string = '';
	let teamInviteCode: string = '';
	let teamTeamMembers: TeamMember[] = [];

	let loading: boolean = true;
	let isSaving: boolean = false;
	let error: string | null = null;

	let formDataTeamId: string = '';
	let formDataTeamInviteCode: string = '';

	let clearErrors: () => {};
	let parseResponse: (response: object) => {};

	async function saveForm() {
		if (formData !== null) {
			isSaving = true;

			const formDataTeam = formData.Team;
			formDataTeamId = formData.Team.Id;

			const updatedTeam = await putTeam(formDataTeam);
			if (updatedTeam.ok) {
				toast.success("You've successfully edited team info.");
				window.location.href = `/#/team/edit/${formDataTeamId}`;
			} else {
				toast.error('Failed to save team');
				console.error(`response status: ${updatedTeam.status}`);
			}
			isSaving = false;
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

		return Promise.all(promises);
	}

	async function loadData(id: string) {
		try {
			const response = await getTeamById(params.id);

			if (!response.ok) {
				loading = false;
				console.error(`Failed to get team information: ${response.status}`);
				return;
			}

			// You can do this all in where you actually do the fetch/services.ts

			const data: CodeJamTeamExtended = await response.json();

			let teamData = data.Team;
			formData = data as CodeJamTeamExtended;

			teamData = data.Team;
            teamId = data.Team.Id;
			teamName = data.Team.Name;
			teamInviteCode = data.Team.InviteCode;
			teamMembers = data.TeamMembers;
			teamEvent = data.Event;
			teamName = data.Team.Name;

			teamVisibility = data.Team.Visibility;
			teamAvailability = data.Team.Availability;
			teamTechnologies = data.Team.Technologies;
			teamDescription = data.Team.Description;
			teamTeamMembers = data.TeamMembers;
			await loadAvatarUrls();
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

	$: if (teamId) {
        console.log("does thissssssss...")
		url = `localhost:8080/#/team/invite/${teamInviteCode}`;
	}

	$: formData;

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
	<Card size="xl" class="w-full mb-10">
		<div class="flex flex-row gap-8"></div>
		{#if loading}
			<div class="p-4">Loading...</div>
		{:else if error}
			<div class="p-4 text-red-500">{error}</div>
		{:else if formData !== null}
			<div class="flex flex-col gap-8 my-8">
				<Form bind:clearErrors bind:parseResponse>
					<FormField label="Team Name:" name="TeamName">
						<Input bind:value={formData.Team.Name}></Input>
					</FormField>
					<div>
						<Radio name="team-type" bind:group={formData.Team.Visibility} value="public"
							>Public Team</Radio
						>
						<Helper class="ml-6 ">(If you want your team to be searchable.)</Helper>
					</div>
					<div>
						<Radio name="team-type" bind:group={formData.Team.Visibility} value="private"
							>Private Team</Radio
						>
						<Helper class="ml-6">(Your team will be invite only)</Helper>
					</div>

					<FormField label="Your general availability:" name="TeamAvailability">
						<Input bind:value={formData.Team.Availability}></Input>
					</FormField>

					<FormField label="Your technologies:" name="TeamTechnologies">
						<Input bind:value={formData.Team.Technologies}></Input>
					</FormField>

					<FormField label="What do you want out of this team?" name="Description">
						<Textarea bind:value={formData.Team.Description} />
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
					{#each teamTeamMembers as Member}
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
											formData.Team.Id &&
											Member.Id &&
											removeMember(formData.Team.Id, Member.Id).then((resTeamId) => {
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
