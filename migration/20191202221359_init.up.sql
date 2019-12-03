create type TodoType as enum (
    'Homework',
    'Coding',
    'Report',
    'Discussion',
    ''
    );

create table Todo
(
    id           bigserial not null,
    content      text      not null,
    due          date,
    estimateCost interval,
    type         TodoType
);

create unique index Todo_id_uindex
    on Todo (id);

alter table Todo
    add constraint Todo_pk
        primary key (id);

create table StudentTodo
(
    student_id char(8) not null,
    todo_id    bigint  not null,
    constraint table_name_pk
        primary key (todo_id)
);

alter table studenttodo
    add constraint studenttodo_todo_id_fk
        foreign key (todo_id) references todo;

alter table studenttodo
    drop constraint studenttodo_todo_id_fk;

alter table studenttodo
    add constraint studenttodo_todo_id_fk
        foreign key (todo_id) references todo;