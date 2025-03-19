conn = new Mongo();

db = conn.getDB("mediation-platform");

// create indexes
db.user.createIndex({ "email": 1 }, { unique: true });
db.user.createIndex({ "phone_number": 1 }, { unique: true });
