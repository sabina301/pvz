CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       email VARCHAR(50) UNIQUE NOT NULL,
                       password_hash VARCHAR(200) NOT NULL,
                       role VARCHAR(20) NOT NULL
);

CREATE TABLE pvzs (
                      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                      registrationDate TIMESTAMP DEFAULT now(),
                      city VARCHAR(20) NOT NULL CHECK ( city IN ('Москва', 'Санкт-Петербург', 'Казань') )
);

CREATE TABLE receptions (
                            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                            createdAt TIMESTAMP NOT NULL DEFAULT now(),
                            pvzId UUID REFERENCES pvzs(id) ON DELETE CASCADE,
                            status VARCHAR(20) NOT NULL
);

CREATE TABLE products (
                          id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                          receivedAt TIMESTAMP NOT NULL DEFAULT now(),
                          receptionId UUID NOT NULL REFERENCES receptions(id) ON DELETE CASCADE,
                          type VARCHAR(30) NOT NULL CHECK (type IN ('электроника', 'одежда', 'обувь'))
);