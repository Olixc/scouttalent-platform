CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS videos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    profile_id UUID NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    blob_url TEXT NOT NULL,
    thumbnail_url TEXT,
    duration INTEGER DEFAULT 0,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'uploading',
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT videos_status_check CHECK (status IN ('uploading', 'processing', 'ready', 'failed'))
);

CREATE INDEX idx_videos_profile_id ON videos(profile_id);
CREATE INDEX idx_videos_status ON videos(status);
CREATE INDEX idx_videos_created_at ON videos(created_at DESC);

CREATE TABLE IF NOT EXISTS uploads (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    upload_id VARCHAR(255) NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL DEFAULT 'uploading',
    progress INTEGER DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT uploads_status_check CHECK (status IN ('uploading', 'processing', 'ready', 'failed')),
    CONSTRAINT uploads_progress_check CHECK (progress >= 0 AND progress <= 100)
);

CREATE INDEX idx_uploads_video_id ON uploads(video_id);
CREATE INDEX idx_uploads_upload_id ON uploads(upload_id);