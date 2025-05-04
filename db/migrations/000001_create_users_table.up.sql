/* Comment for those who want to perform migration without schema change */
CREATE TABLE Users (
    UserID CHAR(36) NOT NULL,
    Name VARCHAR(100),
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (UserID)
);
