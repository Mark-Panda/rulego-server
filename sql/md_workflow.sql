-- 创建索引
create sequence md_workflow_seq increment by 1 minvalue 1 no maxvalue start with 1;

CREATE TABLE "public"."md_workflow" (
    "id" bigint NOT NULL DEFAULT nextval('md_workflow_seq'::regclass),
    "title" varchar(255) COLLATE "pg_catalog"."default",
    "content" text DEFAULT null,
    "desc" varchar(255) COLLATE "pg_catalog"."default",
    "chain_id" varchar(255) COLLATE "pg_catalog"."default",
    "chain_name" varchar(255) COLLATE "pg_catalog"."default",
    "chain_version" int4 NOT NULL DEFAULT 0,
    "created_at" timestamptz(6) NOT NULL DEFAULT now(),
    "updated_at" timestamptz(6) NOT NULL DEFAULT now(),
    CONSTRAINT "md_workflow_pkey" PRIMARY KEY ("id")
);


COMMENT ON TABLE "public"."md_workflow" IS 'md工作流';


COMMENT ON COLUMN "public"."md_workflow"."id" IS '主键ID';
COMMENT ON COLUMN "public"."md_workflow"."title" IS 'md标题';
COMMENT ON COLUMN "public"."md_workflow"."desc" IS 'md描述';
COMMENT ON COLUMN "public"."md_workflow"."content" IS '内容';
COMMENT ON COLUMN "public"."md_workflow"."chain_id" IS '规则链ID';
COMMENT ON COLUMN "public"."md_workflow"."chain_name" IS '规则链名称';
COMMENT ON COLUMN "public"."md_workflow"."chain_version" IS '规则链版本';
COMMENT ON COLUMN "public"."md_workflow"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."md_workflow"."updated_at" IS '更新时间';