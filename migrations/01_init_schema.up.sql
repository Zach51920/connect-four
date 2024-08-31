-- Disable all triggers
SET session_replication_role = 'replica';

-- Start a transaction
BEGIN;


CREATE TABLE IF NOT EXiSTS games
(
    id         UUID PRIMARY KEY NOT NULL,
    created_at TIMESTAMP        NOT NULL DEFAULT NOW(),
    winner     UUID                      DEFAULT NULL,
    is_draw    BOOL             NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS moves
(
    id           UUID PRIMARY KEY NOT NULL,
    game_id      UUID             NOT NULL,
    board_column INT              NOT NULL,
    player       UUID             NOT NULL,
    turn         INT              NOT NULL,
    CONSTRAINT fk_move_game_id FOREIGN KEY (game_id) REFERENCES games (id)
);

-- Commit transaction
COMMIT;

-- Enable all triggers
SET session_replication_role = 'origin';