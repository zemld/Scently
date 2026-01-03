-- +goose Up
-- +goose StatementBegin
ALTER TABLE upper_notes
ADD CONSTRAINT upper_notes_note_fkey
FOREIGN KEY (note) REFERENCES notes(name);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE core_notes
ADD CONSTRAINT core_notes_note_fkey
FOREIGN KEY (note) REFERENCES notes(name);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE base_notes
ADD CONSTRAINT base_notes_note_fkey
FOREIGN KEY (note) REFERENCES notes(name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE upper_notes DROP CONSTRAINT IF EXISTS upper_notes_note_fkey;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE core_notes DROP CONSTRAINT IF EXISTS core_notes_note_fkey;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE base_notes DROP CONSTRAINT IF EXISTS base_notes_note_fkey;
-- +goose StatementEnd