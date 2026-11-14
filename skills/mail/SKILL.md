---
name: mail_finder
description: Finds email
---

# Functionality

This skill searches and finds emails. It does this by executing the "find_mail.go" executable in this directory.

You can call the mail_finder using this invocation:
`go run find_mail.go <topic>` where the <topic> is the search term to search for email. If no search term is specified, then it will get the latest email on the stack.
