import TeamMember from "./TeamMember";

class CodeJamTeam {

    Id : string;
    EventId: string;
    Name: string;
    Visibility: string;
    Technologies: string;
    Availability: string;
    Description: string;
    CreatedOn: Date;
    InviteCode: string;
    TeamMembers: TeamMember[];

    constructor() {
        this.Id = '';
        this.EventId = '';
        this.Name = '';
        this.Visibility = 'public';
        this.Technologies = '';
        this.Availability = '';
        this.Description = '';
        this.CreatedOn = new Date;
        this.InviteCode = '';
        this.TeamMembers = [];
    }
}

export default CodeJamTeam;