# assignment_demo_2023

![Tests](https://github.com/TikTokTechImmersion/assignment_demo_2023/actions/workflows/test.yml/badge.svg)

Completed Tiktok backend assignment
Author: Tan Lip Wei

### Assignment outline
This project serves as a backend for an instant messaging system.

The http-server exposes 2 endpoints
1. /api/send (To send messages)

Body Parameters 
```azure
{
    "chat": string,
    "text": string,
    "sender": string,
    "header": string
}
```
2. /api/pull (To pull messages)

Body Parameters
```azure
{
    "chat": string,
    "cursor": number,
    "limit": number,
    "reverse": bool
}
```