import {activeEventStore, eventStatusStore, userStore} from "../stores/stores";
import CodeJamEvent from "../models/event";

// This shouldn't ever need to be set since dev and prod environments will just use relative endpoints
export let baseApiUrl : string = "";

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
    return fetch(baseApiUrl + "/user/", {method: 'GET', credentials: 'include'})
        .then((response) => {
            if (response.status === 401) {
                userStore.set(null);
            } else {
                response.json()
                    .then((data) => {
                        userStore.set(data);
                    })
                    .catch((err) => {
                        console.error("error deserializing user", err);
                    });
            }
        });
}

export async function logout() {
    return fetch(baseApiUrl + "/user/logout")
        .then(() => {
            userStore.set(null);
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

// Always call at startup to get the initial states
async function initialLoad() {
    getUser();
    getActiveEvent();
    getEventStatuses();
}

initialLoad();