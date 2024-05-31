# my-wolverine-events

Playing around with rec systems. In development. Hosted (sometimes) at [mywolverine.events](https://mywolverine.events/).

## How to run 
Create `.env` file with 
```
GOOGLE_CLIENT_ID=[The app's client ID]
JWT_SECRET=[A random key. For example, generated w/openssl rand -hex 32]
```
As sudo, run `sh serve.sh` or ` nohup ./serve.sh &` to leave running over ssh.


## TODOs
- change prints to logs
- fix datetime in cards