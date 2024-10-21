class TeamMember {
    TeamId: string;
    UserId: string;
    TeamRole: string;
    DisplayName: string;
    AvatarId: string;
    ServiceUserId: string;

    //TODO add Array<Teams> 

    constructor() {

        this.TeamId = '';
        this.UserId = '';
        this.TeamRole = '';
        this.DisplayName = '';
        this.AvatarId = '';
        this.ServiceUserId = '';
    }
}

export default TeamMember;