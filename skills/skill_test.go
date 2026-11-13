package skill

import "testing"

func TestBasicSkill(t *testing.T) {
  testSkill, err := initializeSkill("random")
  testSkillName := "randomizer"

  if err != nil {
    t.Errorf("Skill initializer failed, got non-nil error: %s", err)
  }

  if testSkill.Name() != testSkillName {
    t.Errorf("Skill initializer failed, got skill name [%s] expected [%s]", testSkill.Name(), testSkillName)
  }
}
