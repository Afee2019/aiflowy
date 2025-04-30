INSERT INTO `tb_ai_bot` VALUES (274724831961026560, 1, 1000000, '天气查询', '天气查询', NULL, 273418062806761472, '{\"systemPrompt\":\"你是天气查询助手\"}', NULL, '2025-04-29 15:19:24', 1, '2025-04-29 15:19:24', 1);



INSERT INTO `tb_ai_bot_knowledge` VALUES (6, 267746165268017152, 267751447398232064, NULL);
INSERT INTO `tb_ai_bot_knowledge` VALUES (10, 269226938047168512, 273274084509630464, NULL);
INSERT INTO `tb_ai_bot_knowledge` VALUES (11, 267778325987205120, 273416394929192960, NULL);






INSERT INTO `tb_ai_bot_plugins` VALUES (267770208130498560, 267769906283216896, 267769494146711552, NULL);
INSERT INTO `tb_ai_bot_plugins` VALUES (274649868264542208, 269226938047168512, 267769494146711552, NULL);
INSERT INTO `tb_ai_bot_plugins` VALUES (274727326384918528, 274724831961026560, 273726420990324736, NULL);
INSERT INTO `tb_ai_bot_plugins` VALUES (274728813039194112, 267848016181075968, 267769494146711552, NULL);



INSERT INTO `tb_ai_chat_topic` VALUES (269265123934486528, 1, '未命名', '2025-04-14 13:44:28', NULL);


INSERT INTO `tb_ai_llm` VALUES (273418062806761472, 1, 1000000, '星火大模型', 'spark', NULL, NULL, 1, 1, 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 'wss://spark-api.xf-yun.com/v3.5/chat', 'wss://spark-api.xf-yun.com/v3.5/chat', '', '', NULL);


INSERT INTO `tb_ai_plugin` VALUES (273726420990324736, '和风天气', '每日天气预报API，提供全球城市未来3-30天天气预报，包括：日出日落、月升月落、最高最低温度、天气白天和夜间状况、风力、风速、风向、相对湿度、大气压强、降水量、露点温度、紫外线强度、能见度等。', NULL, 'https://jw54e3j6ha.re.qweatherapi.com', 'apiKey', '2025-04-26 21:12:04', '/attachment/2025/04-26/cfa688d9-9bc5-4d3f-a50c-2501604e8bed.jpg', 'query', '[{\"label\":\"X-QW-Api-Key\",\"value\":\"2ad5de1c5dddd44e589bece44fdad4258\"}]', 'q', '', NULL, NULL, NULL);


INSERT INTO `tb_ai_plugin_tool` VALUES (273931608955039744, 273726420990324736, 'get_weather', '查询某个城市地方未来三天的天气', '/v7/weather/3d', '2025-04-27 10:47:25', NULL, '[{\"name\":\"location\",\"description\":\"城市 LocationID（如 101010100 代表北京），可通过城市列表查询，当用户问及例如北京的天气的时候，大模型需要查询北京转换为对应的城市id\",\"type\":\"String\",\"method\":\"Query\",\"required\":true,\"defaultValue\":\"101010100\",\"enabled\":true,\"key\":\"1745905802390\"}]', NULL, 'Get', 0, 0);
INSERT INTO `tb_ai_plugin_tool` VALUES (275035209345548290, 273726420990324736, 'get_cityId', '城市搜索API提供全球地理位位置、全球城市搜索服务，支持经纬度坐标反查、多语言、模糊搜索等功能。', '/geo/v2/city/lookup', '2025-04-30 11:52:44', NULL, '[{\"name\":\"location\",\"description\":\"需要查询地区的名称，支持文字、以英文逗号分隔的经度,纬度坐标（十进制，最多支持小数点后两位）、LocationID或Adcode（仅限中国城市）。例如 location=北京 或 location=116.41,39.92\",\"type\":\"String\",\"method\":\"Query\",\"required\":true,\"defaultValue\":\"\",\"enabled\":true,\"key\":\"1745985221077\"}]', NULL, 'Get', 0, 0);



