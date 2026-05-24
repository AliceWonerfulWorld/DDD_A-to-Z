-- Modify "user_profiles" table
ALTER TABLE "user_profiles" DROP CONSTRAINT "fk_user_profiles_selected_badge", ADD CONSTRAINT "user_profiles_selected_badge_slug_fkey" FOREIGN KEY ("selected_badge_slug") REFERENCES "badges" ("slug") ON UPDATE NO ACTION ON DELETE NO ACTION;
