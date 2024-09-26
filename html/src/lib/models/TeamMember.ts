class TeamMember {
    TeamId: string;
    UserId: string;
    TeamRole: string;
    DisplayName: string;
    AvatarUrl: string;
    ServiceUserId: string;

    //TODO add Array<Teams> 

    constructor() {

        this.TeamId = '';
        this.UserId = '';
        this.TeamRole = '';
        this.DisplayName = '';
        this.AvatarUrl = '';
        this.ServiceUserId = '';
    }
}

export default TeamMember;