# dcsid
分布式ID生成

## Install:

    go get github.com/cacos-group/dcsid

## 数据表

    CREATE TABLE `dcsid_alloc` (
      `id` bigint(20) NOT NULL AUTO_INCREMENT,
      `biz_tag` varchar(128) NOT NULL DEFAULT '',
      `max_id` bigint(20) NOT NULL DEFAULT '1',
      `step` int(11) NOT NULL,
      `description` varchar(256) DEFAULT NULL,
      `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
      PRIMARY KEY (`id`),
      unique `uniq_biz_tag` (`biz_tag`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8

## 使用
   
    dcsid := New(&Config{
		DSN:    "user:password@tcp(127.0.0.1:3306)/dcsid",
		BizTag: "test",
	})
	id, err := dcsid.NextId()
	if err != nil {
	    
	}
