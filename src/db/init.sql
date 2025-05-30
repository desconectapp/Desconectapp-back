DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT now()
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
  activity TEXT NOT NULL,
  description TEXT,

  day_of_week day_option NOT NULL DEFAULT 'EVERYDAY',

  participants_needed INTEGER DEFAULT 3,

  created_at TIMESTAMP DEFAULT NOW(),
  expires_at TIMESTAMP DEFAULT (NOW() + INTERVAL '7 days')
);
