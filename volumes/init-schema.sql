CREATE TABLE "jwt-tokens"(
    id serial primary key,
    "hash" varchar unique,
    "guid" varchar unique,
    "expiry" bigint not null,
    last_login_ip varchar not null
)