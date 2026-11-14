package skill

import (
  "bufio"
  "errors"
  "log"
  "os"
  "strings"
)

type Skill struct {
  name string
  description string
  payload string
}

func (s *Skill) Name() string {
  return s.name
}

func (s *Skill) Description() string {
  return s.description
}

func (s *Skill) Payload() string {
  return s.payload
}

var COMMENT string = "---"
var NAME string = "name:"
var DESCRIPTION string = "description:"

func initializeSkill(skillDirectory string) (*Skill, error) {
  skillFile, err := os.OpenInRoot(skillDirectory, "SKILL.md")
  if err != nil {
    return nil, err
  }

  skillScanner := bufio.NewScanner(skillFile)
  skillScanner.Split(bufio.ScanLines)
  var name, description, payload string

  for skillScanner.Scan() {
    skillLine := strings.TrimSpace(skillScanner.Text())

    if strings.HasPrefix(skillLine, COMMENT) {
      continue
    }

    if strings.HasPrefix(skillLine, NAME) {
      name = strings.TrimSpace(skillLine[len(NAME) + 1:])
      continue
    }

    if strings.HasPrefix(skillLine, DESCRIPTION) {
      description = strings.TrimSpace(skillLine[len(DESCRIPTION) + 1:])
      continue
    }

    payload += skillLine
  }

  return &Skill{ name, description, payload }, nil
}

func FindSkills(skillsDir string) ([]*Skill, error) {
  var skills []*Skill
  log.Print("----------------------------------------------------------------")
  log.Printf("Searching directory: %s", skillsDir)
  log.Print("----------------------------------------------------------------")
  SkillsDir, err := os.Open(skillsDir)
  if err != nil {
    // No skills directory, so likely no skills
    return skills, errors.New("No skills directory found")
  }

  _skills, err := SkillsDir.Readdir(-1)
  defer SkillsDir.Close()

  if err != nil {
    // Skills directory doesn't contain any skills - exit without doing anything
    return skills, errors.New("No skills found in skill directory")
  }

  for _, _skill := range _skills {
    log.Print("Found: " + _skill.Name())
    if err != nil {
      log.Printf("Couldn't stat info for file: %s\n", _skill.Name())
      continue
    }

    if _skill.IsDir() {
      log.Printf("Found skill: %s", _skill.Name())
      skill, err := initializeSkill(skillsDir + "/" + _skill.Name())

      if err != nil {
        log.Printf("Error while trying to initialize skill: %s", _skill.Name())
        continue
      }

      if skill != nil {
        skills = append(skills, skill)
      } else {
        log.Printf("Could not initialize skill [%s]", _skill.Name())
      }
    }
  }

  return skills, nil
}
