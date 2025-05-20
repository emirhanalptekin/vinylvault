CREATE TABLE IF NOT EXISTS artists (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS genres (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    icon TEXT
);

CREATE TABLE IF NOT EXISTS albums (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    artist_id TEXT NOT NULL,
    release_year TEXT NOT NULL,
    genre_id TEXT NOT NULL,
    notes TEXT,
    rating INTEGER CHECK (rating BETWEEN 1 AND 5),
    condition TEXT CHECK (condition IN ('Mint', 'Excellent', 'Very Good', 'Good', 'Fair', 'Poor')),
    FOREIGN KEY (artist_id) REFERENCES artists(id),
    FOREIGN KEY (genre_id) REFERENCES genres(id)
);

-- Seed data for artists
INSERT INTO artists (id, name) VALUES
    ('art-001', 'Pink Floyd'),
    ('art-002', 'David Bowie'),
    ('art-003', 'Miles Davis'),
    ('art-004', 'Nirvana'),
    ('art-005', 'Radiohead');

-- Seed data for genres
INSERT INTO genres (id, name, icon) VALUES
    ('gen-001', 'Rock', '🎸'),
    ('gen-002', 'Jazz', '🎷'),
    ('gen-003', 'Electronic', '🎹'),
    ('gen-004', 'Hip Hop', '🎤'),
    ('gen-005', 'Classical', '🎻');

-- Seed data for albums
INSERT INTO albums (id, title, artist_id, release_year, genre_id, notes, rating, condition) VALUES
    ('alb-001', 'The Dark Side of the Moon', 'art-001', '1973', 'gen-001', 'Original pressing with posters and stickers', 5, 'Excellent'),
    ('alb-002', 'Kind of Blue', 'art-003', '1959', 'gen-002', 'Columbia Records pressing from the 1970s', 5, 'Very Good'),
    ('alb-003', 'OK Computer', 'art-005', '1997', 'gen-001', 'Special edition with artwork booklet', 4, 'Mint'),
    ('alb-004', 'Nevermind', 'art-004', '1991', 'gen-001', 'First pressing with insert', 4, 'Good'),
    ('alb-005', 'Ziggy Stardust', 'art-002', '1972', 'gen-001', 'Original RCA pressing', 5, 'Very Good');