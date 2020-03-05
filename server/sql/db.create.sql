CREATE TABLE stocks( 
    id SERIAL PRIMARY KEY,
    ts_code VARCHAR(100) NOT NULL,
    symbol VARCHAR(100) NOT NULL,
    stock_name VARCHAR(100) NOT NULL,
    area VARCHAR(100) NOT NULL,
    industry VARCHAR(100) NOT NULL,
    market VARCHAR(100) NOT NULL,
    list_date VARCHAR(100) NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now()
);
