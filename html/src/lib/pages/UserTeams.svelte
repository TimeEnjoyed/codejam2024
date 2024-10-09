<script lang="ts">
	import Page from '../components/Page.svelte';
	import { Avatar, Button, Breadcrumb, BreadcrumbItem, Card } from 'flowbite-svelte';

	import CodeJamTeam from '../models/team';
	import CodeJamEvent from '../models/event';
	import { getUserTeams } from '../services/services';
	import TeamMember from '../models/TeamMember';
	import { onMount } from 'svelte';
	import { loggedInStore, userStore } from '../stores/stores';

	export const params: Record<string, never> = {};

	// TODO:
	// show owner of team
	// show members
	// create edit button if user owns team
	// -- can remove users and edit form inputs

	let loading: boolean = true;
	let error: string | null = null;
	let userTeams: CodeJamTeam[] = [];
	let avatarUrls: Record<string, string> = {};

	async function loadData() {
		try {
			const response = await getUserTeams();
			userTeams = await response.json(); // Array of teams...
            console.log(userTeams)
		} catch (err) {
			error = `Failed to load team data: ${err}`;
			console.error(err);
		} finally {
			loading = false;
		}
	}

	async function getAvatarUrl(member: TeamMember): Promise<string> {
		let ext = member.AvatarUrl.startsWith('a_') ? '.gif' : '.png';
		return `https://cdn.discordapp.com/avatars/${member.ServiceUserId}/${member.AvatarUrl}${ext}`;
	}

	async function loadAvatarUrls() {
		let members: TeamMember[] = [];

		for (let team of userTeams) {
			if ('TeamMembers' in team) {
				members.push(...team.TeamMembers);
			} else {
				console.log('TeamMembers not in team');
			}
		}

		console.log(`MEMBERS: ${members}`);

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

	$: userTeams, loadAvatarUrls();

</script>

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
			<div>Error, please contact admin.</div>
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
							<h4>Team {userTeam.Name}</h4>
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
                            {console.log(userTeam.TeamMembers)}
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
