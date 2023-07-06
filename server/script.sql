CREATE TABLE users (
                       id INT PRIMARY KEY AUTO_INCREMENT,
                       first_name VARCHAR(50) NOT NULL,
                       last_name VARCHAR(50) NOT NULL,
                       email VARCHAR(100) NOT NULL,
                       password VARCHAR(100) NOT NULL,
                       role VARCHAR(50) NOT NULL
);

CREATE TABLE tasks (
                       id INT PRIMARY KEY AUTO_INCREMENT,
                       user_id INT NOT NULL,
                       summary VARCHAR(255) NOT NULL,
                       date DATE NOT NULL,
                       FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE managers (
                          id INT PRIMARY KEY AUTO_INCREMENT,
                          manager_id INT NOT NULL,
                          technician_id INT NOT NULL,
                          FOREIGN KEY (manager_id) REFERENCES users(id),
                          FOREIGN KEY (technician_id) REFERENCES users(id)
);
