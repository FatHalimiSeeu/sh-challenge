# How to run:
I offer fully dockerized setup, it will be initial to only run the command below and you are ready to go:
`docker-compose up --build`
Please ensure the ports 1433, 6379, 8080 since i have configured it to also be served from host, due to easier usage of swagger ui.

# How to test:
Please access `http://localhost:8080/swagger/index.html#`, since initally you will have an empty database please register a user, than log in to obtain the token, than use that token to add on Authorization of swagger to use the endpoints that require the Authorization Token, there are also Get (Read) endpoints that do not require token

# Thinking Out Loud – Get a Peek of Why I Made Decisions

## Clean Architecture
Decided to go with as clean an architecture as possible. That's why the project is divided into:
- `config`
- `controllers`
- `middlewares`
- `models`
- `routes`
- `utils`
- `docker essentials`
- `go-essentials`

## 1. Why Using SQL (Relational Database)?
Decided to use a relational database to create relationships between inventory and restock, in order to maintain performance and structure everything better.

## 2. Why Have a Distinction Between Restock and Inventory?
Upon receiving the assignment, there were clear requirements and validations for inventory items and how to restock them. I decided to have separate models/tables for `Restock` and `Inventory`. Each restock has a foreign key, the `ItemId`, which refers to the `ID` from inventory, making them related.

Thinking about how this app might scale in the future, I opted for separate tables so we can add new relationships as needed, allowing the system to scale. When returning the history of restock items or all restocks, it is more efficient to obtain the name through the relation to the inventory table. The `restock` table remains uncluttered, containing only the `ItemId` (foreign key to inventory) and the `amount` of the restock, while the timestamp is automatically created upon insertion.

## 3. Redis & SQL Setup
Redis and SQL are set up in the `config` directory to keep `main.go` clean. Both Redis and SQL are initialized and fully Dockerized for easier setup.

- Connection to SQL is done with retries because the SQL container doesn’t always start quickly, depending on the machine.
- If the database does not exist, it is created inside the Docker volume.
- Auto-migration is enabled, meaning tables and relationships are created automatically.

## 4. Why Redis?
One of the requirements was that write endpoints should contain JWT Authorization. To follow optimal and widely accepted trends used by companies like Uber, Twitter, Netflix, and Airbnb, I chose this approach:

- When a user logs in, a JWT token is generated with a default secret (unchanged for now).
- The token is saved in Redis along with its expiration time. Redis automatically performs garbage collection for expired tokens.
- Redis is fast and memory-intensive, but if the app is expected to scale to millions of users, this trade-off is acceptable.
- The token remains valid until the same user logs in again, invalidating the old token in Redis.
- Tokens can be removed from Redis to force logout a user and revoke access.

## 5. Why Middleware?
Middleware is necessary for apps that require JWT authentication (or any other authentication). It ensures easy maintenance and scalability:

- The `Auth` middleware authenticates tokens by checking Redis.
- It allows requests to proceed if the token is valid and aborts them otherwise (invalid, expired, non-existent, etc.).

## 6. Usage of Utils (Hash & JWT)
`utils` includes reusable functions like hashing and JWT handling. These are shared services that any controller or service can call, making maintenance and scalability easier.

## 7. Fully Dockerized Setup
Setting up projects that depend on SQL Server, Redis, or other dependencies can be painful across different machines due to architecture, environment, usernames, and passwords. To address this, the setup is fully Dockerized:

- Run `docker-compose up --build` to build and start all necessary containers and volumes.
- Access the endpoints through Swagger or an API client like Thunderstorm or Postman.

### Note:
`.env` is empty since all environment variables are set in `docker-compose`.
