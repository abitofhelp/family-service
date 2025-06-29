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
    '00000000-0000-0000-0000-000000000100',
    'MARRIED',
    '[
        {
            "ID": "00000000-0000-0000-0000-000000000101",
            "FirstName": "John",
            "LastName": "Smith",
            "BirthDate": "1980-05-15T00:00:00Z",
            "DeathDate": null
        },
        {
            "ID": "00000000-0000-0000-0000-000000000102",
            "FirstName": "Jane",
            "LastName": "Smith",
            "BirthDate": "1982-08-22T00:00:00Z",
            "DeathDate": null
        }
    ]',
    '[
        {
            "ID": "00000000-0000-0000-0000-000000000103",
            "FirstName": "Emily",
            "LastName": "Smith",
            "BirthDate": "2010-03-12T00:00:00Z",
            "DeathDate": null
        },
        {
            "ID": "00000000-0000-0000-0000-000000000104",
            "FirstName": "Michael",
            "LastName": "Smith",
            "BirthDate": "2012-11-05T00:00:00Z",
            "DeathDate": null
        }
    ]'
),
(
    '00000000-0000-0000-0000-000000000200',
    'MARRIED',
    '[
        {
            "ID": "00000000-0000-0000-0000-000000000201",
            "FirstName": "Robert",
            "LastName": "Johnson",
            "BirthDate": "1975-12-10T00:00:00Z",
            "DeathDate": null
        },
        {
            "ID": "00000000-0000-0000-0000-000000000202",
            "FirstName": "Maria",
            "LastName": "Johnson",
            "BirthDate": "1978-04-28T00:00:00Z",
            "DeathDate": null
        }
    ]',
    '[
        {
            "ID": "00000000-0000-0000-0000-000000000203",
            "FirstName": "David",
            "LastName": "Johnson",
            "BirthDate": "2008-07-19T00:00:00Z",
            "DeathDate": null
        }
    ]'
),
(
    '00000000-0000-0000-0000-000000000300',
    'DIVORCED',
    '[
        {
            "ID": "00000000-0000-0000-0000-000000000301",
            "FirstName": "Thomas",
            "LastName": "Williams",
            "BirthDate": "1970-09-30T00:00:00Z",
            "DeathDate": null
        }
    ]',
    '[]'
),
(
    '00000000-0000-0000-0000-000000000400',
    'WIDOWED',
    '[
        {
            "ID": "00000000-0000-0000-0000-000000000401",
            "FirstName": "Sarah",
            "LastName": "Brown",
            "BirthDate": "1985-02-14T00:00:00Z",
            "DeathDate": null
        },
        {
            "ID": "00000000-0000-0000-0000-000000000402",
            "FirstName": "James",
            "LastName": "Brown",
            "BirthDate": "1983-11-08T00:00:00Z",
            "DeathDate": "2020-04-15T00:00:00Z"
        }
    ]',
    '[
        {
            "ID": "00000000-0000-0000-0000-000000000403",
            "FirstName": "Olivia",
            "LastName": "Brown",
            "BirthDate": "2015-06-23T00:00:00Z",
            "DeathDate": null
        },
        {
            "ID": "00000000-0000-0000-0000-000000000404",
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

--  Copyright (c) 2025 A Bit of Help, Inc.

-- Print confirmation (this will only work in psql)
-- \echo 'PostgreSQL initialization completed successfully!'
-- \echo 'Created table: families'
-- \echo 'Inserted sample data'
-- SELECT COUNT(*) AS family_count FROM families;
