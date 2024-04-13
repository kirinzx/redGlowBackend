CREATE OR REPLACE FUNCTION create_user(
    username VARCHAR(255),
    password VARCHAR(255),
    email VARCHAR(255)
) RETURNS BIGINT AS
$$
DECLARE
    inserted_id BIGINT;
    err_column_name TEXT;
BEGIN
    BEGIN
        INSERT INTO users("username", "password", "email") VALUES (username,crypt(password, gen_salt('bf')),email) RETURNING id INTO inserted_id;
    EXCEPTION
        WHEN unique_violation THEN
            RAISE EXCEPTION 'That % is already taken', CASE 
                WHEN SQLERRM LIKE '%users_username_key%' THEN 'username'
                WHEN SQLERRM LIKE '%users_email_key%' THEN 'email'
            END;
    END;
    RETURN inserted_id;
END;
$$
LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION commit_user(
    user_email VARCHAR(255)
) RETURNS TABLE (
    username VARCHAR(255),
    photo_path VARCHAR(255),
    background_path VARCHAR(255),
    email VARCHAR(255),
    steam_id BOOLEAN
) AS
$$
BEGIN
    RETURN QUERY
    WITH updated_rows AS (
        UPDATE users
        SET is_active = true
        WHERE email = user_email AND is_active = false
        RETURNING *
    )
    SELECT 
        COALESCE(updated_rows.username, users.username) AS username,
        COALESCE(updated_rows.photo_path, users.photo_path) AS photo_path,
        COALESCE(updated_rows.background_path, users.background_path) AS background_path,
        COALESCE(updated_rows.steam_id IS NOT NULL, false) AS steam_id
    FROM users
    LEFT JOIN updated_rows ON users.email = updated_rows.email
    WHERE users.email = user_email;
END;
$$
LANGUAGE PLPGSQL;



CREATE OR REPLACE FUNCTION check_user_by_credits(
    user_email VARCHAR(255),
    user_password VARCHAR(255)
) RETURNS TABLE (
    username VARCHAR(255),
    photo_path VARCHAR(255),
    background_path VARCHAR(255),
    email VARCHAR(255),
    steam_id BOOLEAN
) AS
$$
BEGIN
    RETURN QUERY
    SELECT 
        users.username,
        users.photo_path,
        users.background_path,
        users.email,
        CASE 
            WHEN users.steam_id IS NOT NULL THEN true
            ELSE false
        END AS steam_id
    FROM users
    WHERE is_active=true AND users.email=user_email AND users.password=crypt(user_password, users.password);
    IF NOT FOUND THEN
        RAISE EXCEPTION 'User not found';
    END IF;
END;
$$
LANGUAGE PLPGSQL;

CREATE OR REPLACE PROCEDURE change_password(
    user_email VARCHAR(255),
    old_password VARCHAR(255),
    new_password VARCHAR(255)
) LANGUAGE PLPGSQL AS
$$
BEGIN
    UPDATE users
    SET password = crypt(new_password, gen_salt('bf'))
    WHERE email = user_email AND password = crypt(old_password, password);

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Incorrect old password';
    END IF;
END;
$$;

CREATE OR REPLACE PROCEDURE change_forgotten_password(
    user_email VARCHAR(255),
    new_password VARCHAR(255)
) LANGUAGE PLPGSQL AS
$$
BEGIN
    UPDATE users
    SET password = crypt(new_password, gen_salt('bf'))
    WHERE email = user_email;
    IF NOT FOUND THEN
        RAISE EXCEPTION 'Incorrect email';
    END IF;
END;
$$;

CREATE OR REPLACE FUNCTION check_email(
    user_email VARCHAR(255)
) RETURNS BOOLEAN AS
$$
BEGIN
    RETURN QUERY SELECT EXISTS(SELECT * FROM users WHERE email=user_email);
END;
$$
LANGUAGE plpgsql;