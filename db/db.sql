create extension if not exists citext;

drop table if exists users cascade;
create unlogged table users (
                       nickname citext primary key,
                       fullname text not null,
                       about text,
                       email citext unique not null
);


drop table if exists forum cascade;
create unlogged table forum (
                       title text,
                       "user" citext not null,
                       slug citext primary key,
                       posts bigint default 0,
                       threads bigint default 0,
                       foreign key ("user") references users(nickname)
);


drop table if exists thread cascade;
create unlogged table thread (
                        id serial primary key,
                        title text not null ,
                        author citext not null,
                        forum citext,
                        message text not null,
                        votes INT default 0,
                        slug citext unique,
                        created timestamptz default now(),
                        foreign key (author) references "users"(nickname),
                        foreign key (forum)  references "forum"(slug)
);


drop table if exists users_forum cascade;
create unlogged table users_forum (
                             nickname citext not null,
                             slug citext not null,
                             primary key (nickname, slug),
                             foreign key (nickname) references users(nickname),
                             foreign key (slug) references forum (slug)
);


drop table if exists post cascade;
create unlogged table post(
                     id bigserial primary key,
                     parent bigint default 0,
                     author citext not null,
                     message text not null,
                     isEdited boolean not null default false,
                     forum citext,
                     thread int,
                     created timestamptz default now(),
                     paths bigint[] default array[]::integer[],
                     foreign key (author) references users(nickname),
                     foreign key (forum) references forum(slug),
                     foreign key (thread) references thread(id)
);

drop table if exists vote cascade ;
create unlogged table vote (
                      id bigserial primary key,
                      nickname citext not null,
                      voice int not null,
                      thread int not null,
                      foreign key (nickname) references users(nickname),
                      foreign key (thread) references thread(id),
                      unique(thread, nickname)
);


create function add_votes() returns trigger as
$add_votes$
begin
    update thread set votes = (votes + new.voice) where id = new.thread;
    return new;
end
$add_votes$ language plpgsql;

create trigger after_insert_add_votes
    after insert
    on vote
    for each row
execute procedure add_votes();


create function update_votes() returns trigger as
$update_votes$
begin
    if old.voice != new.voice then
        update thread set votes = (votes + new.voice * 2) where id = new.thread;
    end if;
    return new;
end
$update_votes$ language plpgsql;

create trigger after_update_votes
    after update
    on vote
    for each row
execute procedure update_votes();


create function new_forum_user() returns trigger as
$new_forum_user$
begin
    insert into users_forum(nickname, slug)
    values (new.author, new.forum)
    on conflict do nothing;
    return new;
end
$new_forum_user$ language plpgsql;

create trigger after_thread_insert_update_user
    after insert
    on thread
    for each row
execute procedure new_forum_user();

create trigger after_post_insert
    after insert
    on post
    for each row
execute procedure new_forum_user();


create function update_post_paths() returns trigger as
$update_post_paths$
declare
    parent_path         bigint[];
    first_parent_thread int;
begin
    if (new.parent = 0) then
        new.paths := array_append(new.paths, new.id);
    else
        select paths from post where id = new.parent into parent_path;
        select thread from post where id = parent_path[1] into first_parent_thread;

        if not FOUND or first_parent_thread != new.thread then
            raise exception 'parent post was created in another thread'
            using errcode = '77777';
        end if;

        new.paths := new.paths || parent_path || new.id;
    end if;

    update forum
    set posts = posts + 1
    where forum.slug = new.forum;
    return new;
END
$update_post_paths$ language plpgsql;

create trigger before_post_insert
    before insert
    on post
    for each row
    execute procedure update_post_paths();


create function increment_threads_counter() returns trigger as
$increment_threads_counter$
begin
    update forum
    set threads = forum.threads + 1
    where slug = new.forum;
return new;
end
$increment_threads_counter$ language plpgsql;

create trigger after_insert_thread
    after insert
    on thread
    for each row
execute procedure increment_threads_counter();

create index if not exists users_email ON users(email);
create index if not exists users_forum_nickname ON users_forum(nickname);
create index if not exists users_forum_slug ON users_forum(slug);
create index if not exists thread_forum ON thread(forum);
create index if not exists thread_created ON thread(created);
create index if not exists thread_slug ON thread(slug);
create index if not exists thread_forum_created ON thread(forum, created);

create unique index if not exists vote_thread_nickname ON vote(thread, nickname);

create index if not exists post_thread_paths_id ON post(thread, paths, id);
create index if not exists post_thread_parent ON post(thread, parent);
create index if not exists paths1_id on post (paths, id);
create index if not exists post_paths1_paths_id ON post ((paths[1]), paths, id);

vacuum;
vacuum analyse;





