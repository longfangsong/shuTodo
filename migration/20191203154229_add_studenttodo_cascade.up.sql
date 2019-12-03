alter table studenttodo
    drop constraint studenttodo_todo_id_fk;

alter table studenttodo
    add constraint studenttodo_todo_id_fk
        foreign key (todo_id) references todo
            on delete cascade;