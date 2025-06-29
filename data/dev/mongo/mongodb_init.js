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
    _id: "00000000-0000-0000-0000-000000000100",
    status: "MARRIED",
    parents: [
      {
        id: "00000000-0000-0000-0000-000000000101",
        firstName: "John",
        lastName: "Smith",
        birthDate: new Date("1980-05-15").toISOString(),
        deathDate: null
      },
      {
        id: "00000000-0000-0000-0000-000000000102",
        firstName: "Jane",
        lastName: "Smith",
        birthDate: new Date("1982-08-22").toISOString(),
        deathDate: null
      }
    ],
    children: [
      {
        id: "00000000-0000-0000-0000-000000000103",
        firstName: "Emily",
        lastName: "Smith",
        birthDate: new Date("2010-03-12").toISOString(),
        deathDate: null
      },
      {
        id: "00000000-0000-0000-0000-000000000104",
        firstName: "Michael",
        lastName: "Smith",
        birthDate: new Date("2012-11-05").toISOString(),
        deathDate: null
      }
    ]
  },
  {
    _id: "00000000-0000-0000-0000-000000000200",
    status: "MARRIED",
    parents: [
      {
        id: "00000000-0000-0000-0000-000000000201",
        firstName: "Robert",
        lastName: "Johnson",
        birthDate: new Date("1975-12-10").toISOString(),
        deathDate: null
      },
      {
        id: "00000000-0000-0000-0000-000000000202",
        firstName: "Maria",
        lastName: "Johnson",
        birthDate: new Date("1978-04-28").toISOString(),
        deathDate: null
      }
    ],
    children: [
      {
        id: "00000000-0000-0000-0000-000000000203",
        firstName: "David",
        lastName: "Johnson",
        birthDate: new Date("2008-07-19").toISOString(),
        deathDate: null
      }
    ]
  },
  {
    _id: "00000000-0000-0000-0000-000000000300",
    status: "DIVORCED",
    parents: [
      {
        id: "00000000-0000-0000-0000-000000000301",
        firstName: "Thomas",
        lastName: "Williams",
        birthDate: new Date("1970-09-30").toISOString(),
        deathDate: null
      }
    ],
    children: []
  },
  {
    _id: "00000000-0000-0000-0000-000000000400",
    status: "WIDOWED",
    parents: [
      {
        id: "00000000-0000-0000-0000-000000000401",
        firstName: "Sarah",
        lastName: "Brown",
        birthDate: new Date("1985-02-14").toISOString(),
        deathDate: null
      },
      {
        id: "00000000-0000-0000-0000-000000000402",
        firstName: "James",
        lastName: "Brown",
        birthDate: new Date("1983-11-08").toISOString(),
        deathDate: new Date("2020-04-15").toISOString()
      }
    ],
    children: [
      {
        id: "00000000-0000-0000-0000-000000000403",
        firstName: "Olivia",
        lastName: "Brown",
        birthDate: new Date("2015-06-23").toISOString(),
        deathDate: null
      },
      {
        id: "00000000-0000-0000-0000-000000000404",
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
