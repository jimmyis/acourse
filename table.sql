create table users (
  id string not null,
  username string not null,
  name string not null,
  email string,
  about_me string not null default '',
  image string not null default '',
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),
  primary key (id)
);
create unique index on users (username);
create unique index on users (email);
create index on users (created_at desc);

create table roles (
  user_id string,
  admin bool not null default false,
  instructor bool not null default false,
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),
  primary key (user_id),
  foreign key (user_id) references users (id)
);
create index on roles (admin);
create index on roles (instructor);

create table courses (
  id serial,
  user_id string not null,
  title string not null,
  short_desc string not null,
  long_desc string not null,
  image string not null,
  start timestamp default null,
  url string default null,
  type int not null default 0,
  price decimal(9,2) not null default 0,
  discount decimal(9,2) default 0,
  enroll_detail string not null default '',
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),
  primary key (id),
  foreign key (user_id) references users (id)
);
create unique index on courses (url);
create index on courses (created_at desc);
create index on courses (updated_at desc);

create table course_options (
  course_id int,
  public bool not null default false,
  enroll bool not null default false,
  attend bool not null default false,
  assignment bool not null default false,
  discount bool not null default false,
  primary key (course_id),
  foreign key (course_id) references courses (id)
);
create index on course_options (public);
create index on course_options (enroll);
create index on course_options (public, enroll);
create index on course_options (public, discount);
create index on course_options (public, discount, enroll);

create table course_contents (
  id serial,
  course_id uuid not null,
  i int not null default 0,
  title string not null default '',
  long_desc string not null default '',
  video_id string not null default '',
  video_type int not null default 0,
  download_url string not null default '',
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),
  primary key (id),
  foreign key (course_id) references courses (id)
);
create index on course_contents (course_id, i);

create table assignments (
  id serial,
  course_id uuid not null,
  i int not null,
  title string not null,
  long_desc string not null,
  open bool not null default false,
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),
  primary key (id),
  foreign key (course_id) references courses (id)
);
create index on assignments (course_id, i);

create table user_assignments (
  id serial,
  user_id string not null,
  assignment_id uuid not null,
  download_url string not null,
  created_at timestamp not null default now(),
  primary key (id),
  foreign key (user_id) references users (id),
  foreign key (assignment_id) references assignments (id)
);
create index on user_assignments (created_at);

create table enrolls (
  user_id string,
  course_id uuid not null,
  created_at timestamp not null default now(),
  primary key (user_id, course_id),
  foreign key (user_id) references users (id),
  foreign key (course_id) references courses (id)
);
create index on enrolls (created_at);
create index on enrolls (user_id, created_at);
create index on enrolls (course_id, created_at);

create table attends (
  id serial,
  user_id string not null,
  course_id uuid not null,
  created_at timestamp not null default now(),
  primary key (id),
  foreign key (user_id) references users (id),
  foreign key (course_id) references courses (id)
);
create index on attends (created_at);
create index on attends (user_id, created_at);
create index on attends (course_id, created_at);
create index on attends (user_id, course_id, created_at);

create table payments (
  id serial,
  user_id string not null,
  course_id uuid not null,
  image string not null,
  price decimal(9, 2) not null,
  original_price decimal(9, 2) not null,
  code string not null,
  status int not null,
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),
  at timestamp default null,
  primary key (id),
  foreign key (user_id) references users (id),
  foreign key (course_id) references courses (id)
);
create index on payments (created_at desc);
create index on payments (code);
create index on payments (course_id, code);
create index on payments (status, created_at desc);
