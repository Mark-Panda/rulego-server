-- 创建索引
create sequence regulation_seq increment by 1 minvalue 1 no maxvalue start with 1;

CREATE TABLE "public"."regulation" (
    "id" bigint NOT NULL DEFAULT nextval('regulation_seq'::regclass),
    "user_name" varchar(64) COLLATE "pg_catalog"."default",
    "root" boolean,
    "disabled" boolean,
    "name" varchar(64) COLLATE "pg_catalog"."default",
    "rule_chain_id" varchar(64) COLLATE "pg_catalog"."default",
    "rule_config" text DEFAULT null,
    "created_at" timestamptz(6) NOT NULL DEFAULT now(),
    "updated_at" timestamptz(6) NOT NULL DEFAULT now(),
    "deleted_at" timestamptz(6),
    CONSTRAINT "regulation_pkey" PRIMARY KEY ("id")
);


COMMENT ON TABLE "public"."regulation" IS '规则配置表';

CREATE UNIQUE INDEX regulation_rule_chain_id_unique_idx ON regulation(rule_chain_id);

COMMENT ON COLUMN "public"."regulation"."id" IS '主键ID';
COMMENT ON COLUMN "public"."regulation"."user_name" IS '用户名';
COMMENT ON COLUMN "public"."regulation"."root" IS '是否根节点';
COMMENT ON COLUMN "public"."regulation"."disabled" IS '是否禁用';
COMMENT ON COLUMN "public"."regulation"."name" IS '规则名称';
COMMENT ON COLUMN "public"."regulation"."rule_chain_id" IS '规则ID';
COMMENT ON COLUMN "public"."regulation"."rule_config" IS '规则配置信息';
COMMENT ON COLUMN "public"."regulation"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."regulation"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."regulation"."deleted_at" IS '删除时间';