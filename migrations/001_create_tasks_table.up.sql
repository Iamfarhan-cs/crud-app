-- Enable UUID generation for task identifiers.
-- The extension is shared database infrastructure, so the down migration does not drop it.
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Titles are required and must contain at least one non-whitespace character.
    title TEXT NOT NULL,

    description TEXT,

    -- Keep task lifecycle values constrained at the database boundary.
    status TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    -- Soft deletes keep historical rows while excluding them from normal reads.
    deleted_at TIMESTAMPTZ,

    CONSTRAINT tasks_title_not_blank CHECK (length(btrim(title)) > 0),
    CONSTRAINT tasks_status_valid CHECK (status IN ('pending', 'in_progress', 'done'))
);

-- Optimize normal list views for active tasks in newest-first order.
CREATE INDEX idx_tasks_active_created_at_desc
    ON tasks (created_at DESC)
    WHERE deleted_at IS NULL;
