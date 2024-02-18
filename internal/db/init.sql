-- START DEFINITIONS BLOCK 

CREATE TABLE Students (
    id VARCHAR(255),
    has_been_suspended BOOLEAN DEFAULT false,
    PRIMARY KEY (id)
);

CREATE TABLE Teachers (
    id VARCHAR(255),
    PRIMARY KEY (id)
);

CREATE TABLE Registrations (
    teacher_id VARCHAR(255),
    student_id VARCHAR(255),
    PRIMARY KEY (teacher_id, student_id),
    FOREIGN KEY (teacher_id) REFERENCES Teachers(id),
    FOREIGN KEY (student_id) REFERENCES Students(id)
);

CREATE TABLE Notifications (
    sender_id VARCHAR(255),
    recipient_id VARCHAR(255),
    generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    notification VARCHAR NOT NULL,
    PRIMARY KEY (sender_id, recipient_id, generated_at),
    FOREIGN KEY (sender_id) REFERENCES Teachers(id),
    FOREIGN KEY (recipient_id) REFERENCES Students(id)
);

-- END DEFINITIONS BLOCK

-- START SEED DATA BLOCK

INSERT INTO Students (id) VALUES ('studenthon@gmail.com');
INSERT INTO Students (id) VALUES ('studentjon@gmail.com');
INSERT INTO Students (id) VALUES ('studentmiche@gmail.com');
INSERT INTO Students (id) VALUES ('commonstudent1@gmail.com');
INSERT INTO Students (id) VALUES ('commonstudent2@gmail.com');
INSERT INTO Students (id) VALUES ('student_only_under_teacher_ken@gmail.com');
INSERT INTO Students (id) VALUES ('studentmary@gmail.com');
INSERT INTO Students (id) VALUES ('studentbob@gmail.com');
INSERT INTO Students (id) VALUES ('studentagnes@gmail.com');

INSERT INTO Teachers (id) VALUES ('teacherjoe@gmail.com');
INSERT INTO Teachers (id) VALUES ('teacherken@gmail.com');

-- END SEED DATA BLOCK
