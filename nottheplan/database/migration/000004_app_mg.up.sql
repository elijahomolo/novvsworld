CREATE TABLE submissions (
                             id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
                             commission_id int NOT NULL,
                                user_id int NOT NULL,
                                demo_link VARCHAR(255) NOT NULL,
                             submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,


);