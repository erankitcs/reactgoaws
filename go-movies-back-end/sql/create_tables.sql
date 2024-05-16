CREATE TABLE public.genres (
    id integer NOT NULL,
    genre character varying(255),
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);

ALTER TABLE public.genres ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.genre_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

CREATE TABLE public.movies (
    id integer NOT NULL,
    title character varying(512),
    release_date date,
    runtime integer,
    mpaa_rating character varying(10),
    description text,
    image character varying(255),
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);

CREATE TABLE public.movies_genres (
    id integer NOT NULL,
    movie_id integer,
    genre_id integer
);

ALTER TABLE public.movies_genres ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.movies_genres_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

ALTER TABLE public.movies ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.movies_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

CREATE TABLE public.users (
    id integer NOT NULL,
    first_name character varying(255),
    last_name character varying(255),
    email character varying(255),
    password character varying(255),
    created_at timestamp without time zone,
    updated_at timestamp without time zone 
);

ALTER TABLE public.users ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

CREATE TABLE public.movies_videos (
    id integer NOT NULL,
    movie_id integer NOT NULL,
    video_path character varying(255) NOT NULL,
    is_latest boolean NOT NULL,
    created_at timestamp without time zone
);

ALTER TABLE public.movies_videos ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.movies_videos_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

--Data

INSERT INTO public.genres (genre, created_at, updated_at)
values ('Comedy', '2022-09-23','2022-09-23');
INSERT INTO public.genres (genre, created_at, updated_at)
values ('Sci-Fi', '2022-09-23','2022-09-23');
INSERT INTO public.genres (genre, created_at, updated_at)
values ('Horror', '2022-09-23','2022-09-23');
INSERT INTO public.genres (genre, created_at, updated_at)
values ('Romance', '2022-09-23','2022-09-23');
INSERT INTO public.genres (genre, created_at, updated_at)
values ('Action', '2022-09-23','2022-09-23');
INSERT INTO public.genres (genre, created_at, updated_at)
values ('Thriller', '2022-09-23','2022-09-23');
INSERT INTO public.genres (genre, created_at, updated_at)
values ('Drama', '2022-09-23','2022-09-23');
INSERT INTO public.genres (genre, created_at, updated_at)
values ('Mystery', '2022-09-23','2022-09-23');
INSERT INTO public.genres (genre, created_at, updated_at)
values ('Crime', '2022-09-23','2022-09-23');
INSERT INTO public.genres (genre, created_at, updated_at)
values ('Animation', '2022-09-23','2022-09-23');
INSERT INTO public.genres (genre, created_at, updated_at)
values ('Adventure', '2022-09-23','2022-09-23');
INSERT INTO public.genres (genre, created_at, updated_at)
values ('Fantasy', '2022-09-23','2022-09-23');
INSERT INTO public.genres (genre, created_at, updated_at)
values ('Superhero', '2022-09-23','2022-09-23');

INSERT INTO public.movies ( title, release_date, runtime, mpaa_rating, description, image, created_at, updated_at)
values ('Highlander', '1986-03-07', 116, 'R', 'Russell Nash, a mysterious Scottish swordsman, realises that he is immortal. Soon he finds himself in a ferocious battle against the powerful and wicked immortals who want to destroy the Earth.', '/highlander.jpg', '2022-09-23','2022-09-23');

INSERT INTO public.movies ( title, release_date, runtime, mpaa_rating, description, image, created_at, updated_at)
values ('Raider of Lost Ark', '1981-06-12', 115, 'PG-13', 'Archaeologist and adventurer Indiana Jones is hired by the U.S. government to find the Ark of the Covenant before the Nazis can obtain its powers, in this action-packed fan-favourite.', '/raiders.jpg', '2022-09-23','2022-09-23');

INSERT INTO public.movies ( title, release_date, runtime, mpaa_rating, description, image, created_at, updated_at)
values ('The Godfather', '1972-03-04', 175, '18A', 'Don Vito Corleone, head of a mafia family, decides to hand over his empire to his youngest son, Michael. However, his decision unintentionally puts the lives of his loved ones in grave danger.', '/thegodfather.jpg', '2022-09-23','2022-09-23');

INSERT INTO public.movies_genres (movie_id, genre_id)
values (1,5);
INSERT INTO public.movies_genres (movie_id, genre_id)
values (1,12);
INSERT INTO public.movies_genres (movie_id, genre_id)
values (2,5);
INSERT INTO public.movies_genres (movie_id, genre_id)
values (2,11);
INSERT INTO public.movies_genres (movie_id, genre_id)
values (3,9);
INSERT INTO public.movies_genres (movie_id, genre_id)
values (3,7);
/* password is secret   */
INSERT INTO public.users (first_name, last_name, email, password, created_at, updated_at)
values ('Admin', 'User', 'admin@example.com','$2a$12$jesVW07pibHnym7QbMN/BOHseAM65uDVetWADC.jeaIoi/6mlGN66', '2022-09-23','2022-09-23');