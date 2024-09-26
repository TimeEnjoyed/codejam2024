import {
    activeEventStore,
    eventStatusStore,
    userStore,
    activeUserStore
} from "../stores/stores";
import CodeJamEvent from "../models/event";
import CodeJamTeam from "../models/team";
import type {ActiveUser, User} from "../models/user";

// This shouldn't ever need to be set since dev and prod environments will just use relative endpoints
export let baseApiUrl: string = "";

const originalFetch = window.fetch;

// Override standard "fetch" so we have a place to generically handle errors
window.fetch = (...args) => {
    return new Promise((resolve) => {
        originalFetch(...args)
            .then((response) => {
                resolve(response);
            })
            .catch((err) => {
                console.error("fetch error:", err, ...args);
                throw err;
            })
    });
}

export async function getUser() {
    return fetch(baseApiUrl + "/user/", { method: 'GET', credentials: 'include' })
        .then((response) => {
            if (response.status === 401) {
                userStore.set(null);
                activeUserStore.set(<ActiveUser>{user: null, loggedIn: false});
            } else {
                response.json()
                    .then((data) => {
                        userStore.set(data);
                        activeUserStore.set(<ActiveUser>{user: data as User, loggedIn: true})
                    })
                    .catch((err) => {
                        console.error("error deserializing user", err);
                    });
            }
        });
}

export async function putProfile(user: User): Promise<Response> {
    return fetch(baseApiUrl + "/user/",
        {
            method: "PUT",
            body: JSON.stringify(user)
        });
}

export async function logout() {
    return fetch(baseApiUrl + "/user/logout")
        .then(() => {
            userStore.set(null);
            activeUserStore.set(<ActiveUser>{user: null, loggedIn: false});
        })
        .catch((err) => {
            console.error("Logout error", err);
        });
}

export async function getActiveEvent() {
    return fetch(baseApiUrl + "/event/active")
        .then((response) => {
            if (response.status === 401) {
                userStore.set(null);
            } else if (response.status === 204) {
                activeEventStore.set(null);
            } else {
                response.json()
                    .then((data) => {
                        activeEventStore.set(data as CodeJamEvent);
                    })
                    .catch((err) => {
                        console.error("error deserializing event", response, err);
                    });
            }
        });
}

export async function getEvents() {
    return fetch(baseApiUrl + "/event/");
}

export async function getEvent(id: string) {
    return fetch(baseApiUrl + "/event/" + id);
}

export async function putEvent(event: CodeJamEvent) {
    return await fetch(baseApiUrl + "/event/" + event.Id,
        {
            method: "PUT",
            body: JSON.stringify(event)
        });
}

export async function getEventStatuses() {
    return fetch(baseApiUrl + "/event/statuses")
        .then((response) => {
            response.json()
                .then((data) => {
                    eventStatusStore.set(data)
                });
        })
}

export async function postTeam(team: CodeJamTeam) {
    // Step 2: Post Team Data API (accepts formData(CodeJamTeam) as argument)
    return await fetch(baseApiUrl + "/team/" + team.Id,
        // formData turned into JSON and sent via POST, retreived at CreateTeam(ctx *gin.Context)
        {
            method: "POST",
            body: JSON.stringify(team)
        });
}

export async function getTeams() {
    return await fetch(baseApiUrl + "/teams/browse")
}


export async function getUserTeams(){
    return fetch(baseApiUrl + "/teams");
}

export async function getTeamById(id: string) {
    return fetch(baseApiUrl + "/team/" + id);
}

export async function getTeamByInvite(inviteCode: string) {
    // stepp 3 pt 1:
    return fetch(baseApiUrl + "/team/invite/" + inviteCode);
}


// any one who has a link joinTeam() works on, can join the team.
// make sure invite_code matches 

export async function joinTeam(teamId: string, inviteCode: string) {
    // making a post to team_members
    return await fetch(baseApiUrl + "/team/" + inviteCode,
        {
            method: "POST",
            body: JSON.stringify({ teamId, inviteCode })
        }
    )
}


// write a joinPublicTeam function that accepts team Id.  
// server: check if team is public.
export async function joinPublicTeam(teamId: string) {
    const response = await fetch(baseApiUrl + "/team/join",
        {
            method: "POST",
            body: JSON.stringify({ teamId })
        }
    )
    return response.json()
}

// Always call at startup to get the initial states
async function initialLoad() {
    await Promise.all([getUser(), getActiveEvent(), getEventStatuses()]);
}

initialLoad();