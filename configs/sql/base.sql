create extension if not exists citext;

create table if not exists users (
    nickname citext primary key,
    fullname text default '',
    about text default '',
    email citext unique not null
);

create table if not exists forums (
    owner citext not null,
    slug citext unique not null primary key,
    title text not null,
    threads integer default 0,
    posts integer default 0,

    foreign key (owner) references users(nickname)
);

create table if not exists forums_users (
    forum_slug citext not null  collate "ucs_basic",
    nickname citext not null  collate "ucs_basic",
    primary key (forum_slug, nickname)
);

create table if not exists threads (
    id      serial primary key,
    author  citext not null,
    title   text not null,
    message text not null,
    votes   integer default 0,
    slug   citext unique,
    created timestamp with time zone default current_timestamp,
    forum   citext,

    foreign key(author) references users(nickname),
    foreign key(forum) references forums(slug)
);


create table if not exists posts (
    id         serial primary key,
    parent     int    default 0 ,
    rootParent int  default 0 ,
    post_path  bigint [] not null  default '{}'::bigint[],
    message    text not null,
    isEdit     boolean default false,
    forum      citext,
    created    timestamp with time zone default current_timestamp,
    thread     bigint default 0,
    author     citext,

    foreign key (author) references users (nickname)
);


create table if not exists votes (
  author citext,
  thread int,
  vote   int default 1,

  foreign key(author) references users(nickname),
  unique (thread, author)
);

create index if not exists threads_forum_created_idx on threads (forum,created);
create index if not exists posts_thread_id_idx on posts (thread, id);
create index if not exists posts_thread_mpath_idx on posts (thread, post_path);

CREATE OR REPLACE FUNCTION upd_forums_users() RETURNS TRIGGER AS $upd_forums_users$
    BEGIN
        insert into forums_users values (new.forum, (select nickname from users where new.author = nickname)) on conflict do nothing;
        RETURN NEW;  -- возвращаемое значение для триггера AFTER игнорируется
    END;
$upd_forums_users$ LANGUAGE plpgsql;

CREATE TRIGGER upd_forums_users_trg
AFTER INSERT ON threads
    FOR EACH ROW EXECUTE PROCEDURE upd_forums_users();

CREATE OR REPLACE FUNCTION upd_forums_users_posts() RETURNS TRIGGER AS $upd_forums_users_posts$
    BEGIN
        insert into forums_users values (new.forum, (select nickname from users where new.author = nickname)) on conflict do nothing;
        RETURN NEW;  -- возвращаемое значение для триггера AFTER игнорируется
    END;
$upd_forums_users_posts$ LANGUAGE plpgsql;

CREATE TRIGGER upd_forums_users_posts_trg
AFTER INSERT ON posts
    FOR EACH ROW EXECUTE PROCEDURE upd_forums_users_posts();

CREATE OR REPLACE FUNCTION upd_forum_threads() RETURNS TRIGGER AS $upd_forum_threads$
    BEGIN
        update forums set threads = threads + 1 where  slug = new.forum;
        RETURN NEW;  -- возвращаемое значение для триггера AFTER игнорируется
    END;
$upd_forum_threads$ LANGUAGE plpgsql;

CREATE TRIGGER upd_forum_threads
AFTER INSERT ON threads
    FOR EACH ROW EXECUTE PROCEDURE upd_forum_threads();



CREATE OR REPLACE FUNCTION upd_forum_posts() RETURNS TRIGGER AS $upd_forum_posts$
    BEGIN
        update forums set posts = posts + 1 where  slug = new.forum;
        RETURN NEW;  -- возвращаемое значение для триггера AFTER игнорируется
    END;
$upd_forum_posts$ LANGUAGE plpgsql;

CREATE TRIGGER upd_forum_posts
AFTER INSERT ON posts
    FOR EACH ROW EXECUTE PROCEDURE upd_forum_posts();


create or replace function upd_votes()
  returns trigger as $upd_votes$
begin
  if (tg_op = 'INSERT')
  then
    update threads set votes = votes + new.vote where id = new.thread;
  elseif (tg_op = 'UPDATE')
    then
      update threads set votes = votes + (new.vote - old.vote) where id = new.thread;
  end if;
  return new;
end;
$upd_votes$
language plpgsql;

CREATE TRIGGER upd_votes_trg
AFTER INSERT OR UPDATE ON votes
    FOR EACH ROW EXECUTE PROCEDURE upd_votes();

GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO postgres;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO postgres;