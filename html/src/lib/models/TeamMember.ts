class TeamMember {
    TeamId: string;
    UserId: string;
    TeamRole: string;
    DisplayName: string;
    AvatarId: string;
    ServiceUserID: string;

    //TODO add Array<Teams> 

    constructor() {

        this.TeamId = '';
        this.UserId = '';
        this.TeamRole = '';
        this.DisplayName = '';
        this.AvatarId = '';
        this.ServiceUserID = '';
    }
}

export default TeamMember;