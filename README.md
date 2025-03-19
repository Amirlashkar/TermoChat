# Why? Why Golang?

I tried to clone the idea of chatting safely on your terminal that i saw on movies like Mr.Robot & The Day of Jackal.
Also i wanted to try new language and golang seemed good choice to create an API for me.

## Done part

- **Database management**
- **Room and users components**
- **Endpoints except websocket**

## Its not complete

- **Fix needed**
    - **Websocket client**: messages are being sent and recieved well but clients can't connect to same room because of one error that causing throwing client out of room.
- **Extra impovements**
    - **SQL injection security**: Database struct is easily exposed into injections
