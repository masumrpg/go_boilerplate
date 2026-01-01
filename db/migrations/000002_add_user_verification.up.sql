-- Add is_verified column to m_users table
ALTER TABLE m_users ADD COLUMN IF NOT EXISTS is_verified BOOLEAN DEFAULT FALSE;

-- Update existing users to be verified (so we don't lock out current users)
UPDATE m_users SET is_verified = TRUE;

-- Create index for faster lookups
CREATE INDEX IF NOT EXISTS idx_m_users_is_verified ON m_users(is_verified);
