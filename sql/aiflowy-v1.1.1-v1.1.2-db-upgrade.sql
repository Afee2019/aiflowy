SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for tb_ai_bot_api_key since v1.1.2
-- ----------------------------
DROP TABLE IF EXISTS `tb_ai_bot_api_key`;
CREATE TABLE `tb_ai_bot_api_key`  (
  `id` bigint NOT NULL COMMENT 'id',
  `api_key` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT 'apiKey，请勿手动修改！',
  `bot_id` bigint NOT NULL COMMENT 'botId',
  `salt` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '加密botId，生成apiKey的盐',
  `options` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '预留拓展配置的字段',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = 'bot apiKey 表' ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;