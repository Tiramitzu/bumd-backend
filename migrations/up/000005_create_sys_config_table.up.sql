CREATE TABLE sys_config (
    id                          SERIAL    PRIMARY KEY,
    otp_expired_minutes         INT NOT NULL DEFAULT 5,
    jwt_expired_minutes         INT NOT NULL DEFAULT 480,
    refresh_token_expired_hour  INT NOT NULL DEFAULT 336
);

INSERT INTO sys_config(otp_expired_minutes, jwt_expired_minutes, refresh_token_expired_hour) VALUES (5,3600,2000);