# Fetches events from umich API and saves to db

import sqlite3
import requests, urllib.parse, argparse, os


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
            # ical = html.split('iCal_href": "')[1].split('"')[0]
            # ical = urllib.parse.urljoin(permalink, ical)
            event['gcal_link'] = gcal
            # event['ical_link'] = ical
        except: pass

def insert_event(cursor, event):
    cursor.execute('''
    INSERT INTO events (nweek, title, event_description, event_date, type, permalink, building_name, building_id, gcal_link, umich_id)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    ''', (
        event.get('nweek'),
        event.get('combined_title'),
        event.get('description'),
        event.get('datetime_start'),
        event.get('event_type'),
        event.get('permalink'),
        event.get('building_name'),
        event.get('building_official_id'),
        event.get('gcal_link'),
        event.get('id')
    ))


if __name__ == '__main__':

    parser = argparse.ArgumentParser(description = 'Gets events from the umich events API and saves to database')
    parser.add_argument('--eventsURL', type = str, help = 'Events json endpoint', default = 'https://events.umich.edu/week/json?v=2', required = False)
    parser.add_argument('--output', type = str, help = 'Output db file to save the events', default = 'data/main.db', required = False)
    args = parser.parse_args()

    # Check if output db exists
    assert os.path.exists(os.path.dirname(args.output)), 'The output db does not exist'
    
    events = get_events(args.eventsURL)
    print(f'Found {len(events)} events. Getting calendar links ... ', end = '')
    get_cal_links(events) # TODO what if this fails
    print('Done')

    print('Inserting events in db ... ', end = '')
    conn = sqlite3.connect(args.output)
    cursor = conn.cursor()

    try:
        cursor.execute('BEGIN TRANSACTION;')

        # Get the week number
        cursor.execute('SELECT MAX(nweek) FROM events;')
        max_nweek = cursor.fetchone()[0] or 0 
        nweek = max_nweek + 1
        print(f'nweek = {nweek} ', end = '')

        # Insert events
        for event in events:
            event['nweek'] = nweek
            insert_event(conn, event)

        conn.commit()
        print('Committed')

    except Exception as e:
        conn.rollback()
        print('Error saving events:', e)
        exit(1)
    
    finally:
        conn.close()
    