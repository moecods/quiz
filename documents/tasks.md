# Quiz App

This App is for managing Quizes

First we have a strcuture for store quizes, each quiz can be have a lot of questions.

Quiz has action list, add, update, delete, get

|action|  test |service|
|------|-------|-------|
| list |&cross;|&cross;|
|  add |&cross;|&check;|
|update|&cross;|&cross;|
|  get |&check;|&check;|
|delete|&cross;|&cross;|

tasks:

- [x] store quiz with questions (&cross;test)
- [x] edit quiz with questions  (&cross;test)
- [x] remove quiz  (&cross;test)
- [x] get list of quiz
    - [ ] taggable
    - [ ] filter by tags
    - [ ] filter by quiz title
    - [ ] pagination
- [ ] cache popular quizes
- [ ] lock quiz after started (prevent update or delete)
- [ ] archive quiz with participant after 6 month
- [ ] create quiz module and import to app

Participant has action 

- [x] register participants (&cross;test)
- [x] save participants answers (&cross;test)
- [ ] save single answer
- [ ] lock participant (stop quiz before end)
- [ ] get list of participant participated in specific quiz
- [ ] use redis for save answers faster

- [ ] score answers with openAi
- [ ] score answers with human