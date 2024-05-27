# Monitors the site and notifies via ntfy if something is up
# Ran with GitHub Actions (see .github/workflows/monitor.yml)

import requests, sys

def notify(msg, url):
    if not CHANNEL:
        print('There is no nfty channel!')
        exit(1)

    print(f'Notifying {msg}')
    requests.post(
        f"https://ntfy.sh/{CHANNEL}",
        data = msg.encode('utf-8'),
        headers = {
            "Click": url,
            "Tags": "warning,mywolverineevents"
    })

def req(url, exit_on_error = True):
    r = None
    try:
        r = requests.get(url)
        if r.status_code != 200:
            raise Exception(f'Status code {r.status_code} for {url}')
    except Exception as e:
        print(f'Error trying to access {url}: {e}')
        notify(f'Site unreachable: {url}', url)
        if exit_on_error: exit(1)
    return r

if __name__ == '__main__':


    BASE_URL = 'https://mywolverine.events/'
    
    if len(sys.argv) != 2:
        print('wrong # of cmd args')
        exit(1)
    
    global CHANNEL
    CHANNEL = sys.argv[1]

    # Check site is up
    req(BASE_URL)