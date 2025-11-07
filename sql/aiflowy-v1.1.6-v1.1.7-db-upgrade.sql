ALTER TABLE `tb_ai_plugin_tool`
ADD COLUMN `english_name` varchar(255) NULL COMMENT '英文名称' AFTER `debug_status`;