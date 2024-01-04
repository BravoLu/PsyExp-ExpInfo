CREATE TABLE sub_infos (
  id INT NOT NULL AUTO_INCREMENT,
  sid VARCHAR(64) UNIQUE NOT NULL,
  eid INT NOT NULL,
  pid INT,
  state INT,
  finished_at DATETIME,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  deleted_at DATETIME,
  PRIMARY KEY (id),
  INDEX idx_sid (sid),
  INDEX idx_eid (eid),
  INDEX idx_pid (pid)
);
