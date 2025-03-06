CREATE TABLE commissions (
    id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
     budget FLOAT NOT NULL,
      currency VARCHAR(255) NOT NULL,
       location VARCHAR(255) NOT NULL,
        deadline DATE NOT NULL,
         status VARCHAR(255) NOT NULL,
          winner VARCHAR(255) NOT NULL,
           entries ENUM ('') NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

);
