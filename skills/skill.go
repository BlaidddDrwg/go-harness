package skill

import (
  "bufio"
  "os"
  "strings"
)

// Given the skill directory, parse out the skill into something that we can use

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
