use analytics;

CREATE TABLE visit
(
id INT NOT NULL AUTO_INCREMENT,
siteId CHAR(12) NOT NULL,
visitor CHAR(15) NOT NULL, 
url TEXT NOT NULL,
userAgent VARCHAR(255),
referrer TEXT,
visitTime DATE,
PRIMARY KEY (id)
) COMMENT='visit table';
