export interface User {
    Id: string;
    DisplayName: string;
    Role: string;
    ServiceName : string;
    ServiceUserId: string;
    ServiceUserName: string;
    AvatarId: string;
    AccountStatus: string;
    LockDisplayName: boolean;
}

export interface ActiveUser {
    user : User | null;
    loggedIn : boolean;
}