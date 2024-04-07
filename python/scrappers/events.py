import requests, argparse


def get_events(url):
    """Gets the events from the url and returns them as a list of dictionaries"""

    assert 'umich.edu' in url, 'The url must be from the umich events page'
    assert 'json' in url, 'The url must be the JSON version of the events page'

    r = requests.get(url)
    try:
        return r.json()
    except:
        raise Exception('Could not parse events from the url. Make sure its the JSON version of the events page.')


if __name__ == '__main__':
    
    EVENTS_URL = 'https://events.umich.edu/week/json?v=2'
    print(get_events(EVENTS_URL)) 