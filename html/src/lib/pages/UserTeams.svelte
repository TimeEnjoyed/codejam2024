<script lang="ts">
	import Page from '../components/Page.svelte';
	import { Avatar, Button, Breadcrumb, BreadcrumbItem, Card } from 'flowbite-svelte';
	import CodeJamTeam from '../models/team';
	import { getUserTeams } from '../services/services';
	import TeamMember from '../models/TeamMember';
	import { onMount } from 'svelte';
	import { loggedInStore, userStore } from '../stores/stores';

	export const params: Record<string, never> = {};

	// TODO:
	// create edit button if user owns team
	// -- can remove users and edit form inputs
	// -- add invite link to c/p

	let loading: boolean = true;
	let error: string | null = null;
	let userTeams: CodeJamTeam[] = [];
	let avatarUrls: Record<string, string> = {};

	async function loadData() {
		try {
			const response = await getUserTeams();
			userTeams = await response.json(); // Array of teams...
			console.log('userTeams: ', userTeams);
		} catch (err) {
			error = `Failed to load team data: ${err}`;
			console.error(err);
		} finally {
			loading = false;
		}
		loadAvatarUrls();
	}

	async function getAvatarUrl(member: TeamMember): Promise<string> {
		let ext = member.AvatarId.startsWith('a_') ? '.gif' : '.png';
		return `https://cdn.discordapp.com/avatars/${member.ServiceUserId}/${member.AvatarId}${ext}`;
	}

	async function loadAvatarUrls() {
		let members: TeamMember[] = [];

		for (let team of userTeams) {
			if (team['TeamMembers']) {
				members.push(...team.TeamMembers);
			} else {
				console.log('TeamMembers not in team');
			}
		}

		const promises = members.map(async (member) => {
			const url = await getAvatarUrl(member);
			avatarUrls[member.UserId] = url;
		});

		await Promise.all(promises);
	}

	function getTeamOwner(teamMembers: TeamMember[]): string {
		let owner = teamMembers.find((member) => member.TeamRole === 'owner');
		if (owner) {
			return owner.DisplayName;
		} else {
			return 'No owner found.';
		}
	}

	onMount(() => {
		loadData();
	});

	// Currently, I'm querying every team where the loggedin user is the owner, and separating all the members.
	// I'm also checking client side whether the loggined user is the owner or not.
	// -- If it's an owner, it shows an edit button.

	// Following, I think I should show a list of all the teams the loggedin user is only a member in.
	// This means querying for every team where the user.id matches the tm.userId of any team, UNLESS theyre an owner.
</script>

<Page>
	<Breadcrumb solid class="mb-4 w-full max-w-screen-xl">
		<BreadcrumbItem href="/#/" home>Home</BreadcrumbItem>
		<BreadcrumbItem href="/#/teams">My Teams</BreadcrumbItem>
	</Breadcrumb>

	<Card size="md" class="w-full flex">
		<h3>Teams You Own</h3>

		{#if loading}
			<div class="p-4">Loading...</div>
		{:else if error}
			<div class="p-4 text-red-500">{error}</div>
		{:else if userTeams === null}
			<div>You currently don't have any teams!</div>
		{:else if userTeams.length === 0}
			<div>
				Looks like you don't have any teams. Go to <a href="/#/teams/browse">browse</a> teams to
				join one or <a href="/#/teams/create">create</a> your own team!
			</div>
		{:else}
			{#each userTeams as userTeam}
				<Card size="xl" class="flex w-full p-8 px-4 py-6 space-y-3">
					{#if getTeamOwner(userTeam.TeamMembers) == $userStore?.DisplayName}
						<center class="p-2">
							<h4>{userTeam.Name}</h4>
							<a href="/#/team/edit/{userTeam.Id}">Edit your team</a>
						</center>
					{:else}
						<center class="p-2">
							<h4>Team {userTeam.Name}</h4>
						</center>
					{/if}
					<span>
						<b>Owner: </b>{getTeamOwner(userTeam.TeamMembers)}
					</span>
					<span>
						<b>Members: </b>
						<div class="flex mb-5 ml-3">
							{#each userTeam.TeamMembers as Member}
								<Avatar src={avatarUrls[Member.UserId]} title={Member.DisplayName} stacked />
							{/each}
						</div>
					</span>

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
