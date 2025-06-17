-- PostgreSQL Initialization Script for Family Service

-- Create database if it doesn't exist
-- Note: This needs to be run as a superuser
-- CREATE DATABASE family_service;

-- Connect to the database
-- \c family_service;

-- Drop tables if they exist to ensure a clean start
DROP TABLE IF EXISTS families;

-- Create tables
CREATE TABLE families (
    id VARCHAR(36) PRIMARY KEY,
    status VARCHAR(20) NOT NULL,
    parents JSONB NOT NULL,
    children JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX idx_families_status ON families(status);
CREATE INDEX idx_families_parents ON families USING GIN (parents);
CREATE INDEX idx_families_children USING GIN (children);

-- Insert sample families
INSERT INTO families (id, status, parents, children) VALUES
(
    'f1a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6',
    'active',
    '[
        {
            "ID": "p1a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
            "FirstName": "John",
            "LastName": "Smith",
            "BirthDate": "1980-05-15T00:00:00Z",
            "DeathDate": null
        },
        {
            "ID": "p2a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
            "FirstName": "Jane",
            "LastName": "Smith",
            "BirthDate": "1982-08-22T00:00:00Z",
            "DeathDate": null
        }
    ]',
    '[
        {
            "ID": "c1a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
            "FirstName": "Emily",
            "LastName": "Smith",
            "BirthDate": "2010-03-12T00:00:00Z",
            "DeathDate": null
        },
        {
            "ID": "c2a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
            "FirstName": "Michael",
            "LastName": "Smith",
            "BirthDate": "2012-11-05T00:00:00Z",
            "DeathDate": null
        }
    ]'
),
(
    'f2a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6',
    'active',
    '[
        {
            "ID": "p3a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
            "FirstName": "Robert",
            "LastName": "Johnson",
            "BirthDate": "1975-12-10T00:00:00Z",
            "DeathDate": null
        },
        {
            "ID": "p4a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
            "FirstName": "Maria",
            "LastName": "Johnson",
            "BirthDate": "1978-04-28T00:00:00Z",
            "DeathDate": null
        }
    ]',
    '[
        {
            "ID": "c3a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
            "FirstName": "David",
            "LastName": "Johnson",
            "BirthDate": "2008-07-19T00:00:00Z",
            "DeathDate": null
        }
    ]'
),
(
    'f3a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6',
    'divorced',
    '[
        {
            "ID": "p5a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
            "FirstName": "Thomas",
            "LastName": "Williams",
            "BirthDate": "1970-09-30T00:00:00Z",
            "DeathDate": null
        }
    ]',
    '[]'
),
(
    'f4a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6',
    'active',
    '[
        {
            "ID": "p6a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
            "FirstName": "Sarah",
            "LastName": "Brown",
            "BirthDate": "1985-02-14T00:00:00Z",
            "DeathDate": null
        },
        {
            "ID": "p7a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
            "FirstName": "James",
            "LastName": "Brown",
            "BirthDate": "1983-11-08T00:00:00Z",
            "DeathDate": "2020-04-15T00:00:00Z"
        }
    ]',
    '[
        {
            "ID": "c4a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
            "FirstName": "Olivia",
            "LastName": "Brown",
            "BirthDate": "2015-06-23T00:00:00Z",
            "DeathDate": null
        },
        {
            "ID": "c5a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
            "FirstName": "William",
            "LastName": "Brown",
            "BirthDate": "2017-09-11T00:00:00Z",
            "DeathDate": null
        }
    ]'
);

-- Create a function to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create a trigger to automatically update the updated_at column
CREATE TRIGGER update_families_updated_at
BEFORE UPDATE ON families
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Print confirmation (this will only work in psql)
-- \echo 'PostgreSQL initialization completed successfully!'
-- \echo 'Created table: families'
-- \echo 'Inserted sample data'
-- SELECT COUNT(*) AS family_count FROM families;