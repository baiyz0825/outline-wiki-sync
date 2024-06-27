CREATE TABLE IF NOT EXISTS file_sync_record
(
    id              INT AUTO_INCREMENT PRIMARY KEY COMMENT '记录ID',
    outline_wiki_id VARCHAR(255) COMMENT '大纲Wiki ID',
    collection_id   VARCHAR(255) COMMENT '集合ID',
    file_name       VARCHAR(255) COMMENT '文件名',
    file_size       VARCHAR(255) COMMENT '文件大小',
    file_path       VARCHAR(255) COMMENT '文件路径',
    file_content    TEXT COMMENT '文件内容',
    sync            BOOLEAN COMMENT '同步标志',
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted         BOOLEAN COMMENT '删除标志',
    INDEX idx_outline_wiki_id (outline_wiki_id),
    INDEX idx_collection_id (collection_id),
    INDEX idx_file_name (file_name)
) COMMENT ='文件同步记录表';


CREATE TABLE IF NOT EXISTS outline_wiki_collection_mapping
(
    id              INT AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    collection_id   VARCHAR(255) COMMENT '集合ID',
    collection_path VARCHAR(255) COMMENT '集合路径',
    collection_name VARCHAR(255) COMMENT '集合名称',
    sync            BOOLEAN COMMENT '同步标志',
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted         BOOLEAN COMMENT '删除标志',
    INDEX idx_collection_id (collection_id),
    INDEX idx_collection_path (collection_path),
    INDEX idx_collection_name (collection_name)
) COMMENT ='大纲Wiki集合映射表';