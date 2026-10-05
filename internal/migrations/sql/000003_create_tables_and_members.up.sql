CREATE TABLE tables
(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    dm_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT tables_name_lenght CHECK ( char_length(name) BETWEEN 1 AND 128)
);

CREATE INDEX idx_tables_dm_id ON tables (dm_id);

CREATE TABLE table_members
(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    table_id UUID NOT NULL REFERENCES tables (id) ON DELETE CASCADE,
    player_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT table_members_status_chek CHECK ( status IN ('pending', 'accepted', 'declined')),
    CONSTRAINT table_members_table_player_unique UNIQUE (table_id, player_id)
);

CREATE INDEX idx_table_members_player_id ON table_members (player_id);