CREATE DATABASE portal;
CREATE USER 'developerportal'@'%' IDENTIFIED BY 'changeme';
GRANT ALL PRIVILEGES ON portal.* TO 'developerportal'@'%' WITH GRANT OPTION;
