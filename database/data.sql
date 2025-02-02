-- Insert Genres
INSERT
   OR IGNORE INTO genres (name)
VALUES ('Action'),
   ('Sci-Fi'),
   ('Drama');
-- Insert Movies
INSERT
   OR IGNORE INTO movies (title, description, release_date, image_url)
VALUES (
      'Inception',
      'A mind-bending thriller',
      '2010-07-16',
      'inception.jpg'
   ),
   (
      'The Matrix',
      'A sci-fi classic',
      '1999-03-31',
      'matrix.jpg'
   );
-- Insert Movie_Genre
INSERT
   OR IGNORE INTO movie_genre (movie_id, genre_id)
VALUES (1, 2),
   (2, 2),
   (2, 1);
