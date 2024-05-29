DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS coupons CASCADE;

CREATE TABLE users 
(
    userid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(255) ,
    email VARCHAR(255) ,
    phone_number VARCHAR(255) ,
    password VARCHAR(255) NOT NULL CHECK (octet_length(password) <> 0),
    full_name VARCHAR(255) NOT NULL CHECK (full_name <> ''),
    birth_date DATE DEFAULT NULL,
    last_login TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    discount UUID[] REFERENCES coupons(couponid) DEFAULT NULL,
)

CREATE TABLE coupons {
    couponid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(255) NOT NULL CHECK (code <> ''),
    discount DECIMAL(5,2) NOT NULL CHECK (discount > 0),
    expiration_date DATE DEFAULT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
}