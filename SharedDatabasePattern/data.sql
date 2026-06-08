-- lozinke su hashovane BCryptom jer Spring Security to ocekuje
-- zapravo su : marija marija123, nenad nenad123, teodora teodora123, srdjan srdjan123

-- Auth service table (minimal, owns auth)
CREATE TABLE IF NOT EXISTS registry (
    id       BIGSERIAL PRIMARY KEY,
    username VARCHAR(32)  UNIQUE NOT NULL,
    password TEXT         NOT NULL,
    email    VARCHAR(255) UNIQUE NOT NULL,
    role     VARCHAR(32)  NOT NULL
);

-- Stakeholder service table (full profile)
CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL PRIMARY KEY,
    registry_id   BIGINT       UNIQUE NOT NULL REFERENCES registry(id) ON DELETE CASCADE,
    username      VARCHAR(32)  UNIQUE NOT NULL,
    email         VARCHAR(255) UNIQUE NOT NULL,
    role          VARCHAR(50)  NOT NULL,
    first_name    VARCHAR(255),
    last_name     VARCHAR(255),
    profile_image VARCHAR(255),
    biography     TEXT,
    motto         VARCHAR(255),
    is_blocked    BOOLEAN      NOT NULL DEFAULT false
);

-- Seed registry
INSERT INTO registry (username, password, email, role)
SELECT 'marija', '$2a$12$HqxKvHNhjWgv88TGuCny5.DAofILnpudYaVj.YSCAVl7pWa69DDBm', 'marija@example.com', 'TOURIST'
WHERE NOT EXISTS (SELECT 1 FROM registry WHERE username = 'marija');

INSERT INTO registry (username, password, email, role)
SELECT 'nenad', '$2a$12$xXUXHPskplCBK4m/vuyzXuTwAS1RvEsAuPUvROnOBnz4VIhyGboRK', 'nenad@example.com', 'ADMIN'
WHERE NOT EXISTS (SELECT 1 FROM registry WHERE username = 'nenad');

INSERT INTO registry (username, password, email, role)
SELECT 'teodora', '$2a$12$2q8BsNnF6OQHfFScUyqCFOTL66S7xKGZIRbDY4PvK0FZIsJzYyPBy', 'teodora@example.com', 'TOURIST'
WHERE NOT EXISTS (SELECT 1 FROM registry WHERE username = 'teodora');

INSERT INTO registry (username, password, email, role)
SELECT 'srdjan', '$2a$12$OFJ9nh7RHS6.EoTR1JIYAesR/StIV.HVVJnGeSCIqqUwJTUVDWlh.', 'srdjan@example.com', 'GUIDE'
WHERE NOT EXISTS (SELECT 1 FROM registry WHERE username = 'srdjan');

-- Seed users (linked to registry)
INSERT INTO users (registry_id, username, email, role, first_name, last_name, profile_image, biography, motto, is_blocked)
SELECT r.id, r.username, r.email, r.role, 'Marija', 'Todorovic', null, null, null, false
FROM registry r WHERE r.username = 'marija'
AND NOT EXISTS (SELECT 1 FROM users WHERE username = 'marija');

INSERT INTO users (registry_id, username, email, role, first_name, last_name, profile_image, biography, motto, is_blocked)
SELECT r.id, r.username, r.email, r.role, 'Nenad', 'Lukic', null, null, null, false
FROM registry r WHERE r.username = 'nenad'
AND NOT EXISTS (SELECT 1 FROM users WHERE username = 'nenad');

INSERT INTO users (registry_id, username, email, role, first_name, last_name, profile_image, biography, motto, is_blocked)
SELECT r.id, r.username, r.email, r.role, 'Teodora', 'Pesic', null, null, null, false
FROM registry r WHERE r.username = 'teodora'
AND NOT EXISTS (SELECT 1 FROM users WHERE username = 'teodora');

INSERT INTO users (registry_id, username, email, role, first_name, last_name, profile_image, biography, motto, is_blocked)
SELECT r.id, r.username, r.email, r.role, 'Srđan', 'Sancanin', null, 'Ja sam jos u stomaku sanjao o tome da budem turisticki vodic i evo me sada!', 'No pain, no gain', false
FROM registry r WHERE r.username = 'srdjan'
AND NOT EXISTS (SELECT 1 FROM users WHERE username = 'srdjan');
