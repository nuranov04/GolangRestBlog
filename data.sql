CREATE TABLE public.user
(
    id            SERIAL       NOT NULL PRIMARY KEY unique ,
    username      VARCHAR(100) NOT NULL unique,
    email         VARCHAR(100) NOT NULL unique ,
    password_hash VARCHAR(500) NOT NULL
);

CREATE TABLE public.post
(
    id          SERIAL       NOT NULL PRIMARY KEY,
    title       VARCHAR(100) NOT NULL,
    description TEXT         NOT NULL,
    owner_id    INTEGER      NOT NULL,
    CONSTRAINT owner FOREIGN KEY (owner_id) REFERENCES public.user(id)

);

INSERT INTO public.user (username, email, password_hash) VALUES ('admin', 'admin@gmail.com', 'admin');
INSERT INTO public.user (username, email, password_hash) VALUES ('admin1', 'admin1@gmail.com', 'admin');


INSERT INTO public.post (title, description, owner_id) VALUES ('title', 'description', '1');
INSERT INTO public.post (title, description, owner_id) VALUES ('title', 'description', '2');