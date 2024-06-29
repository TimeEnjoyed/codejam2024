-- capture the original user so admins can see it without having to lookup a discord id
ALTER TABLE users ADD COLUMN IF NOT EXISTS service_user_name TEXT;

-- allow for a user to be disabled or banned, or have some other sort of status that's TBD
ALTER TABLE users ADD COLUMN IF NOT EXISTS account_status TEXT DEFAULT 'ACTIVE';

-- allow for a user to be prevented from modifying their display name
ALTER TABLE users ADD COLUMN IF NOT EXISTS lock_display_name BOOLEAN DEFAULT FALSE;

-- assumes no real users are in the system yet
UPDATE users SET service_user_name = display_name;

