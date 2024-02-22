CREATE TABLE IF NOT EXISTS todolist (
id bigserial PRIMARY KEY,
title text NOT NULL,
description text NOT NULL,
completed boolean NOT NULL,
);
