CREATE OR REPLACE DATABASE `einyx`;
CREATE OR REPLACE USER 'grom' IDENTIFIED BY 'changemeplease';
GRANT ALL PRIVILEGES ON einyx.* TO 'grom'@'%';
FLUSH PRIVILEGES;