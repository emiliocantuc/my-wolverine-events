-- Script that creates tables
-- Ran with: sqlite3 main.db < tables.sql

-- Creating the users table
CREATE TABLE users (
    user_id INTEGER PRIMARY KEY AUTOINCREMENT,
    email VARCHAR(100) NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Creating the preferences table
CREATE TABLE preferences (
    user_id INT NOT NULL,
    educationLevel VARCHAR(255) DEFAULT '',
    program VARCHAR(255) DEFAULT '' CHECK (educationLevel IN ('Undergraduate', 'Graduate', '')),
    campusLocation VARCHAR (255) DEFAULT 'Central Campus' CHECK (campusLocation IN ('South Campus', 'Central Campus', 'North Campus')),
    interests TEXT DEFAULT '',
    incSeminars INTEGER DEFAULT 1 CHECK (incSeminars IN (0, 1)),
    incSports INTEGER DEFAULT 1 CHECK (incSports IN (0, 1)),
    incSocial INTEGER DEFAULT 1 CHECK (incSocial IN (0, 1)),
    sendEmail INTEGER DEFAULT 1 CHECK (incSocial IN (0, 1)),
    keywordsToAvoid TEXT DEFAULT '',
    PRIMARY KEY (user_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

-- Creating the events table
CREATE TABLE events (
    event_id INTEGER PRIMARY KEY AUTOINCREMENT,
    nweek INTEGER,
    title VARCHAR(100) NOT NULL,
    event_description TEXT,
    event_date DATETIME,
    type VARCHAR(50),
    permalink VARCHAR(255),
    building_name VARCHAR(100),
    building_id INTEGER,
    gcal_link VARCHAR(255),
    umich_id VARCHAR(100),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Creating the recommendations table
CREATE TABLE recommended_events (
    recommendation_id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    event_id INTEGER NOT NULL,
    method VARCHAR(30),
    params VARCHAR(100),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (event_id) REFERENCES events(event_id) ON DELETE CASCADE
);

-- Creating the votes table
CREATE TABLE votes (
    user_id INTEGER NOT NULL,
    event_id INTEGER NOT NULL,
    vote_type CHAR(1) CHECK (vote_type IN ('U', 'D', 'C')),
    voted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, event_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (event_id) REFERENCES events(event_id) ON DELETE CASCADE
);

-- Creating the statistics table
CREATE TABLE statistics (
    nweek INT PRIMARY KEY,
    nusers INT,
    nevents INT
);
INSERT INTO statistics (nweek, nusers, nevents) VALUES (1, 0, 0);

-- Auto create preferences for new user
CREATE TRIGGER after_user_insert
AFTER INSERT ON users
FOR EACH ROW
BEGIN
    INSERT INTO preferences (user_id) VALUES (NEW.user_id);
END;

-- Indexing
CREATE INDEX idx_user_id ON votes (user_id);
CREATE INDEX idx_event_id ON votes (event_id);
CREATE INDEX idx_vote_type ON votes (vote_type);
