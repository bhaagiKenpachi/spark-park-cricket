-- Create oauth_states table for storing OAuth state parameters
CREATE TABLE IF NOT EXISTS oauth_states (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    state VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    used_at TIMESTAMP WITH TIME ZONE NULL
);

-- Create index on state for fast lookups
CREATE INDEX IF NOT EXISTS idx_oauth_states_state ON oauth_states(state);

-- Create index on expires_at for cleanup
CREATE INDEX IF NOT EXISTS idx_oauth_states_expires_at ON oauth_states(expires_at);

-- Create index on used_at for cleanup
CREATE INDEX IF NOT EXISTS idx_oauth_states_used_at ON oauth_states(used_at);
