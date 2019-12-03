ALTER TYPE todotype ADD VALUE '';
alter table todo
    alter column type set default '';