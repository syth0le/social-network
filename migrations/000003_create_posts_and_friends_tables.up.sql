CREATE TABLE IF NOT EXISTS post_table
(
    id         TEXT                     NOT NULL,
    user_id    TEXT                     NOT NULL,
    text       TEXT                     NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT pk_token_table PRIMARY KEY (id),
    CONSTRAINT fk_token_table_user_table FOREIGN KEY (user_id) REFERENCES user_table (id)
);

CREATE TABLE IF NOT EXISTS friend_table
(
    id             TEXT                     NOT NULL,
    first_user_id  TEXT                     NOT NULL,
    second_user_id TEXT                     NOT NULL,

    -- status (expected, accepted, declined, revoked)
    -- expected - 1st user send friend request to 2nd user
    -- accepted - 2nd user accepted friend request of 1st user
    -- declined - 2nd user declined friend request of 1st user
    -- revoked - 1st user revoked friend request to 2nd user (2nd user will be follower of 1st user in this case)

    -- friend for user X -> the person which (second_user_id = X or first_user_id = X) and status = accepted
    -- follower for user X -> (the person which second_user_id = X and status = expected or declined) or
    --            (the person which first_user = X and status = revoked)
    -- subscription for user X -> (the person which first_user_id = X and status = expected or declined) or
    --            (the person which second_user_id = id and status = revoked)

    -- get feed for user X (get feed of friends and subscriptions) ->
    --            (the person which (second_user_id = X or first_user_id = X) and status = accepted) or
    --            (
    --              (the person which first_user_id = X and status = expected or declined) or
    --              (the person which second_user_id = id and status = revoked)
    --            )
--     status         TEXT                     NOT NULL
--         CONSTRAINT friendship_status_field CHECK (status in ('expected', 'accepted', 'declined', 'revoked')),


    created_at     TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at     TIMESTAMP WITH TIME ZONE NOT NULL,

    CONSTRAINT pk_token_table PRIMARY KEY (id),
    CONSTRAINT fk_token_table_user_table_first_user FOREIGN KEY (first_user_id) REFERENCES user_table (id),
    CONSTRAINT fk_token_table_user_table_second_user FOREIGN KEY (second_user_id) REFERENCES user_table (id)
);
