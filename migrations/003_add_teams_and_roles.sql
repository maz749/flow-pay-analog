-- Migration for Teams and Roles system

-- Create roles enum type
CREATE TYPE user_role AS ENUM ('admin', 'user');

-- Add role column to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS role user_role DEFAULT 'user';

-- Update existing users to admin (first user should be admin)
UPDATE users SET role = 'admin' WHERE id = (SELECT MIN(id) FROM users);

-- Create teams table
CREATE TABLE IF NOT EXISTS teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    owner_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    budget_amount DECIMAL(10, 2),
    budget_currency VARCHAR(10) DEFAULT 'RUB',
    budget_period VARCHAR(20), -- monthly, yearly
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create team_members table (many-to-many relationship)
CREATE TABLE IF NOT EXISTS team_members (
    id SERIAL PRIMARY KEY,
    team_id INTEGER REFERENCES teams(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(team_id, user_id)
);

-- Add team_id to subscriptions table
ALTER TABLE subscriptions ADD COLUMN IF NOT EXISTS team_id INTEGER REFERENCES teams(id) ON DELETE SET NULL;

-- Create invitations table
CREATE TABLE IF NOT EXISTS invitations (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    role user_role DEFAULT 'user',
    team_id INTEGER REFERENCES teams(id) ON DELETE CASCADE,
    invited_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    token VARCHAR(255) UNIQUE NOT NULL,
    status VARCHAR(20) DEFAULT 'pending', -- pending, accepted, expired
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_teams_owner ON teams(owner_id);
CREATE INDEX idx_team_members_team ON team_members(team_id);
CREATE INDEX idx_team_members_user ON team_members(user_id);
CREATE INDEX idx_subscriptions_team ON subscriptions(team_id);
CREATE INDEX idx_invitations_email ON invitations(email);
CREATE INDEX idx_invitations_token ON invitations(token);
CREATE INDEX idx_invitations_status ON invitations(status);

-- Create trigger for teams updated_at
CREATE TRIGGER update_teams_updated_at BEFORE UPDATE ON teams
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert default "Personal" team for existing users
INSERT INTO teams (name, description, owner_id)
SELECT
    'Личные подписки',
    'Личные подписки пользователя',
    id
FROM users
WHERE NOT EXISTS (
    SELECT 1 FROM teams WHERE owner_id = users.id
);

-- Assign existing subscriptions to user's personal team
UPDATE subscriptions s
SET team_id = t.id
FROM teams t
WHERE s.user_id = t.owner_id
AND s.team_id IS NULL;

-- Add all users to their personal teams
INSERT INTO team_members (team_id, user_id)
SELECT t.id, t.owner_id
FROM teams t
WHERE NOT EXISTS (
    SELECT 1 FROM team_members WHERE team_id = t.id AND user_id = t.owner_id
);
