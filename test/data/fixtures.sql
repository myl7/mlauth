TRUNCATE users;
INSERT INTO users (uid, username, password, email, display_name, is_active, is_super, created_at)
VALUES (1, 'testusername', '$2a$12$0Ew1ypxmezSM0YV9TQJe8.kygAS8XFGBYvnCXkVv.mi3vjOuCt0/m', 'testemail@gmail.com',
        'test display name', TRUE, FALSE,
        '1999-01-08 04:05:06'),
       (2, 'anotherusername', '$2a$12$rZ.iRNlhtM9UUzk89hUoKedVVw6yy4LgRRIu75R1OYO913KPfBKSu',
        'anotheremail@mail.ustc.edu.cn', 'test display name', TRUE, FALSE,
        '2021-07-08 00:00:00'),
       (3, 'username3', '$2a$12$08blyfUu0siB40qhslSNrujccQi3Xg0qd4NeBYQkcuzYoYB1v/jfu', 'hello@qq.com', 'User 3', TRUE,
        FALSE, '2020-01-01 01:02:03');
SELECT pg_catalog.SETVAL(PG_GET_SERIAL_SEQUENCE('users', 'uid'), MAX(uid))
FROM users;
