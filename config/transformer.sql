CREATE TABLE IF NOT EXISTS jobs (
  job_id VARCHAR(50) NOT NULL,
  status TINYINT(1) DEFAULT 1 NOT NULL,
  PRIMARY KEY (job_id)
) ENGINE=InnoDB AUTO_INCREMENT=500 DEFAULT CHARSET=utf8;