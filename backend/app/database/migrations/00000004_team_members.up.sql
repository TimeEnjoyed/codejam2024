CREATE TABLE IF NOT EXISTS team_members (
    team_id UUID NOT NULL references teams(id) ON DELETE CASCADE,
    user_id UUID NOT NULL references users(id) ON DELETE CASCADE,
    team_role TEXT NOT NULL,
    created_on TIMESTAMP WITH TIME ZONE DEFAULT (now() AT TIME ZONE('utc')),
    PRIMARY KEY (team_id, user_id)
);

ALTER TABLE teams DROP COLUMN IF EXISTS owner_user_id;

