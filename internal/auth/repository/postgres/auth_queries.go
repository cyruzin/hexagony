package postgres

const sqlGetUser = "SELECT * from users WHERE email = $1"
