class TeamMember {
    TeamId: string;
    Id: string;
    TeamRole: string;
    DisplayName: string;
    AvatarId: string;
    ServiceUserId: string;

    //TODO add Array<Teams> 

    constructor() {

        this.TeamId = '';
        this.Id = '';
        this.TeamRole = '';
        this.DisplayName = '';
        this.AvatarId = '';
        this.ServiceUserId = '';
    }
}

export default TeamMember;