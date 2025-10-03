DROP TABLE IF EXISTS students;
DROP TABLE IF EXISTS tests;
DROP TABLE IF EXISTS questions;

CREATE TABLE students (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    age INT NOT NULL
);

CREATE TABLE tests (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE questions (
    id VARCHAR(32) PRIMARY KEY,
    question VARCHAR(255) NOT NULL,
    answer VARCHAR(255) NOT NULL,
    test_id VARCHAR(32) NOT NULL,
    FOREIGN KEY (test_id) REFERENCES tests(id)
);

INSERT INTO students (id, name, age) VALUES ('student1', 'John Doe', 20);
INSERT INTO students (id, name, age) VALUES ('student2', 'Jane Doe', 21);
INSERT INTO students (id, name, age) VALUES ('student3', 'John Smith', 22);
INSERT INTO students (id, name, age) VALUES ('student4', 'Jane Smith', 23);

INSERT INTO tests (id, name) VALUES ('test1', 'Test 1');
INSERT INTO tests (id, name) VALUES ('test2', 'Test 2');
INSERT INTO tests (id, name) VALUES ('test3', 'Test 3');
INSERT INTO tests (id, name) VALUES ('test4', 'Test 4');
INSERT INTO tests (id, name) VALUES ('test5', 'Test 5');

INSERT INTO questions (id, question, answer, test_id) VALUES ('question1', 'Question 1', 'Answer 1', 'test1');
INSERT INTO questions (id, question, answer, test_id) VALUES ('question2', 'Question 2', 'Answer 2', 'test2');
INSERT INTO questions (id, question, answer, test_id) VALUES ('question3', 'Question 3', 'Answer 3', 'test3');
INSERT INTO questions (id, question, answer, test_id) VALUES ('question4', 'Question 4', 'Answer 4', 'test4');
INSERT INTO questions (id, question, answer, test_id) VALUES ('question5', 'Question 5', 'Answer 5', 'test5');