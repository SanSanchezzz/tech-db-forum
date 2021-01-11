\c forums;

drop table if exists users, forums, threads, posts, votes;

create extension if not exists citext;

create unlogged table users (
                                nickname citext collate "C" primary key,
                                fullname varchar(256) not null,
                                about text,
                                email citext unique not null
);

create unlogged table forums (
                                 slug citext primary key,
                                 title text not null,
                                 nickname citext not null,
                                 foreign key (nickname) references users (nickname) on delete cascade,
                                 threads int not null default 0,
                                 posts int not null default 0
);

create unlogged table threads (
                                  id serial primary key,
                                  title varchar(256) not null,
                                  nickname citext not null,
                                  foreign key (nickname) references users (nickname) on delete cascade,
                                  forum citext not null,
                                  foreign key (forum) references forums (slug) on delete cascade,
                                  message text not null,
                                  votes int not null default 0,
                                  slug citext not null,
                                  created timestamptz not null
);

create unlogged table posts (
                                id serial primary key ,
                                parent integer not null,
                                author citext not null,
                                foreign key (author) references users (nickname) on delete cascade,
                                message text not null,
                                is_edited boolean default false,
                                forum citext not null,
                                foreign key (forum) references forums (slug) on delete cascade,
                                thread integer not null,
                                foreign key (thread) references threads (id) on delete cascade,
                                created timestamptz not null default now(),
                                path int[] not null
);

create unlogged table votes (
                                nickname citext not null,
                                foreign key (nickname) references users (nickname) on delete cascade,
                                thread integer not null,
                                foreign key (thread) references threads (id) on delete cascade,
                                voice integer not null,
                                primary key (nickname, thread)
);

create unlogged table forum_user (
                                     nickname citext,
                                     foreign key (nickname) references users (nickname) on delete cascade,
                                     forum citext,
                                     foreign key (forum) references forums (slug) on delete cascade,
                                     primary key (nickname, forum)
);

--
create index on users using hash (email);
create index on users using hash (nickname);

create index on posts (id, thread);
-- create index on posts (path, id);
-- create index on posts (path, created, id);
-- create index on posts (created, id);
create index on posts (thread);
create index on posts (path);
create index on posts(forum);
-- create index on posts(author);

create index on forums USING hash (slug);
create index on forums (nickname);

create index on threads USING hash (id);
create index on threads USING hash (slug);
create index on threads (forum, created);
--

create or replace function vote_after() returns trigger as
$$
begin
    update threads set votes = votes + new.voice where threads.id = new.thread;
    return new;
end;
$$ language plpgsql;

create trigger vote_after after insert on votes
    for each row execute procedure vote_after();

create or replace function vote_update() returns trigger as
$$
begin
    if new.voice != old.voice then
        update threads set votes = votes + new.voice - old.voice where threads.id = new.thread;
    end if;
    return old;
end;
$$ language plpgsql;

create trigger vote_update after update on votes
    for each row execute procedure vote_update();

create or replace function post_path() returns trigger as
$$
begin
    if new.path[1] <> 0 then
        new.path := (
                        select path
                        from posts
                        where id = new.path[1] and thread = new.thread
                    ) || array [new.id];
        if cardinality(new.path) = 1 then
            raise 'Parent exception' using errcode = '00404';
        end if;
    else
        new.path[1] := new.id;
    end if;
    return new;
end;
$$ language plpgsql;

drop trigger if exists post_path on posts;

create trigger post_path before insert on posts
    for each row execute procedure post_path();

create or replace function forums_thread_insert() returns trigger as
$$
begin
    update forums set threads = threads + 1 where slug = new.forum;
    insert into forum_user values (new.nickname, new.forum) on conflict do nothing;
    return new;
end;
$$ language plpgsql;

drop trigger if exists thread_count on threads;
create trigger thread_count after insert on threads
    for each row execute procedure forums_thread_insert();

create or replace function forums_post_insert() returns trigger as
$$
begin
    update forums set posts = posts + 1 where slug = new.forum;
    insert into forum_user values (new.author, new.forum) on conflict do nothing;
    return new;
end;
$$ language plpgsql;

drop trigger if exists post_count on posts;
create trigger post_count after insert on posts
    for each row execute procedure forums_post_insert();
