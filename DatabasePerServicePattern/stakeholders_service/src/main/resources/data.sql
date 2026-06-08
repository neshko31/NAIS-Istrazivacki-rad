-- Stakeholder service table
CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL PRIMARY KEY,
    registry_id   BIGINT       UNIQUE NOT NULL,
    username      VARCHAR(32)  UNIQUE NOT NULL,
    email         VARCHAR(255) UNIQUE NOT NULL,
    role          VARCHAR(32)  NOT NULL,
    first_name    VARCHAR(255),
    last_name     VARCHAR(255),
    profile_image VARCHAR(255),
    biography     TEXT,
    motto         VARCHAR(255),
    is_blocked    BOOLEAN      NOT NULL DEFAULT false
);

-- Seed users; registry_id mirrors the auth-service user IDs
INSERT INTO users (registry_id, username, email, role, first_name, last_name, profile_image, biography, motto, is_blocked)
SELECT 1, 'marija', 'marija@example.com', 'TOURIST', 'Marija', 'Todorovic', null, null, null, false
WHERE NOT EXISTS (SELECT 1 FROM users WHERE username = 'marija');

INSERT INTO users (registry_id, username, email, role, first_name, last_name, profile_image, biography, motto, is_blocked)
SELECT 2, 'nenad', 'nenad@example.com', 'ADMIN', 'Nenad', 'Lukic', null, null, null, false
WHERE NOT EXISTS (SELECT 1 FROM users WHERE username = 'nenad');

INSERT INTO users (registry_id, username, email, role, first_name, last_name, profile_image, biography, motto, is_blocked)
SELECT 3, 'teodora', 'teodora@example.com', 'TOURIST', 'Teodora', 'Pesic', null, null, null, false
WHERE NOT EXISTS (SELECT 1 FROM users WHERE username = 'teodora');

INSERT INTO users (registry_id, username, email, role, first_name, last_name, profile_image, biography, motto, is_blocked)
SELECT 4, 'srdjan', 'srdjan@example.com', 'GUIDE', 'Srđan', 'Sancanin', null, 'Ja sam jos u stomaku sanjao o tome da budem turisticki vodic i evo me sada!', 'No pain, no gain', false
WHERE NOT EXISTS (SELECT 1 FROM users WHERE username = 'srdjan');
