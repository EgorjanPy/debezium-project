-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    role TEXT[] NOT NULL DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create index on email for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Insert sample data for testing
INSERT INTO users (email, name, last_name, role) VALUES
    ('john.doe@example.com', 'John', 'Doe', ARRAY['user', 'admin']),
    ('jane.smith@example.com', 'Jane', 'Smith', ARRAY['user']),
    ('bob.wilson@example.com', 'Bob', 'Wilson', ARRAY['user', 'moderator'])
ON CONFLICT (email) DO NOTHING;
