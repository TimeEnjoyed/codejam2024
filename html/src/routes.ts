import HomePage from "./lib/pages/HomePage.svelte";
import EventEdit from "./lib/pages/EventEdit.svelte";
import EventList from "./lib/pages/EventList.svelte";
import TeamOptions from "./lib/pages/TeamOptions.svelte";
import TeamsBrowse from "./lib/pages/TeamsBrowse.svelte";
import TeamsCreate from "./lib/pages/TeamsCreate.svelte";
import MyTeam from "./lib/pages/MyTeam.svelte";
import Invite from "./lib/pages/Invite.svelte";
import UserTeams from "./lib/pages/UserTeams.svelte";
import ProfilePage from "./lib/pages/ProfilePage.svelte";
import TeamEdit from "./lib/pages/TeamEdit.svelte";

export default {
    '/': HomePage,
    '/home': HomePage,
    '/admin/events': EventList,
    '/admin/event/:id': EventEdit,
    '/team': TeamOptions,
    '/team/:id': MyTeam, // link to one of your teams (sharable) We get an id here in this route...
    '/team/invite/:invitecode': Invite, // sharable
    '/teams': UserTeams, // displays all the user's teams (private)
    '/teams/browse': TeamsBrowse,
    '/teams/create': TeamsCreate,
    '/profile': ProfilePage,
    '/team/edit/:id': TeamEdit
}