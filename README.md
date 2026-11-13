# Introduction

This is a personal project for me to understand how tool calling, skills, and agents work. I am building my own, self-contained harness from basic principles, that is going to be used to drive my locally hosted model (currently qwen3.5:9b). The first-cut functionality provides support for just basic tool calling.

I am also going to be using this project to note down my [observations](#Observations) about how tool calling, skill calling, and other features work. The central question driving my study is - how efficient can we make our harness, and can we operate it more efficiently than some of the open source harnesses out there?

## Organization
The central program for the harness is just called `harness.go`. I would like this to be as modular as possible (so that I can build different implementations that I can use to experiment) so I am building skill support in the 'skills' directory. The skills directory is also going to be used to hold all the skills.

## Observations
Different harnesses [that we tested](https://github.com/BlaidddDrwg/harness-bench) handle tool calling, skills, and agents differently. This shows up as a vast difference in the context size and the token overhead for supporting this functionality.
