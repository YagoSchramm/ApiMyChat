-- Banco inicial da ApiMyChat
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    uid UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS rooms (
    name TEXT NOT NULL,
    id TEXT PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS medias (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type VARCHAR(50) NOT NULL,
    message_id UUID NOT NULL,
    uid UUID NOT NULL,
    url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS room_users (
    room_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    joined_at TIMESTAMP DEFAULT NOW(),
    left_at TIMESTAMP NULL,
    PRIMARY KEY (room_id, user_id),
    FOREIGN KEY (room_id) REFERENCES rooms(id)
);

CREATE INDEX IF NOT EXISTS idx_room_users_user ON room_users(user_id);
CREATE INDEX IF NOT EXISTS idx_room_users_room ON room_users(room_id);

CREATE TABLE IF NOT EXISTS messages (
   id TEXT PRIMARY KEY,
   sender_id TEXT NOT NULL,
   room_id TEXT NOT NULL,
   content TEXT NOT NULL,
   status VARCHAR(20) NOT NULL DEFAULT 'sent',
   created_at TIMESTAMP NOT NULL DEFAULT NOW(),
   CONSTRAINT fk_room
       FOREIGN KEY(room_id)
       REFERENCES rooms(id)
       ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_devices (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 user_id TEXT NOT NULL,
 fcm_token TEXT NOT NULL,
 created_at TIMESTAMP DEFAULT NOW()
);
