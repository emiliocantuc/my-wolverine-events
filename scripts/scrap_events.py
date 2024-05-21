import requests, urllib.parse, json, argparse, os


def get_events(url):
    """Gets the events from the url and returns them as a list of dictionaries"""

    assert 'umich.edu' in url, 'The url must be from the umich events page'
    assert 'json' in url, 'The url must be the JSON version of the events page'

    r = requests.get(url)
    try:
        return r.json()
    except:
        raise Exception('Could not parse events from the url. Make sure its the JSON version of the events page.')

def get_cal_links(events_json):

    for event in events_json:
        try:
            permalink = event['permalink']
            html = requests.get(permalink).text
            gcal = html.split('googleCal_href": "')[1].split('"')[0]
            ical = html.split('iCal_href": "')[1].split('"')[0]
            ical = urllib.parse.urljoin(permalink, ical)
            event['gcal_link'] = gcal
            event['ical_link'] = ical
        except: pass



if __name__ == '__main__':

    parser = argparse.ArgumentParser(description = 'Get events from the umich events page')
    parser.add_argument('--eventsURL', type = str, help = 'Events json endpoint', default = 'https://events.umich.edu/week/json?v=2', required = False)
    parser.add_argument('--output', type = str, help = 'Output file to save the events', default = 'data/events.json', required = False)
    args = parser.parse_args()

    # Check if output directory exists
    assert os.path.exists(os.path.dirname(args.output)), 'The output directory does not exist'
    
    events = get_events(args.eventsURL)
    get_cal_links(events)

    print(f'Found {len(events)} events')
    with open(args.output, 'w+') as f:
        print(f'Saving events to {args.output}')
        f.write(json.dumps(events, indent = 4))