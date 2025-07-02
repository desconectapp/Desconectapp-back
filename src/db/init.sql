DROP TYPE IF EXISTS current_situation_options;

CREATE TYPE current_situation_options AS ENUM (
  'WORKING FULLTIME', 'WORKING PARTIME', 'STUDENT', 'STAY AT HOME', 'ENTREPRENEUR', 'LOTS OF FREE TIME', 'LOOKING FOR WORK', 'OTHER'
);

DROP TABLE IF EXISTS users CASCADE;

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  age INTEGER NOT NULL DEFAULT 18,
  city TEXT NOT NULL DEFAULT 'Buenos Aires, Argentina',
  current_situation current_situation_options NOT NULL DEFAULT 'OTHER', 
  profile_complete BOOLEAN NOT NULL DEFAULT false
);

DROP TYPE IF EXISTS categories;

CREATE TYPE categories as ENUM (
  'SPORT', 'CREATIVE', 'OUTDOOR', 'INDOOR', 'GAME', 'SOCIAL', 'WELLNESS'
);

DROP TABLE IF EXISTS activities;

CREATE TABLE activities(
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  emoji TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  category categories NOT NULL DEFAULT 'SOCIAL'
);

DROP TABLE IF EXISTS users_preference;

CREATE TABLE users_preference (
  user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
  activity_id INTEGER REFERENCES activities(id) ON DELETE CASCADE,
  PRIMARY KEY (user_id, activity_id)
);

DROP TYPE IF EXISTS day_option;

CREATE TYPE day_option AS ENUM (
  'EVERYDAY', 'WEEKDAYS', 'WEEKEND',
  'MONDAY', 'TUESDAY', 'WEDNESDAY', 'THURSDAY', 'FRIDAY', 'SATURDAY', 'SUNDAY'
);


DROP TABLE IF EXISTS activity_requests;

CREATE TABLE activity_requests (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
  activity_id INTEGER REFERENCES activities(id) ON DELETE CASCADE,
  description TEXT,

  day_of_week day_option NOT NULL DEFAULT 'EVERYDAY',

  participants_needed INTEGER DEFAULT 3,

  created_at TIMESTAMP DEFAULT NOW(),
  expires_at TIMESTAMP DEFAULT (NOW() + INTERVAL '7 days')
);

CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    refresh_token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    refresh_expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sessions_token ON sessions(token);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);

COPY activities(name, emoji, category)
FROM '/activities.csv'
WITH (FORMAT csv, HEADER true);


CREATE OR REPLACE FUNCTION cleanup_expired_sessions()
RETURNS void AS $$
BEGIN
    DELETE FROM sessions WHERE expires_at < CURRENT_TIMESTAMP;
END;
$$ LANGUAGE plpgsql; 