INSERT INTO `tb_ai_plugins` VALUES (267769494146711552, 1, 1000000, 1, 'get_ip_info', '根据ip获取区域信息', '{\"method\":\"GET\",\"url\":\"https://qifu-api.baidubce.com/ip/geo/v1/district\",\"headers\":[],\"params\":[{\"key\":\"ip\",\"value\":\"\",\"desc\":\"IP地址\"}],\"body\":{}}', 1, '2025-04-10 10:41:22', 1, '2025-04-29 10:21:07', 1, '', 0, NULL);
INSERT INTO `tb_ai_plugins` VALUES (269229432299077632, 1, 1000000, 1, 'test', 'test', NULL, 1, '2025-04-14 11:22:39', 1, '2025-04-14 11:22:39', 1, '', 0, '/attachment/2025/04-14/e2f3e5a5-a4fa-4abe-8fe6-49d10c3d8dbe.jpeg');




INSERT INTO `tb_sys_account` VALUES (1, 1, 1000000, 'admin', '$2a$10$S8HVnrS8m7iygQBS7r1dYuOstEUl5q/W1yhgFcS1uyL6o2/23yUYO', 99, '超级管理员', '15555555555', 'aaa@qq.com', '/attachment/2025/04-10/59866709-5bc5-4e9f-9445-ecb603ff2d82.jpg', 1, NULL, 1, '2025-04-10 16:33:48', 1, '2025-04-10 17:56:17', 1, '', 0);


INSERT INTO `tb_sys_account_position` VALUES (267858187553452032, 1, 259067270360543232);


INSERT INTO `tb_sys_account_role` VALUES (267858187456983040, 1, 1);



INSERT INTO `tb_sys_api_key` VALUES (275081657634275328, '8e8ca1ff71d1430d9ad2a074f5901ed8', '2025-04-30 14:57:18', 1, 1, 1000000, '2025-05-30 00:00:00', 1);



INSERT INTO `tb_sys_dept` VALUES (1, 1000000, 0, '0', '总公司', 'root_dept', 0, 1, '2025-03-17 09:09:57', 1, '2025-03-17 09:10:00', 1, '', 0);



INSERT INTO `tb_sys_dict` VALUES (268213400717598720, 'test', 'test', 'test', 1, NULL, NULL, '{}', '2025-04-11 16:05:18', '2025-04-11 16:05:18');





