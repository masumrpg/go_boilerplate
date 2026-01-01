-- Drop tables in reverse order of creation

DROP TABLE IF EXISTS t_oauth_accounts CASCADE;
DROP TABLE IF EXISTS t_sessions CASCADE;
DROP TABLE IF EXISTS m_users CASCADE;
DROP TABLE IF EXISTS m_roles CASCADE;

-- Note: We generally don't drop extensions as they might be used by other schemas
