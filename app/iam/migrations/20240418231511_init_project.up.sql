CREATE TABLE users
(
    id         INT AUTO_INCREMENT PRIMARY KEY,
    username   varchar(30) NOT NULL UNIQUE,
    name       VARCHAR(30) NOT NULL,
    school     VARCHAR(100),
    class      VARCHAR(30),
    role       VARCHAR(30) NOT NULL,
    created_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP            DEFAULT NULL
);

CREATE TABLE user_passwords
(
    id         INT AUTO_INCREMENT PRIMARY KEY,
    user_id    INT       NOT NULL,
    password   VARCHAR(128),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP          DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id)
);




