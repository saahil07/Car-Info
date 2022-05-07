USE CarDealership;

Drop table if exists car, engine;

create table car (
                     id varchar(36) NOT NULL,
                     name varchar(50) NOT null,
                     year int(4) not null,
                     brand varchar(50) NOT NULL,
                     fuel varchar(50)NOT null,
                     engineId varchar(36) NOT NULL,
                     PRIMARY KEY (id)
);

create table engine(
                       engineId varchar(36) NOT NULL,
                       displacement int,
                       noOfCylinders int,
                       engineRange int,
                       PRIMARY KEY (engineId)
);