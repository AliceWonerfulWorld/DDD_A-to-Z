ALTER TABLE user_profiles ADD COLUMN selected_badge_slug TEXT;
ALTER TABLE user_profiles ADD CONSTRAINT fk_user_profiles_selected_badge
  FOREIGN KEY (user_id, selected_badge_slug) REFERENCES user_badges(user_id, badge_slug);