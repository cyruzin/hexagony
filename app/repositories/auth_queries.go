package postgres

const sqlAuthGetUser = "SELECT * from users WHERE email = $1"
