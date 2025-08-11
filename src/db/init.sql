DROP TABLE IF EXISTS users CASCADE;

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL
);

DROP TABLE IF EXISTS profiles CASCADE;

CREATE TABLE profiles (
  user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  age INTEGER NOT NULL DEFAULT 18,
  city TEXT NOT NULL DEFAULT 'Buenos Aires',
  current_situation TEXT NOT NULL DEFAULT 'OTHER', 
  gender TEXT NOT NULL DEFAULT 'OTHER',
  profile_complete BOOLEAN NOT NULL DEFAULT false
);

DROP TYPE IF EXISTS categories;

CREATE TYPE categories as ENUM (
  'SPORT', 'CREATIVE', 'OUTDOOR', 'INDOOR', 'GAME', 'SOCIAL', 'WELLNESS'
);

DROP TABLE IF EXISTS activities CASCADE;

CREATE TABLE activities(
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  icon TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  category categories NOT NULL DEFAULT 'SOCIAL'
);

DROP TABLE IF EXISTS users_preference;

CREATE TABLE users_preference (
  user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
  activity_id INTEGER REFERENCES activities(id) ON DELETE CASCADE,
  PRIMARY KEY (user_id, activity_id)
);


DROP TABLE IF EXISTS activity_requests;

CREATE TABLE activity_requests (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
  activity_id INTEGER REFERENCES activities(id) ON DELETE CASCADE,
  description TEXT,
  week_hours INTEGER[],
  participants_needed INTEGER DEFAULT 3,
  maximum_participants INTEGER DEFAULT 10,
  latitude DOUBLE PRECISION,
  longitude DOUBLE PRECISION,
  search_radius INTEGER DEFAULT 10,
  created_at TIMESTAMP DEFAULT NOW(),
  expires_at TIMESTAMP DEFAULT (NOW() + INTERVAL '7 days')
);

DROP TABLE IF EXISTS partial_matches;

CREATE TABLE partial_matches (
    id SERIAL PRIMARY KEY,
    activity_id INTEGER REFERENCES activities(id) ON DELETE CASCADE,
    description TEXT,
    week_hours INTEGER[],
    participants_needed INTEGER DEFAULT 3,
    maximum_participants INTEGER DEFAULT 10,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    search_radius INTEGER DEFAULT 10,
    members_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP DEFAULT (NOW() + INTERVAL '7 days')
);

DROP TABLE IF EXISTS partial_match_members;

CREATE TABLE partial_match_members (
    partial_match_id INTEGER REFERENCES partial_matches(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (partial_match_id, user_id)
);

DROP TABLE IF EXISTS groups;

CREATE TABLE groups (
    id SERIAL PRIMARY KEY,
    name TEXT,
    description TEXT,
    location TEXT,
    activity_id INTEGER NOT NULL REFERENCES activities(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

DROP TABLE IF EXISTS group_members;

CREATE TABLE group_members (
    group_id INTEGER NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (group_id, user_id)
);

COPY activities(name, icon, category)
FROM '/activities.csv'
WITH (FORMAT csv, HEADER true);

COPY users(email, password)
FROM '/users.csv'
WITH (FORMAT csv, HEADER true);

COPY profiles(user_id, name, age, city, current_situation, profile_complete)
FROM '/profiles.csv'
WITH (FORMAT csv, HEADER true);

COPY users_preference(user_id, activity_id)
FROM '/users_preference.csv'
WITH (FORMAT csv, HEADER true);
