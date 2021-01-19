db.createUser({
    user: "olimar",
    pwd: "password",
    roles: [{
        role: "readWrite",
        db: "pikmin-database"
    }]
})