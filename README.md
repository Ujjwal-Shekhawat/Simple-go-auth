# Simple-go-auth

#### A simple authentication system made with golang

### `How to run`

- Open your terminal and paste this line `git clone https://* github.com/Ujjwal-Shekhawat/Simple-go-auth`
- Then cd into Simple-go-auth dir and create .env file there using this command `touch .env` and then refer to **env vars** section
- Then type this in your terminal `go run server.go` . This will automatically install all the dependencies and start the server on [localhost:8080](localhost:8080)

### Dependencies

- github.com/dgrijalva/jwt-go v3.2.0
- github.com/gin-gonic/gin v1.6.3
- github.com/joho/godotenv v1.3.0
- go.mongodb.org/mongo-driver v1.4.2

### Env Vars

- `JWT_SECERET` your json web token seceret for eg (JWT_SECERET=somerandomstuff)
- `MONGO_URI` your mongo uri for

## License

MIT

[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com) [![forthebadge](https://forthebadge.com/images/badges/built-with-love.svg)](https://forthebadge.com)
