-- Drop tables in reverse order of creation

DROP TABLE IF EXISTS t_oauth_accounts;
DROP TABLE IF EXISTS t_refresh_tokens;
DROP TABLE IF EXISTS m_users;
DROP TABLE IF EXISTS m_roles;

-- Note: We generally don't drop extensions as they might be used by other schemas
