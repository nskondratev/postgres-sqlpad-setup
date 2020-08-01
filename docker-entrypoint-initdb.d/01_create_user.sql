-- Create a group
CREATE ROLE ro;

-- Grant access to existing tables
GRANT USAGE ON SCHEMA public TO ro;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO ro;

-- Grant access to future tables
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO ro;

-- Grant all privileges only to sqlpad user
REVOKE ALL ON SCHEMA public FROM public;
GRANT ALL ON SCHEMA public TO sqlpad;

-- Create a final user with password
CREATE USER sqlpad_user WITH PASSWORD 'sqlpad_user';
GRANT ro TO sqlpad_user;
