-- Lozinke su zapravo: 
-- marija: marija123, 
-- nenad: nenad123, 
-- teodora: teodora123, 
-- srdjan: srdjan123

CREATE TABLE IF NOT EXISTS users (
    id       BIGSERIAL PRIMARY KEY,
    username VARCHAR(32) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    email    VARCHAR(255) UNIQUE NOT NULL,
    role     VARCHAR(32) NOT NULL
);

INSERT INTO users (id, username, password, email, role)
SELECT 1, 'marija', '$2a$12$HqxKvHNhjWgv88TGuCny5.DAofILnpudYaVj.YSCAVl7pWa69DDBm', 'marija@example.com', 'TOURIST'
WHERE NOT EXISTS (SELECT 1 FROM users WHERE username = 'marija');

INSERT INTO users (id, username, password, email, role)
SELECT 2, 'nenad', '$2a$12$xXUXHPskplCBK4m/vuyzXuTwAS1RvEsAuPUvROnOBnz4VIhyGboRK', 'nenad@example.com', 'ADMIN'
WHERE NOT EXISTS (SELECT 1 FROM users WHERE username = 'nenad');

INSERT INTO users (id, username, password, email, role)
SELECT 3, 'teodora', '$2a$12$2q8BsNnF6OQHfFScUyqCFOTL66S7xKGZIRbDY4PvK0FZIsJzYyPBy', 'teodora@example.com', 'TOURIST'
WHERE NOT EXISTS (SELECT 1 FROM users WHERE username = 'teodora');

INSERT INTO users (id, username, password, email, role)
SELECT 4, 'srdjan', '$2a$12$OFJ9nh7RHS6.EoTR1JIYAesR/StIV.HVVJnGeSCIqqUwJTUVDWlh.', 'srdjan@example.com', 'GUIDE'
WHERE NOT EXISTS (SELECT 1 FROM users WHERE username = 'srdjan');
