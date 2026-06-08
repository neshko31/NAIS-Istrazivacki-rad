CREATE TABLE IF NOT EXISTS users (
    id       SERIAL PRIMARY KEY,
    username VARCHAR(32) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    email VARCHAR(32) UNIQUE NOT NULL,
    role VARCHAR(32) NOT NULL
);

INSERT INTO users (username, password, email, role)
SELECT 'marija', '$2a$12$HqxKvHNhjWgv88TGuCny5.DAofILnpudYaVj.YSCAVl7pWa69DDBm', 'marija@example.com', 'tourist'
WHERE NOT EXISTS (SELECT 1 FROM users WHERE username = 'marija');

INSERT INTO users (username, password, email, role)
SELECT 'nenad', '$2a$12$xXUXHPskplCBK4m/vuyzXuTwAS1RvEsAuPUvROnOBnz4VIhyGboRK', 'nenad@example.com', 'admin'
WHERE NOT EXISTS (SELECT 1 FROM users WHERE username = 'nenad');

INSERT INTO users (username, password, email, role)
SELECT 'teodora', '$2a$12$2q8BsNnF6OQHfFScUyqCFOTL66S7xKGZIRbDY4PvK0FZIsJzYyPBy', 'teodora@example.com', 'tourist'
WHERE NOT EXISTS (SELECT 1 FROM users WHERE username = 'teodora');

INSERT INTO users (username, password, email, role)
SELECT 'srdjan', '$2a$12$OFJ9nh7RHS6.EoTR1JIYAesR/StIV.HVVJnGeSCIqqUwJTUVDWlh.', 'srdjan@example.com', 'guide'
WHERE NOT EXISTS (SELECT 1 FROM users WHERE username = 'srdjan');
