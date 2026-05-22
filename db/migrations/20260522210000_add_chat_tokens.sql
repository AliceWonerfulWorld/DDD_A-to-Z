-- Create "chat_tokens" table
CREATE TABLE "chat_tokens" ("token_hash" text NOT NULL, "user_id" text NOT NULL, "guild_id" text NOT NULL, "expires_at" timestamptz NOT NULL, "used_at" timestamptz NULL, "created_at" timestamptz NOT NULL, PRIMARY KEY ("token_hash"), CONSTRAINT "chat_tokens_guild_id_fkey" FOREIGN KEY ("guild_id") REFERENCES "guilds" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT "chat_tokens_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- Create index "chat_tokens_expires_at_idx" to table: "chat_tokens"
CREATE INDEX "chat_tokens_expires_at_idx" ON "chat_tokens" ("expires_at");
-- Create "guild_chat_messages" table
CREATE TABLE "guild_chat_messages" ("id" text NOT NULL, "guild_id" text NOT NULL, "user_id" text NOT NULL, "body" text NOT NULL, "created_at" timestamptz NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "guild_chat_messages_guild_id_fkey" FOREIGN KEY ("guild_id") REFERENCES "guilds" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT "guild_chat_messages_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT "guild_chat_messages_body_check" CHECK ((length(body) > 0) AND (length(body) <= 1000)));
-- Create index "guild_chat_messages_guild_id_created_at_idx" to table: "guild_chat_messages"
CREATE INDEX "guild_chat_messages_guild_id_created_at_idx" ON "guild_chat_messages" ("guild_id", "created_at" DESC, "id" DESC);
