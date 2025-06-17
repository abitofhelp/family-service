// MongoDB Initialization Script for Family Service

// Connect to the database (or create it if it doesn't exist)
db = db.getSiblingDB('family_service');

// Drop collections if they exist to ensure a clean start
db.families.drop();

// Create indexes for better query performance
db.families.createIndex({ "_id": 1 });
db.families.createIndex({ "parents.id": 1 });
db.families.createIndex({ "children.id": 1 });

// Insert sample families
db.families.insertMany([
  {
    _id: "f1a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
    status: "active",
    parents: [
      {
        id: "p1a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
        firstName: "John",
        lastName: "Smith",
        birthDate: new Date("1980-05-15").toISOString(),
        deathDate: null
      },
      {
        id: "p2a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
        firstName: "Jane",
        lastName: "Smith",
        birthDate: new Date("1982-08-22").toISOString(),
        deathDate: null
      }
    ],
    children: [
      {
        id: "c1a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
        firstName: "Emily",
        lastName: "Smith",
        birthDate: new Date("2010-03-12").toISOString(),
        deathDate: null
      },
      {
        id: "c2a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
        firstName: "Michael",
        lastName: "Smith",
        birthDate: new Date("2012-11-05").toISOString(),
        deathDate: null
      }
    ]
  },
  {
    _id: "f2a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
    status: "active",
    parents: [
      {
        id: "p3a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
        firstName: "Robert",
        lastName: "Johnson",
        birthDate: new Date("1975-12-10").toISOString(),
        deathDate: null
      },
      {
        id: "p4a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
        firstName: "Maria",
        lastName: "Johnson",
        birthDate: new Date("1978-04-28").toISOString(),
        deathDate: null
      }
    ],
    children: [
      {
        id: "c3a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
        firstName: "David",
        lastName: "Johnson",
        birthDate: new Date("2008-07-19").toISOString(),
        deathDate: null
      }
    ]
  },
  {
    _id: "f3a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
    status: "divorced",
    parents: [
      {
        id: "p5a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
        firstName: "Thomas",
        lastName: "Williams",
        birthDate: new Date("1970-09-30").toISOString(),
        deathDate: null
      }
    ],
    children: []
  },
  {
    _id: "f4a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
    status: "active",
    parents: [
      {
        id: "p6a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
        firstName: "Sarah",
        lastName: "Brown",
        birthDate: new Date("1985-02-14").toISOString(),
        deathDate: null
      },
      {
        id: "p7a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
        firstName: "James",
        lastName: "Brown",
        birthDate: new Date("1983-11-08").toISOString(),
        deathDate: new Date("2020-04-15").toISOString()
      }
    ],
    children: [
      {
        id: "c4a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
        firstName: "Olivia",
        lastName: "Brown",
        birthDate: new Date("2015-06-23").toISOString(),
        deathDate: null
      },
      {
        id: "c5a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
        firstName: "William",
        lastName: "Brown",
        birthDate: new Date("2017-09-11").toISOString(),
        deathDate: null
      }
    ]
  }
]);

// Print confirmation
print("MongoDB initialization completed successfully!");
print("Created database: family_service");
print("Created collection: families");
print("Inserted sample data: " + db.families.count() + " families");