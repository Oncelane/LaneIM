-- Active: 1725601096974@@127.0.0.1@3306
CREATE DATABASE laneIM DEFAULT CHARACTER SET = 'utf8mb4';

USe laneIM

SELECT `room_mgrs`.`room_id`, `room_mgrs`.`online_count`
FROM
    `room_mgrs`
    JOIN `room_users` ON room_id = room_users.room_mgr_room_id
WHERE
    room_users.user_mgr_user_id = 2

INSERT INTO
    `room_comets` (
        `room_mgr_room_id`,
        `room_comet_room_id`,
        `room_comet_comet_addr`
    )
VALUES (
        1833435292313321472,
        0,
        '127.0.0.1:50050'
    )
ON DUPLICATE KEY UPDATE
    `room_mgr_room_id` = `room_mgr_room_id`

INSERT INTO
    `room_comets` (
        `room_mgr_room_id`,
        `room_comet_room_id`,
        `room_comet_comet_addr`
    )
VALUES (
        1833437093628477440,
        0,
        '127.0.0.1:50050'
    )
ON DUPLICATE KEY UPDATE
    `room_mgr_room_id` = `room_mgr_room_id`