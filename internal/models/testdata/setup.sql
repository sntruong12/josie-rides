CREATE TABLE rides (
    id                  INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title               VARCHAR(255) NOT NULL,
    description         LONGTEXT NOT NULL,
    trail_name          VARCHAR(255) NULL,
    distance_miles      DECIMAL(5,2) NOT NULL,
    duration            INT NOT NULL,
    rode_at             DATETIME NOT NULL,
    created_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    media               JSON NULL,

    INDEX idx_rode_at (rode_at),
    INDEX idx_trail_name (trail_name),

    CONSTRAINT chk_distance_miles_positive CHECK (distance_miles > 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created_at DATETIME NOT NULL
);

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);

INSERT INTO users (name, email, hashed_password, created_at) VALUES (
    'Alice Jones',
    'alice@example.com',
    '$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG',
    '2022-01-01 10:00:00'
);