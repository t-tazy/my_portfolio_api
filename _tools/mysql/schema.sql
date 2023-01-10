CREATE TABLE users
(
  id       BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ユーザーの識別子',
  name     VARCHAR(255) NOT NULL COMMENT 'ユーザー名',
  password VARCHAR(80) NOT NULL COMMENT 'パスワードハッシュ',
  role     VARCHAR(80) NOT NULL COMMENT 'ロール',
  created  DATETIME(6) NOT NULL COMMENT 'レコード作成日時',
  modified DATETIME(6) NOT NULL COMMENT 'レコード修正日時',
  PRIMARY KEY (id),
  UNIQUE KEY uix_name (name) USING BTREE
) Engine=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='ユーザー';

CREATE TABLE exercises
(
  id          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'エクササイズの識別子',
  user_id     BIGINT UNSIGNED NOT NULL COMMENT 'エクササイズを作成したユーザーの識別子',
  title       VARCHAR(128) NOT NULL COMMENT 'エクササイズのタイトル',
  description TEXT COMMENT 'エクササイズの説明',
  created     DATETIME(6) NOT NULL COMMENT 'レコード作成日時',
  modified    DATETIME(6) NOT NULL COMMENT 'レコード修正日時',
  PRIMARY KEY (id),
  CONSTRAINT fk_user_id
    FOREIGN KEY (user_id) REFERENCES users (id)
      ON DELETE RESTRICT ON UPDATE RESTRICT
) Engine=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='エクササイズ';
