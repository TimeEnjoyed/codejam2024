import {baseApiUrl} from "./services";
import type {User} from "../models/user";


export async function getAllUsers() {
    return await fetch(baseApiUrl + "/admin/user/all");
}


interface PutAccountStatusRequest {
    AccountStatus: string,
}

export async function putAccountStatus(userId: string, status: string) : Promise<User> {
    const requestInit : RequestInit = {
        method: 'PUT',
        body: JSON.stringify(
            <PutAccountStatusRequest>{
                AccountStatus: status
            }
        )
    }
    const response: Response = await fetch(`${baseApiUrl}/admin/user/${userId}/account_status/`, requestInit);
    return await response.json() as User;
}

export async function BanUser(userId: string) :Promise<User> {
    const response : Response = await fetch(`${baseApiUrl}/admin/user/${userId}/ban`, { method: 'PUT' });
    return await response.json() as User;
}

export async function UnbanUser(userId: string) : Promise<User> {
    const response : Response = await fetch(`${baseApiUrl}/admin/user/${userId}/unban`, { method: 'PUT' });
    return await response.json() as User;
}


interface PutDisplayNameLockRequest {
    Lock: boolean,
}

export async function updateDisplayNameLock(userId: string, lock: boolean) {
    const requestInit : RequestInit = {
        method: 'PUT',
        body: JSON.stringify(
            <PutDisplayNameLockRequest>{
                Lock: lock
            }
        )
    }
    const response: Response = await fetch(`${baseApiUrl}/admin/user/${userId}/display_name_lock`, requestInit);
    return await response.json() as User;
}


interface PutDisplayNameRequest {
    DisplayName: string,
}

export async function UpdateDisplayName(userId: string, displayName: string) {
    const requestInit : RequestInit = {
        method: 'PUT',
        body: JSON.stringify(
            <PutDisplayNameRequest>{
                DisplayName: displayName
            }
        )
    }
    const response : Response  = await fetch(`${baseApiUrl}/admin/user/${userId}/display_name`, requestInit);
    return await response.json() as User;
}