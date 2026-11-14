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

func TestSkillFinder(t *testing.T) {
  skills, err := FindSkills()
  expectedSkillCount := 1
  expectedSkillName := "randomizer"

  if err != nil {
    t.Errorf("Skill finder failed, got non-nil error: %s", err)
  }

  if len(skills) != expectedSkillCount {
    t.Errorf("Skill finder failed, got %d skills, expected %d", len(skills), expectedSkillCount)
  }

  if skills[0].Name() != expectedSkillName {
    t.Errorf("Skill finder fetched skill %s expected %s", skills[0].Name(), expectedSkillName)
  }
}
