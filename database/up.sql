DROP TABLE IF EXISTS students;

CREATE TABLE students (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    age INT NOT NULL
);

INSERT INTO students (name, age) VALUES ('John Doe', 20);
INSERT INTO students (name, age) VALUES ('Jane Doe', 21);
INSERT INTO students (name, age) VALUES ('John Smith', 22);
INSERT INTO students (name, age) VALUES ('Jane Smith', 23);