package database

const createSchema = `
CREATE TABLE sites (
	id int NOT NULL AUTO_INCREMENT,
	url varchar(100) NOT NULL,
    frequency float NOT NULL,
    last_execution_date datetime,
    sucess bool,
    response_time float,
    response_average_time float,
    creation_date datetime NOT NULL,
    last_updated_date datetime,
    PRIMARY KEY (id)
)
`

var insertSiteSchema = `
insert into sites (url, frequency, creation_date) values ('https://google.com', 2, '2022-11-23 17:31:26')
`
