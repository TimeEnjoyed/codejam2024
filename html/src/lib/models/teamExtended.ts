import CodeJamEvent from "./event";
import CodeJamTeam from "./team";
import TeamMember from "./TeamMember";

class CodeJamTeamExtended {
    Event: CodeJamEvent;
    Team: CodeJamTeam;
    TeamMembers: TeamMember[]; 

    constructor() {
        this.Event = new CodeJamEvent;
        this.Team = new CodeJamTeam;
        this.TeamMembers = [];
    }
}

export default CodeJamTeamExtended;