INSERT INTO `tb_sys_menu` VALUES (258052082618335232, 0, 0, '系统管理', '', '', 'ControlOutlined', 1, '', 999, 0, '2025-03-14 15:07:51', 1, '2025-04-08 11:12:00', 1, '', 0);
INSERT INTO `tb_sys_menu` VALUES (258052774330368000, 258052082618335232, 0, '用户管理', '/sys/sysAccount', '', 'UserOutlined', 1, '', 0, 0, '2025-03-14 15:10:36', 1, '2025-03-14 15:10:36', 1, '', 0);
INSERT INTO `tb_sys_menu` VALUES (258075705244676096, 258052082618335232, 0, '角色管理', '/sys/sysRole', '', 'FrownOutlined', 1, '', 11, 0, '2025-03-14 16:41:43', 1, '2025-03-14 16:41:43', 1, '', 0);
INSERT INTO `tb_sys_menu` VALUES (258075850434703360, 258052082618335232, 0, '菜单管理', '/sys/sysMenu', '', 'ApartmentOutlined', 1, '', 21, 0, '2025-03-14 16:42:18', 1, '2025-03-14 16:42:18', 1, '', 0);
INSERT INTO `tb_sys_menu` VALUES (258077347000119296, 258052774330368000, 1, '用户查询', '', '', NULL, 0, 'sysUser:list', 1, 0, '2025-03-14 16:48:14', 1, '2025-03-14 16:48:14', 1, '', 0);
INSERT INTO `tb_sys_menu` VALUES (258079445137780736, 258052774330368000, 1, '用户保存', '', '', NULL, 1, 'sysUser:save', 2, 0, '2025-03-14 16:56:35', 1, '2025-03-14 16:56:35', 1, '', 0);
INSERT INTO `tb_sys_menu` VALUES (259048038847483904, 258052082618335232, 0, '部门管理', '/sys/sysDept', '', 'BankOutlined', 1, '', 31, 0, '2025-03-17 09:05:25', 1, '2025-03-17 09:05:25', 1, '', 0);
INSERT INTO `tb_sys_menu` VALUES (259065916854448128, 258052082618335232, 0, '岗位管理', '/sys/sysPosition', '', 'AuditOutlined', 1, '', 32, 0, '2025-03-17 10:16:28', 1, '2025-03-17 10:16:28', 1, '', 0);
INSERT INTO `tb_sys_menu` VALUES (259166589327630336, 0, 0, '欢迎回来', '/index', '', 'BorderlessTableOutlined', 1, '', 0, 0, '2025-03-17 16:56:30', 1, '2025-03-17 16:56:30', 1, '', 1);
INSERT INTO `tb_sys_menu` VALUES (259168688849412096, 0, 0, '系统配置', '', '', NULL, 1, '', 99999, 0, '2025-03-17 17:04:51', 1, '2025-03-17 17:04:51', 1, '', 0);
INSERT INTO `tb_sys_menu` VALUES (259168810450673664, 259168688849412096, 0, 'AI大模型', '/config/model', '', 'DeploymentUnitOutlined', 1, '', 0, 0, '2025-03-17 17:05:20', 1, '2025-03-17 17:05:20', 1, '', 1);
INSERT INTO `tb_sys_menu` VALUES (259168916721754112, 259168688849412096, 0, '系统设置', 'sys/settings', '', 'SettingOutlined', 1, '', 12, 0, '2025-03-17 17:05:45', 1, '2025-04-17 15:54:39', 1, '', 0);
INSERT INTO `tb_sys_menu` VALUES (259169177682960384, 258052082618335232, 0, '数据字典', '/sys/dicts', '', 'AppstoreOutlined', 1, '', 51, 0, '2025-03-17 17:06:47', 1, '2025-04-14 13:35:41', 1, '', 0);
INSERT INTO `tb_sys_menu` VALUES (259169318720626688, 258052082618335232, 0, '日志管理', '/sys/logs', '', 'ExclamationCircleOutlined', 1, '', 61, 0, '2025-03-17 17:07:21', 1, '2025-04-14 13:35:34', 1, '', 0);
INSERT INTO `tb_sys_menu` VALUES (259169540360232960, 0, 0, 'AI能力', '', '', NULL, 1, '', 21, 0, '2025-03-17 17:08:14', 1, '2025-03-17 17:08:14', 1, '', 0);
INSERT INTO `tb_sys_menu` VALUES (259169689438380032, 259169540360232960, 0, '大模型商店', '/ai/ollama', '', 'SlackSquareFilled', 1, '', 1111, 0, '2025-03-17 17:08:49', 1, '2025-04-08 12:24:03', 1, '', 0);
INSERT INTO `tb_sys_menu` VALUES (259169837824466944, 259169540360232960, 0, 'Bots', '/ai/bots', '', 'TwitchOutlined', 1, '', 11, 0, '2025-03-17 17:09:24', 1, '2025-03-17 17:09:24', 1, '', 0);
INSERT INTO `tb_sys_menu` VALUES (259169982154661888, 259169540360232960, 0, '插件', '/ai/plugin', '', 'RobotOutlined', 1, '', 21, 0, '2025-03-17 17:09:59', 1, '2025-04-25 11:23:15', 1, '', 0);
INSERT INTO `tb_sys_menu` VALUES (259170117110587392, 259169540360232960, 0, '工作流', '/ai/workflow', '', 'BranchesOutlined', 1, '', 31, 0, '2025-03-17 17:10:31', 1, '2025-03-17 17:10:31', 1, '', 0);
INSERT INTO `tb_sys_menu` VALUES (259170422338478080, 259169540360232960, 0, '知识库', '/ai/knowledge', '', 'ReadOutlined', 1, '', 51, 0, '2025-03-17 17:11:44', 1, '2025-03-17 17:11:44', 1, '', 0);
INSERT INTO `tb_sys_menu` VALUES (259170538264846336, 259169540360232960, 0, '大模型', '/ai/llms', '', 'SlidersOutlined', 1, '', 61, 0, '2025-03-17 17:12:11', 1, '2025-04-14 10:49:05', 1, '', 0);
INSERT INTO `tb_sys_menu` VALUES (269220820365377536, 259170422338478080, 0, '知识库', '/ai/knowledge/document', '', 'ZhihuSquareFilled', 1, '', 1, 0, '2025-04-14 10:48:25', 1, '2025-04-14 10:50:19', 1, '', 1);
INSERT INTO `tb_sys_menu` VALUES (269221948243083264, 259170422338478080, 1, '知识库', '/ai/knowledge/Document', '', 'BookFilled', 1, '', 1, 0, '2025-04-14 10:52:54', 1, '2025-04-14 10:52:54', 1, '', 1);
INSERT INTO `tb_sys_menu` VALUES (270761213536096256, 258052082618335232, 0, 'apiKey', '/sys/sysApiKey', '', 'PoundOutlined', 1, '', 22, 0, '2025-04-18 16:49:24', 1, '2025-04-30 14:55:15', 1, '', 0);





