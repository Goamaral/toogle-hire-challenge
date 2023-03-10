CREATE TABLE IF NOT EXISTS questions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	body TEXT NOT NULL,
	author_id INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS question_options (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	body TEXT NOT NULL,
  correct INTEGER NOT NULL,
  question_id INTEGER NOT NULL,
  FOREIGN KEY(question_id) REFERENCES questions(id) ON DELETE CASCADE
);