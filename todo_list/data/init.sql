CREATE TABLE todos (
    id INT AUTO_INCREMENT PRIMARY KEY,
    description VARCHAR(255) NOT NULL,
    done tinyint(1) NOT NULL DEFAULT 0
);
