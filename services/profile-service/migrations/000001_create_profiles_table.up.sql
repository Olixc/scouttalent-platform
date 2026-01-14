-- Create user type enum
CREATE TYPE user_type AS ENUM ('player', 'scout', 'academy');

-- Create trust level enum
CREATE TYPE trust_level AS ENUM ('newcomer', 'established', 'verified', 'pro');

-- Create profiles table
CREATE TABLE IF NOT EXISTS profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE,
    type user_type NOT NULL,
    
    -- Common fields
    display_name VARCHAR(100) NOT NULL,
    bio TEXT,
    avatar_url VARCHAR(500),
    location_country VARCHAR(100),
    location_city VARCHAR(100),
    
    -- Trust system
    trust_level trust_level NOT NULL DEFAULT 'newcomer',
    profile_completion_score INTEGER NOT NULL DEFAULT 0,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT valid_completion_score CHECK (profile_completion_score BETWEEN 0 AND 100)
);

-- Create indexes
CREATE INDEX idx_profiles_user_id ON profiles(user_id);
CREATE INDEX idx_profiles_type ON profiles(type);
CREATE INDEX idx_profiles_trust_level ON profiles(trust_level);
CREATE INDEX idx_profiles_location ON profiles(location_country, location_city);

-- Full-text search on display_name and bio
CREATE INDEX idx_profiles_search ON profiles 
    USING gin(to_tsvector('english', display_name || ' ' || COALESCE(bio, '')));

-- Player-specific details
CREATE TABLE IF NOT EXISTS player_details (
    profile_id UUID PRIMARY KEY REFERENCES profiles(id) ON DELETE CASCADE,
    position VARCHAR(50) NOT NULL,
    date_of_birth DATE,
    height_cm INTEGER,
    weight_kg INTEGER,
    preferred_foot VARCHAR(10),
    current_team VARCHAR(100),
    
    -- AI-generated scores (updated by AI workers)
    skill_scores JSONB DEFAULT '{}',
    overall_score DECIMAL(3,1),
    last_scored_at TIMESTAMP,
    
    CONSTRAINT valid_position CHECK (position IN ('goalkeeper', 'defender', 'midfielder', 'forward')),
    CONSTRAINT valid_preferred_foot CHECK (preferred_foot IN ('left', 'right', 'both'))
);

-- Create indexes for player details
CREATE INDEX idx_player_position ON player_details(position);
CREATE INDEX idx_player_overall_score ON player_details(overall_score DESC NULLS LAST);

-- Scout-specific details
CREATE TABLE IF NOT EXISTS scout_details (
    profile_id UUID PRIMARY KEY REFERENCES profiles(id) ON DELETE CASCADE,
    organization VARCHAR(200),
    organization_type VARCHAR(50),
    regions_of_interest TEXT[],
    positions_of_interest TEXT[],
    
    -- Verification
    verified_at TIMESTAMP,
    verified_by UUID,
    verification_documents JSONB DEFAULT '[]'
);

-- Create indexes for scout details
CREATE INDEX idx_scout_regions ON scout_details USING gin(regions_of_interest);
CREATE INDEX idx_scout_positions ON scout_details USING gin(positions_of_interest);