INSERT INTO `tb_sys_position` VALUES (259067270360543232, 1000000, 1, '总部CTO', '', 0, 1, '2025-03-17 10:21:50', 1, '2025-03-17 10:21:50', 1, '', 0);



INSERT INTO `tb_sys_role` VALUES (1, 1000000, '超级管理员', 'super_admin', 1, '2025-03-14 14:52:37', 1, '2025-03-14 14:52:37', 1, '', 0);



INSERT INTO `tb_sys_role_menu` VALUES (259111372649250817, 1, 258052774330368000);
INSERT INTO `tb_sys_role_menu` VALUES (259111372649250818, 1, 258075705244676096);
INSERT INTO `tb_sys_role_menu` VALUES (259111372649250819, 1, 258075850434703360);
INSERT INTO `tb_sys_role_menu` VALUES (259111372649250820, 1, 258077347000119296);
INSERT INTO `tb_sys_role_menu` VALUES (259111372649250821, 1, 258079445137780736);
INSERT INTO `tb_sys_role_menu` VALUES (259111372649250822, 1, 259048038847483904);
INSERT INTO `tb_sys_role_menu` VALUES (259111372649250824, 1, 259065916854448128);
INSERT INTO `tb_sys_role_menu` VALUES (259111372649250825, 1, 258052082618335232);
INSERT INTO `tb_sys_role_menu` VALUES (259168688966852608, 1, 259168688849412096);
INSERT INTO `tb_sys_role_menu` VALUES (259168916826611712, 1, 259168916721754112);
INSERT INTO `tb_sys_role_menu` VALUES (259169177817178112, 1, 259169177682960384);
INSERT INTO `tb_sys_role_menu` VALUES (259169318829678592, 1, 259169318720626688);
INSERT INTO `tb_sys_role_menu` VALUES (259169540477673472, 1, 259169540360232960);
INSERT INTO `tb_sys_role_menu` VALUES (259169689576792064, 1, 259169689438380032);
INSERT INTO `tb_sys_role_menu` VALUES (259169837941907456, 1, 259169837824466944);
INSERT INTO `tb_sys_role_menu` VALUES (259169982280491008, 1, 259169982154661888);
INSERT INTO `tb_sys_role_menu` VALUES (259170117223833600, 1, 259170117110587392);
INSERT INTO `tb_sys_role_menu` VALUES (259170422447529984, 1, 259170422338478080);
INSERT INTO `tb_sys_role_menu` VALUES (259170538378092544, 1, 259170538264846336);
INSERT INTO `tb_sys_role_menu` VALUES (270761213603205120, 1, 270761213536096256);


