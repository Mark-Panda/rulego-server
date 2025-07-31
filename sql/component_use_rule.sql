-- 创建索引
create sequence component_use_rule_seq increment by 1 minvalue 1 no maxvalue start with 1;

CREATE TABLE "public"."component_use_rule" (
    "id" bigint NOT NULL DEFAULT nextval('component_use_rule_seq'::regclass),
    "component_name" varchar(64) COLLATE "pg_catalog"."default",
    "component_type" varchar(64) COLLATE "pg_catalog"."default",
    "disabled" boolean DEFAULT false,
    "use_desc" text DEFAULT null,
    "use_rule_desc" text DEFAULT null,
    "created_at" timestamptz(6) NOT NULL DEFAULT now(),
    "updated_at" timestamptz(6) NOT NULL DEFAULT now(),
    CONSTRAINT "component_use_rule_pkey" PRIMARY KEY ("id")
);


COMMENT ON TABLE "public"."component_use_rule" IS '组件使用规则配置表';

CREATE UNIQUE INDEX component_use_rule_component_name_unique_idx ON component_use_rule(component_name);


COMMENT ON COLUMN "public"."component_use_rule"."id" IS '主键ID';
COMMENT ON COLUMN "public"."component_use_rule"."component_name" IS '组件名称';
COMMENT ON COLUMN "public"."component_use_rule"."component_type" IS '组件类型';
COMMENT ON COLUMN "public"."component_use_rule"."disabled" IS '是否禁用';
COMMENT ON COLUMN "public"."component_use_rule"."use_desc" IS '使用描述';
COMMENT ON COLUMN "public"."component_use_rule"."use_rule_desc" IS '使用规则描述';
COMMENT ON COLUMN "public"."component_use_rule"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."component_use_rule"."updated_at" IS '更新时间';
