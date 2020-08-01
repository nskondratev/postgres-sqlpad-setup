create table users_evaluation_of_satisfaction (
    request_id varchar(255) primary key,
    result_mentioned_by_user varchar(255)
);


create table support_tickets (
    user_id int,
    ticket_category varchar(255),
    ticket_subcategory varchar(255),
    current_state varchar(255),
    request_id varchar(255),
    activity_start_dt timestamp,
    fact_reaction_dt timestamp,

    primary key(user_id, request_id)
);

create table new_items_by_support_users (
    user_id int,
    user_registration_time timestamp,
    user_first_listing_date timestamp,
    item_id int primary key,
    item_starttime timestamp,
    item_category varchar(255),
    item_subcategory varchar(255)
);
