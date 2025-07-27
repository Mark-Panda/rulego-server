-- 创建索引
create sequence run_log_seq increment by 1 minvalue 1 no maxvalue start with 1;

CREATE TABLE "public"."run_log" (
    "id" bigint NOT NULL DEFAULT nextval('run_log_seq'::regclass),
    "run_id" varchar(255) COLLATE "pg_catalog"."default",
    "chain_id" varchar(255) COLLATE "pg_catalog"."default",
    "chain_name" varchar(255) COLLATE "pg_catalog"."default",
    "node_log" text DEFAULT null,
    "additional_info" text DEFAULT null,
    "rule_chain_info" text DEFAULT null,
    "metadata" text DEFAULT null,
    "start_ts" bigint,
    "end_ts" bigint,
    "created_at" timestamptz(6) NOT NULL DEFAULT now(),
    "updated_at" timestamptz(6) NOT NULL DEFAULT now(),
    CONSTRAINT "run_log_pkey" PRIMARY KEY ("id")
);


COMMENT ON TABLE "public"."run_log" IS '运行日志表';

CREATE UNIQUE INDEX run_log_rule_chain_id_unique_idx ON run_log(run_id);

COMMENT ON COLUMN "public"."run_log"."id" IS '主键ID';
COMMENT ON COLUMN "public"."run_log"."run_id" IS '运行ID';
COMMENT ON COLUMN "public"."run_log"."chain_id" IS '规则链ID';
COMMENT ON COLUMN "public"."run_log"."chain_name" IS '规则链名称';
COMMENT ON COLUMN "public"."run_log"."node_log" IS '节点运行日志';
COMMENT ON COLUMN "public"."run_log"."additional_info" IS '附加信息';
COMMENT ON COLUMN "public"."run_log"."rule_chain_info" IS '规则链信息';
COMMENT ON COLUMN "public"."run_log"."metadata" IS '元数据';
COMMENT ON COLUMN "public"."run_log"."start_ts" IS '开始时间';
COMMENT ON COLUMN "public"."run_log"."end_ts" IS '结束时间';
COMMENT ON COLUMN "public"."run_log"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."run_log"."updated_at" IS '更新时间';