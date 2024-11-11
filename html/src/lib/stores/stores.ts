
import {writable, derived} from 'svelte/store';
import CodeJamEvent from "../models/event";
import {type ActiveUser, type User} from "../models/user";

export const activeContent = writable('');

/**
 * Deprecated - use activeUserStore.user instead
 */

export const userStore = writable<User | null>(null);

/**
 * Deprecated - use activeUserStore.loggedIn instead
 */

export const loggedInStore = derived(userStore, (userData) => userData != null);

/**
 * activeUserStore will allow components to use a single object / store to handle loggedIn logic as well as User
 * data handling.  This solves a problem with being able to correctly initialize component states based on logged-in and
 * User values without causing flickering due to DOM updates.
 */
export const activeUserStore = writable<ActiveUser | null>(null);

export const activeEventStore = writable<CodeJamEvent | null>(null);
export const eventStatusStore = writable(null);
