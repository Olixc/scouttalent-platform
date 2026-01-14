-- Create moderation_results table
CREATE TABLE IF NOT EXISTS moderation_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    video_id UUID NOT NULL,
    approved BOOLEAN NOT NULL,
    confidence DECIMAL(3,2) NOT NULL,
    flags TEXT[] DEFAULT '{}',
    reason TEXT NOT NULL,
    result_data JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (video_id) REFERENCES videos(id) ON DELETE CASCADE
);

CREATE INDEX idx_moderation_video_id ON moderation_results(video_id);
CREATE INDEX idx_moderation_approved ON moderation_results(approved);
CREATE INDEX idx_moderation_created_at ON moderation_results(created_at DESC);