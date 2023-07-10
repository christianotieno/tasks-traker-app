CREATE TABLE tasks (
                       id VARCHAR(36) PRIMARY KEY,
                       summary VARCHAR(255) NOT NULL,
                       date DATE NOT NULL,
                       user_id VARCHAR(30) NOT NULL,
                       FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE users (
                       id VARCHAR(36) PRIMARY KEY,
                       first_name VARCHAR(50) NOT NULL,
                       last_name VARCHAR(50) NOT NULL,
                       email VARCHAR(100) NOT NULL,
                       password VARBINARY(255) NOT NULL
);

CREATE TABLE managers (
                          id VARCHAR(36) PRIMARY KEY,
                          manager_id VARCHAR(36) NOT NULL,
                          technician_id VARCHAR(16) NOT NULL,
                          FOREIGN KEY (manager_id) REFERENCES users(id),
                          FOREIGN KEY (technician_id) REFERENCES users(id)
);
