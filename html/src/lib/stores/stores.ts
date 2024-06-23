
import {writable, derived} from 'svelte/store';
import CodeJamEvent from "../models/event";

export const activeContent = writable('');
export const userStore = writable(null);
export const activeEventStore = writable<CodeJamEvent | null>(null);
export const eventStatusStore = writable(null);
export const loggedInStore = derived(userStore, (userData) => userData != null);