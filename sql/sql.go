package sql

const _updateSql = "UPDATE dcsid_alloc SET max_id=max_id+step WHERE biz_tag=?"
const _selectSql = "SELECT max_id, step FROM dcsid_alloc where biz_tag = ?"
const _insertSql = "INSERT INTO `dcsid_alloc` (`biz_tag`, `max_id`, `step`, `description`) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE `step` = VALUES(`step`);"
